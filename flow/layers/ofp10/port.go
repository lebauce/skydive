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
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

type PortConfig uint32
type PortState uint32

const (
	OFP_MAX_PORT_NAME_LEN = 16

	OFPPC_PORT_DOWN    = 1 << 0
	OFPPC_NO_STP       = 1 << 1
	OFPPC_NO_RECV      = 1 << 2
	OFPPC_NO_RECV_STP  = 1 << 3
	OFPPC_NO_FLOOD     = 1 << 4
	OFPPC_NO_FWD       = 1 << 5
	OFPPC_NO_PACKET_IN = 1 << 6

	// ofp_port_state
	OFPPS_LINK_DOWN   = 1 << 0
	OFPPS_STP_LISTEN  = 0 << 8
	OFPPS_STP_LEARN   = 1 << 8
	OFPPS_STP_FORWARD = 2 << 8
	OFPPS_STP_BLOCK   = 3 << 8
	OFPPS_STP_MASK    = 3 << 8

	// ofp_port
	OFPP_MAX        = 0xff00
	OFPP_IN_PORT    = 0xfff8
	OFPP_TABLE      = 0xfff9
	OFPP_NORMAL     = 0xfffa
	OFPP_FLOOD      = 0xfffb
	OFPP_ALL        = 0xfffc
	OFPP_CONTROLLER = 0xfffd
	OFPP_LOCAL      = 0xfffe
	OFPP_NONE       = 0xffff

	// ofp_port_features
	OFPPF_10MB_HD    = 1 << 0
	OFPPF_10MB_FD    = 1 << 1
	OFPPF_100MB_HD   = 1 << 2
	OFPPF_100MB_FD   = 1 << 3
	OFPPF_1GB_HD     = 1 << 4
	OFPPF_1GB_FD     = 1 << 5
	OFPPF_10GB_FD    = 1 << 6
	OFPPF_COPPER     = 1 << 7
	OFPPF_FIBER      = 1 << 8
	OFPPF_AUTONEG    = 1 << 9
	OFPPF_PAUSE      = 1 << 10
	OFPPF_PAUSE_ASYM = 1 << 11
)

/* Description of a physical port
   PhysicalPort is modeled after ofp_phy_port */
type PhysicalPort struct {
	PortNumber uint16
	MAC        net.HardwareAddr
	Name       string
	Config     uint32
	State      uint32
	Current    uint32
	Advertised uint32
	Supported  uint32
	Peer       uint32
}

func (p *PhysicalPort) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(48)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], p.PortNumber)

	if err := openflow.SerializeMAC(data[2:8], p.MAC); err != nil {
		return err
	}

	copy(data[8:8+OFP_MAX_PORT_NAME_LEN], p.Name)

	binary.BigEndian.PutUint32(data[24:28], p.Config)
	binary.BigEndian.PutUint32(data[28:32], p.State)

	binary.BigEndian.PutUint32(data[32:36], p.Current)
	binary.BigEndian.PutUint32(data[36:40], p.Advertised)
	binary.BigEndian.PutUint32(data[40:44], p.Supported)
	binary.BigEndian.PutUint32(data[44:48], p.Peer)

	return nil
}

func (p *PhysicalPort) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 48 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	p.PortNumber = binary.BigEndian.Uint16(data[0:2])
	p.MAC = net.HardwareAddr(data[2:8])

	p.Name = string(bytes.Trim(data[8:8+OFP_MAX_PORT_NAME_LEN], "\x00"))

	p.Config = binary.BigEndian.Uint32(data[24:28])
	p.State = binary.BigEndian.Uint32(data[28:32])

	p.Current = binary.BigEndian.Uint32(data[32:36])
	p.Advertised = binary.BigEndian.Uint32(data[36:40])
	p.Supported = binary.BigEndian.Uint32(data[40:44])
	p.Peer = binary.BigEndian.Uint32(data[44:48])

	return nil
}

type PortStatusReason uint8

const (
	OFPPR_ADD PortStatusReason = iota
	OFPPR_DELETE
	OFPPR_MODIFY
)

// PortStatus is modeled after ofp_port_status
type PortStatus struct {
	Reason      PortStatusReason
	Description PhysicalPort
}
