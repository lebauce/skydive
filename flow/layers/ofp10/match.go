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

type FlowWildcards uint32

const (
	OFPFW_IN_PORT      = 1 << 0
	OFPFW_DL_VLAN      = 1 << 1
	OFPFW_DL_SRC       = 1 << 2
	OFPFW_DL_DST       = 1 << 3
	OFPFW_DL_TYPE      = 1 << 4
	OFPFW_NW_PROTO     = 1 << 5
	OFPFW_TP_SRC       = 1 << 6
	OFPFW_TP_DST       = 1 << 7
	OFPFW_NW_SRC_SHIFT = 8
	OFPFW_NW_SRC_BITS  = 6
	OFPFW_NW_SRC_MASK  = ((1 << OFPFW_NW_SRC_BITS) - 1) << OFPFW_NW_SRC_SHIFT
	OFPFW_NW_SRC_ALL   = 32 << OFPFW_NW_SRC_SHIFT
	OFPFW_NW_DST_SHIFT = 14
	OFPFW_NW_DST_BITS  = 6
	OFPFW_NW_DST_MASK  = ((1 << OFPFW_NW_DST_BITS) - 1) << OFPFW_NW_DST_SHIFT
	OFPFW_NW_DST_ALL   = 32 << OFPFW_NW_DST_SHIFT
	OFPFW_DL_VLAN_PCP  = 1 << 20
	OFPFW_NW_TOS       = 1 << 21
	OFPFW_ALL          = ((1 << 22) - 1)
)

type Match struct {
	Wildcards  FlowWildcards
	InPort     uint16
	SrcMAC     net.HardwareAddr
	DstMAC     net.HardwareAddr
	VLAN       uint16
	VLANPCP    uint8
	Type       uint16
	IPTOS      uint8
	IPProtocol uint8
	SrcIP      net.IP
	DstIP      net.IP
	SrcPort    uint16
	DstPort    uint16
}

func (m *Match) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(40)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], uint32(m.Wildcards))
	binary.BigEndian.PutUint16(data[4:6], m.InPort)

	if err := openflow.SerializeMAC(data[6:12], m.SrcMAC); err != nil {
		return err
	}

	if err := openflow.SerializeMAC(data[12:18], m.DstMAC); err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[18:20], m.VLAN)
	data[20] = m.VLANPCP

	binary.BigEndian.PutUint16(data[22:24], m.Type)
	data[24] = m.IPTOS
	data[25] = m.IPProtocol

	copy(data[28:32], m.SrcIP.To4())
	copy(data[32:36], m.DstIP.To4())

	binary.BigEndian.PutUint16(data[36:38], m.SrcPort)
	binary.BigEndian.PutUint16(data[38:40], m.DstPort)

	return nil
}

func (m *Match) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 40 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Wildcards = FlowWildcards(binary.BigEndian.Uint32(data[0:4]))
	m.InPort = binary.BigEndian.Uint16(data[4:6])
	m.SrcMAC = net.HardwareAddr(data[6:12])
	m.DstMAC = net.HardwareAddr(data[12:18])
	m.VLAN = binary.BigEndian.Uint16(data[18:20])
	m.VLANPCP = data[20]
	m.Type = binary.BigEndian.Uint16(data[22:24])
	m.IPTOS = data[24]
	m.IPProtocol = data[25]
	m.SrcIP = net.IP(data[28:32])
	m.DstIP = net.IP(data[32:36])
	m.SrcPort = binary.BigEndian.Uint16(data[36:38])
	m.DstPort = binary.BigEndian.Uint16(data[38:40])

	return nil
}

type NXMEntry interface {
	openflow.Serializable
}

type NXMatch struct {
	Entries []NXMEntry
}

func (m *NXMatch) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	for _, entry := range m.Entries {
		if err := entry.SerializeTo(b, opts); err != nil {
			return err
		}
	}
	return nil
}

func (m *NXMatch) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	for index := 0; index < len(data); {
		header := binary.BigEndian.Uint32(data[0:4])
		length := uint8(header)
		if entry, found := nxmEntries[header]; found {
			if err := entry.DecodeFromBytes(data[4:4+length], df); err != nil {
				return err
			}
			m.Entries = append(m.Entries, entry)
		}
		data = data[4+int(length):]
	}

	return nil
}

type EthSrcField struct {
	MAC net.HardwareAddr
}

func (f *EthSrcField) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(openflow.OFP_ETH_ALEN)
	if err != nil {
		return err
	}

	return openflow.SerializeMAC(data, f.MAC)
}

func (f *EthSrcField) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < openflow.OFP_ETH_ALEN {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	f.MAC = net.HardwareAddr(data[0:6])
	return nil
}

func registerNXMEntry(class uint16, field uint8, hasMask bool, length uint8, entry NXMEntry) {
	value := (uint32(class) << 16) | (uint32(field) << 9) | uint32(length)
	if hasMask {
		value |= 0x100
	}
	nxmEntries[value] = entry
}

func init() {
	nxmEntries = make(map[uint32]NXMEntry)
	registerNXMEntry(0x0, 2, false, 6, &EthSrcField{})
}
