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
	"fmt"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

const (
	OFPQT_NONE = iota
	OFPQT_MIN_RATE
)

// PacketQueue is modeled after ofp_packet_queue
type PacketQueue struct {
	QueueID    uint32
	Length     uint16
	Properties []QueuePropHeader
}

func (q *PacketQueue) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(data[0:4], q.QueueID)
	binary.BigEndian.PutUint16(data[4:6], q.Length)
	binary.BigEndian.PutUint16(data[6:8], 0)

	for _, property := range q.Properties {
		if err := property.SerializeTo(b, opts); err != nil {
			return err
		}
	}

	return nil
}

func (q *PacketQueue) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	q.QueueID = binary.BigEndian.Uint32(data[0:4])
	q.Length = binary.BigEndian.Uint16(data[4:6])
	data = data[8:]

	for len(data) > 0 {
		var queuePropHeader QueuePropHeader
		if err := queuePropHeader.DecodeFromBytes(data, df); err != nil {
			return err
		}
		q.Properties = append(q.Properties, queuePropHeader)
		data = data[queuePropHeader.Length:]
	}

	return nil
}

// QueuePropHeader is modeled after ofp_queue_prop_header
type QueuePropHeader struct {
	Property uint16
	Length   uint16
	Content  openflow.Serializable
}

func (h *QueuePropHeader) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], h.Property)
	binary.BigEndian.PutUint16(data[2:4], h.Length)
	binary.BigEndian.PutUint32(data[4:8], 0)

	return h.Content.SerializeTo(b, opts)
}

func (h *QueuePropHeader) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	h.Property = binary.BigEndian.Uint16(data[0:2])
	h.Length = binary.BigEndian.Uint16(data[2:4])
	data = data[8:]

	switch h.Property {
	case OFPQT_MIN_RATE:
		h.Content = &QueuePropMinRate{}
		return h.Content.DecodeFromBytes(data, df)
	default:
		return fmt.Errorf("Invalid queue property type: '%d'", h.Property)
	}
}

// QueuePropMinRate is modeled after ofp_queue_prop_min_rate
type QueuePropMinRate struct {
	Rate uint16
}

func (r *QueuePropMinRate) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], r.Rate)
	copy(data[2:8], []byte{0, 0, 0, 0, 0, 0})

	return nil
}

func (r *QueuePropMinRate) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	r.Rate = binary.BigEndian.Uint16(data[0:2])

	return nil
}

// QueueGetConfigRequest is modeled after ofp_queue_get_config_request
type QueueGetConfigRequest struct {
	Port uint16
}

func (r *QueueGetConfigRequest) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(2)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], r.Port)

	return nil
}

func (r *QueueGetConfigRequest) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 2 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	r.Port = binary.BigEndian.Uint16(data[0:2])

	return nil
}

// QueueGetConfigReply is modeled after ofp_queue_get_config_reply
type QueueGetConfigReply struct {
	Port   uint16
	Queues []PacketQueue
}

func (r *QueueGetConfigReply) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], r.Port)
	copy(data[2:8], []byte{0, 0, 0, 0, 0, 0})

	for _, packetQueue := range r.Queues {
		if err := packetQueue.SerializeTo(b, opts); err != nil {
			return err
		}
	}
	return nil
}

func (r *QueueGetConfigReply) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 8 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	r.Port = binary.BigEndian.Uint16(data[0:2])
	data = data[8:]

	for len(data) > 0 {
		var packetQueue PacketQueue
		if err := packetQueue.DecodeFromBytes(data, df); err != nil {
			return err
		}
		data = data[packetQueue.Length:]
		r.Queues = append(r.Queues, packetQueue)
	}

	return nil
}
