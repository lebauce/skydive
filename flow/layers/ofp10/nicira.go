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
	"time"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

const (
	NX_VENDOR_ID = 0x00002320

	NXFMF_INITIAL = 1 << 0
	NXFMF_ADD     = 1 << 1
	NXFMF_DELETE  = 1 << 2
	NXFMF_MODIFY  = 1 << 3
	NXFMF_ACTIONS = 1 << 4
	NXFMF_OWN     = 1 << 5

	NXFME_ADDED    = 0
	NXFME_DELETED  = 1
	NXFME_MODIFIED = 2
	NXFME_ABBREV   = 3

	NXST_FLOW = 2
)

var nxmEntries map[uint32]NXMEntry

// NiciraStatsMsg is modeled after nicira10_stats_msg
type NiciraStatsMsg struct {
	Subtype uint32
	Content openflow.Serializable
}

func (m *NiciraStatsMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], m.Subtype)
	binary.BigEndian.PutUint32(data[4:8], 0)

	return m.Content.SerializeTo(b, opts)
}

func (m *NiciraStatsMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Subtype = binary.BigEndian.Uint32(data[0:4])

	switch m.Subtype {
	case NXST_FLOW:
		m.Content = &NxFlowUpdateReply{}
	default:
		m.Content = &openflow.UnhandledContent{}
	}

	return m.Content.DecodeFromBytes(data[8:], df)
}

// NxFlowMonitorRequest is modeled after nx_flow_monitor_request
type NxFlowMonitorRequest struct {
	MonitorID uint32
	Flags     uint16
	OutPort   uint16
	MatchLen  uint16
	TableID   uint8
	NXMatch   NXMatch
}

func (r *NxFlowMonitorRequest) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(16)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], r.MonitorID)
	binary.BigEndian.PutUint16(data[4:6], r.Flags)
	binary.BigEndian.PutUint16(data[6:8], r.OutPort)
	binary.BigEndian.PutUint16(data[8:10], r.MatchLen)
	data[10] = r.TableID
	copy(data[11:16], []byte{0, 0, 0, 0, 0})

	if err := r.NXMatch.SerializeTo(b, opts); err != nil {
		return err
	}

	padding := int((r.MatchLen+7)/8*8 - r.MatchLen)
	data, err = b.AppendBytes(padding)
	if err != nil {
		return err
	}
	copy(data, bytes.Repeat([]byte{0}, padding))

	return nil
}

func (r *NxFlowMonitorRequest) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 16 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	r.MonitorID = binary.BigEndian.Uint32(data[0:4])
	r.Flags = binary.BigEndian.Uint16(data[4:6])
	r.OutPort = binary.BigEndian.Uint16(data[6:8])
	r.MatchLen = binary.BigEndian.Uint16(data[8:10])

	if r.MatchLen != 0 {
		if err := r.NXMatch.DecodeFromBytes(data[10:10+r.MatchLen], df); err != nil {
			return err
		}
	}

	return nil
}

type NxFlowUpdateReply struct {
	Updates []NxFlowUpdateEvent
}

func (r *NxFlowUpdateReply) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	for _, update := range r.Updates {
		if err := update.SerializeTo(b, opts); err != nil {
			return err
		}
	}

	return nil
}

func (r *NxFlowUpdateReply) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	for len(data) > 0 {
		var updateEvent NxFlowUpdateEvent
		if err := updateEvent.DecodeFromBytes(data, df); err != nil {
			return err
		}
		data = data[updateEvent.Length:]
		r.Updates = append(r.Updates, updateEvent)
	}
	return nil
}

// NxFlowUpdateEvent is modeled after nx_flow_update_header
type NxFlowUpdateEvent struct {
	Length  uint16
	Event   uint16
	Content openflow.Serializable
}

func (e *NxFlowUpdateEvent) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(4)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], e.Length)
	binary.BigEndian.PutUint16(data[2:4], e.Event)

	return e.Content.SerializeTo(b, opts)
}

func (e *NxFlowUpdateEvent) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 4 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	e.Length = binary.BigEndian.Uint16(data[0:2])
	e.Event = binary.BigEndian.Uint16(data[2:4])

	switch e.Event {
	case NXFME_ADDED, NXFME_MODIFIED, NXFME_DELETED:
		e.Content = &NxFlowUpdateFull{}
	case NXFME_ABBREV:
		e.Content = &openflow.UnhandledContent{}
	}

	return e.Content.DecodeFromBytes(data[4:], df)
}

// NxFlowUpdateFull is modeled after nx_flow_update_full
type NxFlowUpdateFull struct {
	Reason      uint16
	Priority    uint16
	IdleTimeout time.Duration
	HardTimeout time.Duration
	MatchLen    uint16
	TableID     uint8
	Cookie      uint64
	NXMatch     NXMatch
}

func (e *NxFlowUpdateFull) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(24)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], e.Reason)
	binary.BigEndian.PutUint16(data[2:4], e.Priority)
	binary.BigEndian.PutUint16(data[4:6], uint16(e.IdleTimeout.Seconds()))
	binary.BigEndian.PutUint16(data[6:8], uint16(e.HardTimeout.Seconds()))
	binary.BigEndian.PutUint16(data[8:10], e.MatchLen)
	data[10] = e.TableID
	data[11] = 0
	binary.BigEndian.PutUint64(data[12:20], uint64(e.Cookie))

	return nil
}

func (e *NxFlowUpdateFull) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 24 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	e.Reason = binary.BigEndian.Uint16(data[0:2])
	e.Priority = binary.BigEndian.Uint16(data[2:4])
	e.IdleTimeout = time.Second * time.Duration(binary.BigEndian.Uint16(data[4:6]))
	e.HardTimeout = time.Second * time.Duration(binary.BigEndian.Uint16(data[6:8]))
	e.MatchLen = binary.BigEndian.Uint16(data[8:10])
	e.TableID = data[10]
	e.Cookie = binary.BigEndian.Uint64(data[12:20])

	if e.MatchLen != 0 {
		if err := e.NXMatch.DecodeFromBytes(data[20:20+e.MatchLen], df); err != nil {
			return err
		}
	}

	return nil
}
