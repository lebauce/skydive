/*
 * Copyright (C) 2018 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package ovsdb

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers"
	"github.com/skydive-project/skydive/flow/layers/ofp10"
	"github.com/skydive-project/skydive/flow/layers/openflow"
	"github.com/skydive-project/skydive/logging"
)

const (
	echoDuration = 5
)

var (
	// ErrContextDone is returned what the context was done or canceled
	ErrContextDone = errors.New("Context was terminated")
	// ErrConnectionTimeout is returned when a timeout was reached when trying to connect
	ErrConnectionTimeout = errors.New("Timeout while connecting")
	// ErrReaderChannelClosed is returned when the read channel was closed
	ErrReaderChannelClosed = errors.New("Reader channel was closed")
)

// Monitor describe a OVS monitor similar to the
// ovs-vsctl monitor watch: command
type Monitor struct {
	sync.RWMutex
	addr      string
	conn      net.Conn
	reader    *bufio.Reader
	ctx       context.Context
	msgChan   chan (openflow.Message)
	listeners []MonitorListener
}

type MonitorListener interface {
	OnMessage(openflow.Message)
}

func (m *Monitor) connect(addr string) (net.Conn, error) {
	split := strings.SplitN(addr, ":", 2)
	if len(split) < 2 {
		return nil, fmt.Errorf("Invalid monitor scheme: '%s'", addr)
	}
	scheme, addr := split[0], split[1]

	switch scheme {
	case "tcp":
		return net.Dial(scheme, addr)
	default:
		return nil, fmt.Errorf("Unsupported connection scheme '%s'", scheme)
	}
}

func (m *Monitor) handshake() (openflow.Message, error) {
	m.sendHello()

	select {
	case <-m.ctx.Done():
		return nil, ErrContextDone
	case <-time.After(30 * time.Second):
		return nil, ErrConnectionTimeout
	case msg, ok := <-m.msgChan:
		if !ok {
			return nil, ErrReaderChannelClosed
		}

		helloMsg, ok := msg.(*openflow.HelloMsg)
		if !ok {
			return nil, fmt.Errorf("Expected a first message of type Hello")
		}

		return helloMsg, nil
	}
}

func (m *Monitor) handleLoop() error {
	for {
		select {
		case <-m.ctx.Done():
			return ErrContextDone
		case <-time.After(30 * time.Second):
			return ErrConnectionTimeout
		case msg, ok := <-m.msgChan:
			if !ok {
				return ErrReaderChannelClosed
			}

			m.RLock()
			for _, listener := range m.listeners {
				listener.OnMessage(msg)
			}
			m.RUnlock()
		}
	}
}

func (m *Monitor) readLoop() {
	type Timeout interface {
		Timeout() bool
	}

	echoTicker := time.NewTicker(time.Second * echoDuration)
	defer echoTicker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-echoTicker.C:
			m.sendEcho()
		default:
			data, err := m.reader.Peek(8)
			if err != nil {
				if _, ok := err.(Timeout); !ok {
					logging.GetLogger().Errorf("Failed to read packet: %s", err)
					return
				}
				continue
			}

			header := openflow.Header{}
			if err := header.DecodeFromBytes(data, gopacket.NilDecodeFeedback); err != nil {
				logging.GetLogger().Debugf("Ignoring non OpenFlow packet: %s", err)
				continue
			}

			data = make([]byte, header.Length)
			n, err := m.reader.Read(data)
			if n != int(header.Length) {
				logging.GetLogger().Errorf("Failed to read full OpenFlow message: %s", err)
				continue
			}

			var msg openflow.Message
			if header.Type.IsImmutable() {
				msg, err = openflow.DecodeImmutableMessage(header.Type, data[8:], gopacket.NilDecodeFeedback)
			} else {
				switch header.Version {
				case openflow.OPENFLOW_1_0:
					msg, err = ofp10.DecodeMessage(header.Type, data[8:], gopacket.NilDecodeFeedback)
					if err != nil {
						logging.GetLogger().Errorf("Failed to decode OpenFlow message: %s", err)
						continue
					}
				default:
					logging.GetLogger().Warningf("Unsupported OpenFlow version '%s", header.Version)
					continue
				}
			}

			logging.GetLogger().Infof("Decoded message %+v (%+v)", msg, header)
			m.msgChan <- msg
		}
	}
}

func (m *Monitor) sendMessage(msg openflow.Message) error {
	b := gopacket.NewSerializeBuffer()
	ofl := layers.OpenFlowLayer{
		Header: openflow.Header{
			Version: openflow.OPENFLOW_1_0,
			Type:    msg.MessageType(),
		},
		Message: msg,
	}
	if err := ofl.SerializeTo(b, gopacket.SerializeOptions{}); err != nil {
		return err
	}
	_, err := m.conn.Write(b.Bytes())
	return err
}

func (m *Monitor) sendEcho() error {
	return m.sendMessage(&openflow.EchoRequestMsg{})
}

func (m *Monitor) sendHello() error {
	return m.sendMessage(&openflow.HelloMsg{})
}

func (m *Monitor) RegisterListener(listener MonitorListener) {
	m.Lock()
	defer m.Unlock()

	m.listeners = append(m.listeners, listener)
}

// Start monitoring the OVS bridge
func (m *Monitor) Start(ctx context.Context) (err error) {
	m.conn, err = m.connect(m.addr)
	if err != nil {
		return err
	}

	m.reader = bufio.NewReader(m.conn)
	m.ctx = ctx

	go m.readLoop()

	_, err = m.handshake()
	if err != nil {
		return err
	}

	go m.handleLoop()

	logging.GetLogger().Infof("Successfully connected to OVS")

	return nil
}

// Stop monitoring the OVS bridge
func (m *Monitor) Stop() error {
	return nil
}

// NewMonitor returns a new OVS monitor using either a UNIX socket or a TCP socket
func NewMonitor(addr string) (*Monitor, error) {
	return &Monitor{
		addr:    addr,
		msgChan: make(chan openflow.Message, 500),
	}, nil
}
