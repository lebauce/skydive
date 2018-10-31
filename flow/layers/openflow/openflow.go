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

package openflow

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
)

type Version uint8

const (
	OPENFLOW_1_0 Version = 0x01
	OPENFLOW_1_1         = 0x02
	OPENFLOW_1_2         = 0x03
	OPENFLOW_1_3         = 0x04
	OPENFLOW_1_4         = 0x05
	OPENFLOW_1_5         = 0x06
)

const (
	/* Immutable messages. */
	OFPT_HELLO MessageType = iota
	OFPT_ERROR
	OFPT_ECHO_REQUEST
	OFPT_ECHO_REPLY
	OFPT_VENDOR
)

const (
	OFP_ETH_ALEN = 6

	OfMinimumRecordSizeInBytes = 8
)

type MessageType uint8

func (t MessageType) IsImmutable() bool {
	return t >= OFPT_HELLO && t <= OFPT_VENDOR
}

type Header struct {
	Version Version
	Type    MessageType
	Length  uint16
	Xid     uint32
}

func (h *Header) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < OfMinimumRecordSizeInBytes {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	h.Version = Version(data[0])
	h.Type = MessageType(data[1])
	h.Length = binary.BigEndian.Uint16(data[2:4])
	h.Xid = binary.BigEndian.Uint32(data[4:8])

	return nil
}

func (h *Header) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.PrependBytes(OfMinimumRecordSizeInBytes)
	if err != nil {
		return err
	}

	data[0] = uint8(h.Version)
	data[1] = uint8(h.Type)
	binary.BigEndian.PutUint32(data[4:8], h.Xid)

	return nil
}

type Serializable interface {
	SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error
	DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error
}

type Message interface {
	Serializable
	MessageType() MessageType
}

func SerializeMAC(data []byte, mac net.HardwareAddr) error {
	if len(mac) != OFP_ETH_ALEN {
		return fmt.Errorf("Invalid MAC address: %s", mac.String())
	}
	copy(data, mac)
	return nil
}

type EmptyMsg struct {
}

func (m *EmptyMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (m *EmptyMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

type BasicMsg struct {
	Body []byte
}

func (m *BasicMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(len(m.Body))
	if err != nil {
		return err
	}
	copy(data, m.Body)
	return nil
}

func (m *BasicMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	m.Body = data
	return nil
}

type HelloMsg struct {
	BasicMsg
}

func (m *HelloMsg) MessageType() MessageType {
	return OFPT_HELLO
}

func NewHelloMsg(data []byte) *HelloMsg {
	return &HelloMsg{BasicMsg{Body: data}}
}

type EchoRequestMsg struct {
	BasicMsg
}

func (m *EchoRequestMsg) MessageType() MessageType {
	return OFPT_ECHO_REQUEST
}

type EchoReplyMsg struct {
	BasicMsg
}

func (m *EchoReplyMsg) MessageType() MessageType {
	return OFPT_ECHO_REPLY
}

type VendorMsg struct {
	Vendor uint32
	Body   []byte
}

func (m *VendorMsg) MessageType() MessageType {
	return OFPT_VENDOR
}

func (m *VendorMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 4 {
		df.SetTruncated()
		return errors.New("OpenFlow vendor message too short")
	}

	m.Vendor = binary.BigEndian.Uint32(data[0:4])
	m.Body = data[4:]

	return nil
}

func (m *VendorMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.PrependBytes(4 + len(m.Body))
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], m.Vendor)
	copy(data[4:], m.Body)

	return nil
}

type UnhandledMessage struct {
	BasicMsg
	Type MessageType
}

func (m *UnhandledMessage) MessageType() MessageType {
	return m.Type
}

type UnhandledContent struct {
}

func (c *UnhandledContent) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (c *UnhandledContent) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func DecodeImmutableMessage(msgType MessageType, data []byte, df gopacket.DecodeFeedback) (msg Message, err error) {
	switch msgType {
	case OFPT_HELLO:
		msg = &HelloMsg{}
	case OFPT_ERROR:
		msg = &ErrorMsg{}
	case OFPT_ECHO_REQUEST:
		msg = &EchoRequestMsg{}
	case OFPT_ECHO_REPLY:
		msg = &EchoReplyMsg{}
	case OFPT_VENDOR:
		msg = &VendorMsg{}
	default:
		return nil, fmt.Errorf("Invalid immutable message type: %s", msgType)
	}

	if err := msg.DecodeFromBytes(data, df); err != nil {
		return nil, err
	}

	return msg, nil
}
