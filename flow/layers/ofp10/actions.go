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
)

// Action is modeled after ofp_action_type
type ActionType uint16

const (
	OFPAT_OUTPUT ActionType = iota
	OFPAT_SET_VLAN_VID
	OFPAT_SET_VLAN_PCP
	OFPAT_STRIP_VLAN
	OFPAT_SET_DL_SRC
	OFPAT_SET_DL_DST
	OFPAT_SET_NW_SRC
	OFPAT_SET_NW_DST
	OFPAT_SET_NW_TOS
	OFPAT_SET_TP_SRC
	OFPAT_SET_TP_DST
	OFPAT_ENQUEUE
)

// ActionHeader is modeled after ofp_action_header
type ActionHeader struct {
	Type   ActionType
	Length uint16
}

func (h *ActionHeader) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], uint16(h.Type))
	binary.BigEndian.PutUint16(data[2:4], h.Length)
	binary.BigEndian.PutUint32(data[4:8], 0)

	return nil
}

func (h *ActionHeader) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	h.Type = ActionType(binary.BigEndian.Uint16(data[0:2]))
	h.Length = binary.BigEndian.Uint16(data[2:4])

	return nil
}

// ActionOutput is modeled after ofp10_action_output
type ActionOutput struct {
	Port      uint16
	MaxLength uint16
}

// ActionEnqueue is modeled after ofp10_action_enqueue
type ActionEnqueue struct {
	Port    uint16
	QueueID uint32
}

// ActionVLANID is modeled after ofp_action_vlan_vid
type ActionVLANID struct {
	VLANID uint16
}

// ActionVLANPCP is modeled after ofp_action_vlan_pcp
type ActionVLANPCP struct {
	VLANPCP uint8
}

// ActionDLAddr is modeled after ofp_action_dl_addr
type ActionDLAddr struct {
	DLMAC net.HardwareAddr
}

// ActionNWAddr is modeled after ofp_action_nw_addr
type ActionNWAddr struct {
	NWAddr net.IP
}

// ActionTPPort is modeled after ofp_action_tp_port
type ActionTPPort struct {
	TPPort uint16
}

// ActionVendor is modeled after ofp_action_vendor_header
type ActionVendor struct {
	Vendor uint32
}
