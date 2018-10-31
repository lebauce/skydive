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

	"github.com/google/gopacket"
)

type ErrorType uint16

const (
	OFPET_HELLO_FAILED ErrorType = iota
	OFPET_BAD_REQUEST
	OFPET_BAD_ACTION
	OFPET_FLOW_MOD_FAILED
	OFPET_PORT_MOD_FAILED
	OFPET_QUEUE_OP_FAILED
)

// ofp_hello_failed_code
const (
	OFPHFC_INCOMPATIBLE = iota
	OFPHFC_EPERM
)

// ofp_bad_request_code
const (
	OFPBRC_BAD_VERSION = iota
	OFPBRC_BAD_TYPE
	OFPBRC_BAD_STAT
	OFPBRC_BAD_VENDOR
	OFPBRC_BAD_SUBTYPE
	OFPBRC_EPERM
	OFPBRC_BAD_LEN
	OFPBRC_BUFFER_EMPTY
	OFPBRC_BUFFER_UNKNOWN
)

// ofp_bad_action_code
const (
	OFPBAC_BAD_TYPE = iota
	OFPBAC_BAD_LEN
	OFPBAC_BAD_VENDOR
	OFPBAC_BAD_VENDOR_TYPE
	OFPBAC_BAD_OUT_PORT
	OFPBAC_BAD_ARGUMENT
	OFPBAC_EPERM
	OFPBAC_TOO_MANY
	OFPBAC_BAD_QUEUE
)

// ofp_flow_mod_failed_code
const (
	OFPFMFC_ALL_TABLES_FULL = iota
	OFPFMFC_OVERLAP
	OFPFMFC_EPERM
	OFPFMFC_BAD_EMERG_TIMEOUT
	OFPFMFC_BAD_COMMAND
	OFPFMFC_UNSUPPORTED
)

// ofp_port_mod_failed_code
const (
	OFPPMFC_BAD_PORT = iota
	OFPPMFC_BAD_HW_ADDR
)

// ofp_queue_op_failed_code
const (
	OFPQOFC_BAD_PORT = iota
	OFPQOFC_BAD_QUEUE
	OFPQOFC_EPERM
)

type ErrorMsg struct {
	Type ErrorType
	Code uint16
}

func (m *ErrorMsg) MessageType() MessageType {
	return OFPT_ERROR
}

func (m *ErrorMsg) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	data, err := b.AppendBytes(4)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(data[0:2], uint16(m.Type))
	binary.BigEndian.PutUint16(data[2:4], m.Code)

	return nil
}

func (m *ErrorMsg) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 4 {
		df.SetTruncated()
		return errors.New("OpenFlow packet too short")
	}

	m.Type = ErrorType(binary.BigEndian.Uint16(data[0:2]))
	m.Code = binary.BigEndian.Uint16(data[2:4])

	return nil
}

func DecodeErrorMsg(data []byte, df gopacket.DecodeFeedback) (Message, error) {
	errorMsg := &ErrorMsg{}
	if err := errorMsg.DecodeFromBytes(data, df); err != nil {
		return nil, err
	}
	return errorMsg, nil
}
