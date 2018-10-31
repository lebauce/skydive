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
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

type StatsType uint16

const (
	OFPST_DESC StatsType = iota
	OPFST_FLOW
	OFPST_AGGREGATE
	OFPST_TABLE
	OFPST_PORT
	OFPST_QUEUE
	OFPST_VENDOR = 0xffff

	DESC_STR_LEN           = 256
	SERIAL_NUM_LEN         = 32
	OFP_MAX_TABLE_NAME_LEN = 32
)

type StatsRequestMsg struct {
	Type    StatsType
	Flags   uint16
	Content openflow.Serializable
}

func (m *StatsRequestMsg) MessageType() openflow.MessageType {
	return OFPT_STATS_REQUEST
}

func (m *StatsRequestMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(4)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], uint16(m.Type))
	binary.BigEndian.PutUint16(data[2:4], m.Flags)

	return m.Content.SerializeTo(b, opts)
}

func (m *StatsRequestMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 4 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Type = StatsType(binary.BigEndian.Uint16(data[0:2]))
	m.Flags = binary.BigEndian.Uint16(data[2:4])

	switch m.Type {
	case OFPST_DESC:
		m.Content = &DescStatsReply{}
	case OPFST_FLOW:
		m.Content = &FlowStatsReply{}
	case OFPST_AGGREGATE:
		m.Content = &AggregateStatsReply{}
	case OFPST_TABLE:
		m.Content = &TableStats{}
	case OFPST_PORT:
		m.Content = &PortStats{}
	case OFPST_QUEUE:
		m.Content = &QueueStats{}
	case OFPST_VENDOR:
		m.Content = &VendorStatsMsg{}
	default:
		return fmt.Errorf("Invalid stats message '%d'", m.Type)
	}

	return m.Content.DecodeFromBytes(data[4:], df)
}

type StatsReplyMsg struct {
	Type    StatsType
	Flags   uint16
	Content openflow.Serializable
}

func (m *StatsReplyMsg) MessageType() openflow.MessageType {
	return OFPT_STATS_REPLY
}

func (m *StatsReplyMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(4)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], uint16(m.Type))
	binary.BigEndian.PutUint16(data[2:4], m.Flags)

	return m.Content.SerializeTo(b, opts)
}

func (m *StatsReplyMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 4 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Type = StatsType(binary.BigEndian.Uint16(data[0:2]))
	m.Flags = binary.BigEndian.Uint16(data[2:4])

	switch m.Type {
	case OFPST_DESC:
		m.Content = &DescStatsReply{}
	case OPFST_FLOW:
		m.Content = &FlowStatsReply{}
	case OFPST_AGGREGATE:
		m.Content = &AggregateStatsReply{}
	case OFPST_TABLE:
		m.Content = &TableStats{}
	case OFPST_PORT:
		m.Content = &PortStats{}
	case OFPST_QUEUE:
		m.Content = &QueueStats{}
	case OFPST_VENDOR:
		m.Content = &VendorStatsMsg{}
	default:
		return fmt.Errorf("Invalid stats message '%d'", m.Type)
	}

	return m.Content.DecodeFromBytes(data[4:], df)
}

type VendorStatsMsg struct {
	Vendor  uint32
	Content openflow.Serializable
}

func (m *VendorStatsMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(4)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], m.Vendor)

	return m.Content.SerializeTo(b, opts)
}

func (m *VendorStatsMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 4 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Vendor = binary.BigEndian.Uint32(data[0:4])

	switch m.Vendor {
	case NX_VENDOR_ID:
		m.Content = &NiciraStatsMsg{}
	default:
		m.Content = &openflow.UnhandledContent{}
	}

	return m.Content.DecodeFromBytes(data[4:], df)
}

// FlowStatsReply is modeled after ofp_flow_stats
type FlowStatsReply struct {
	Length      uint16
	TableID     uint8
	Match       Match
	Duration    time.Duration
	Priority    uint16
	IdleTimeout time.Duration
	HardTimeout time.Duration
	Cookie      uint64
	PacketCount uint64
	ByteCount   uint64
	Actions     []Action
}

func (m *FlowStatsReply) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(4)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], uint16(m.Length))
	data[2] = m.TableID
	data[3] = 0

	if err := m.Match.SerializeTo(b, opts); err != nil {
		return err
	}

	if data, err = b.AppendBytes(44); err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], uint32(m.Duration/time.Second))
	binary.BigEndian.PutUint32(data[4:8], uint32(m.Duration%time.Second))
	binary.BigEndian.PutUint16(data[8:10], m.Priority)
	binary.BigEndian.PutUint16(data[10:12], uint16(m.IdleTimeout.Seconds()))
	binary.BigEndian.PutUint16(data[12:14], uint16(m.HardTimeout.Seconds()))
	copy(data[14:20], []byte{0, 0, 0, 0, 0, 0})
	binary.BigEndian.PutUint64(data[20:28], m.Cookie)
	binary.BigEndian.PutUint64(data[28:36], m.PacketCount)
	binary.BigEndian.PutUint64(data[36:44], m.ByteCount)

	return nil
}

func (m *FlowStatsReply) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Length = binary.BigEndian.Uint16(data[0:2])
	m.TableID = data[2]

	if err := m.Match.DecodeFromBytes(data[4:44], df); err != nil {
		return err
	}

	m.Duration = time.Second*time.Duration(binary.BigEndian.Uint32(data[44:48])) +
		time.Nanosecond*time.Duration(binary.BigEndian.Uint32(data[48:52]))
	m.Priority = binary.BigEndian.Uint16(data[52:54])

	m.IdleTimeout = time.Second * time.Duration(binary.BigEndian.Uint16(data[54:56]))
	m.HardTimeout = time.Second * time.Duration(binary.BigEndian.Uint16(data[56:58]))

	m.Cookie = binary.BigEndian.Uint64(data[64:72])
	m.PacketCount = binary.BigEndian.Uint64(data[72:80])
	m.ByteCount = binary.BigEndian.Uint64(data[80:88])

	return nil
}

// AggregateStatsRequest is modeled after ofp_aggregate_stats_request
type AggregateStatsRequest struct {
	Match   Match
	TableID uint8
	OutPort uint16
}

// AggregateStatsReply is modeled after ofp_aggregate_stats_reply
type AggregateStatsReply struct {
	openflow.UnhandledContent
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
}

type DescStatsReply struct {
	Manufacturer string
	Hardware     string
	Software     string
	SerialNumber string
	Datapath     string
}

func (m *DescStatsReply) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(1056)
	if err != nil {
		return err
	}

	copy(data[0:DESC_STR_LEN], m.Manufacturer)
	copy(data[DESC_STR_LEN:DESC_STR_LEN*2], m.Hardware)
	copy(data[DESC_STR_LEN*2:DESC_STR_LEN*3], m.Software)
	copy(data[DESC_STR_LEN*3:DESC_STR_LEN*3+SERIAL_NUM_LEN], m.SerialNumber)
	copy(data[DESC_STR_LEN*3+SERIAL_NUM_LEN:], m.Datapath)

	return nil
}

func (m *DescStatsReply) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 1056 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	if len(data) != 1056 {
		return errors.New("OFPST_DESC message should be 1056 bytes long")
	}

	m.Manufacturer = string(bytes.Trim(data[0:DESC_STR_LEN], "\x00"))
	m.Hardware = string(bytes.Trim(data[DESC_STR_LEN:DESC_STR_LEN*2], "\x00"))
	m.Software = string(bytes.Trim(data[DESC_STR_LEN*2:DESC_STR_LEN*3], "\x00"))
	m.SerialNumber = string(bytes.Trim(data[DESC_STR_LEN*3:DESC_STR_LEN*3+SERIAL_NUM_LEN], "\x00"))
	m.Datapath = string(bytes.Trim(data[DESC_STR_LEN*3+SERIAL_NUM_LEN:], "\x00"))

	return nil
}

// TableStats is modeled after ofp_table_stats
type TableStats struct {
	openflow.UnhandledContent
	TableID   uint8
	Name      string
	Wildcards uint32
}

// PortStatsRequest is modeled after ofp_port_stats_request
type PortStatsRequest struct {
	PortNumber uint16
}

type PortStats struct {
	openflow.UnhandledContent
	PortNumber uint16
	RxPackets  uint64
	TxPackets  uint64
	RxBytes    uint64
	TxBytes    uint64
	RxDropped  uint64
	TxDropped  uint64
	RxErrors   uint64
	TxErrors   uint64
	RxFrameErr uint64
	RxOverErr  uint64
	RxCrcErr   uint64
	Collisions uint64
}

type PortStatsReply struct {
	PortStats []PortStats
}

// QueueStatsRequest is modeled after ofp_queue_stats_request
type QueueStatsRequest struct {
	PortNumber uint16
	QueueID    uint32
}

// QueueStats is modeled after ofp_queue_stats
type QueueStats struct {
	openflow.UnhandledContent
	PortNumber uint16
	QueueID    uint32
	TxBytes    uint64
	TxPackets  uint64
	TxErrors   uint64
}
