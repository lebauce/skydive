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

package ofp10

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

// Capabilities is modeled after ofp_capabilities
type Capabilities uint32

const (
	OFPC_FLOW_STATS   = 1 << 0
	OFPC_TABLE_STATS  = 1 << 1
	OFPC_PORT_STATS   = 1 << 2
	OFPC_STP          = 1 << 3
	OFPC_RESERVED     = 1 << 4
	OFPC_IP_REASM     = 1 << 5
	OFPC_QUEUE_STATS  = 1 << 6
	OFPC_ARP_MATCH_IP = 1 << 7
)

type DatapathID uint64

func (DatapathID) MAC() net.HardwareAddr {
	return net.HardwareAddr{}
}

type ActionTypes uint32

type FeaturesRequestMsg struct {
	openflow.EmptyMsg
}

func (m *FeaturesRequestMsg) MessageType() openflow.MessageType {
	return OFPT_FEATURES_REQUEST
}

// FeaturesReplyMsg is modeled after ofp_switch_features
type FeaturesReplyMsg struct {
	DatapathID      DatapathID
	BufferedPackets uint32
	TableCount      uint8
	Capabilities    Capabilities
	Actions         ActionTypes
	Ports           []PhysicalPort
}

func (m *FeaturesReplyMsg) MessageType() openflow.MessageType {
	return OFPT_FEATURES_REPLY
}

func (m *FeaturesReplyMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(24)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(data[0:8], uint64(m.DatapathID))
	binary.BigEndian.PutUint32(data[8:12], m.BufferedPackets)
	data[12] = m.TableCount

	binary.BigEndian.PutUint32(data[16:20], uint32(m.Capabilities))
	binary.BigEndian.PutUint32(data[20:24], uint32(m.Actions))

	for _, port := range m.Ports {
		if err := port.SerializeTo(b, opts); err != nil {
			return nil
		}
	}

	return nil
}

func (m *FeaturesReplyMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 24 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.DatapathID = DatapathID(binary.BigEndian.Uint64(data[0:8]))
	m.BufferedPackets = binary.BigEndian.Uint32(data[8:12])
	m.TableCount = data[12]

	m.Capabilities = Capabilities(binary.BigEndian.Uint32(data[16:20]))
	m.Actions = ActionTypes(binary.BigEndian.Uint32(data[20:24]))

	if ((len(data) - 24) % 48) != 0 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	portCount := (len(data) - 24) / 48
	m.Ports = make([]PhysicalPort, portCount)
	for i := range m.Ports {
		var port PhysicalPort
		port.DecodeFromBytes(data[24+(i*48):24+((i+1)*48)], df)
		m.Ports[i] = port
	}

	return nil
}
