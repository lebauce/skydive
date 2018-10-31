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
	"fmt"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

const (
	// Switch configuration messages
	OFPT_FEATURES_REQUEST = iota + openflow.OFPT_VENDOR + 1
	OFPT_FEATURES_REPLY
	OFPT_GET_CONFIG_REQUEST
	OFPT_GET_CONFIG_REPLY
	OFPT_SET_CONFIG

	// Asynchronous messages
	OFPT_PACKET_IN
	OFPT_FLOW_REMOVED
	OFPT_PORT_STATUS

	// Controller command messages
	OFPT_PACKET_OUT
	OFPT_FLOW_MOD
	OFPT_PORT_MOD

	// Statistics messages
	OFPT_STATS_REQUEST
	OFPT_STATS_REPLY

	// Barrier messages
	OFPT_BARRIER_REQUEST
	OFPT_BARRIER_REPLY

	// Queue Configuration messages
	OFPT_QUEUE_GET_CONFIG_REQUEST
	OFPT_QUEUE_GET_CONFIG_REPLY
)

const OF_ALL_TABLES = 0xFF

type MsgFlowStatsRequest struct {
	Match   Match
	TableID uint8
	OutPort uint16
}

type Action struct {
}

// stopped at ofp_aggregate_stats_request, page 3

func DecodeMessage(msgType openflow.MessageType, data []byte, df gopacket.DecodeFeedback) (msg openflow.Message, err error) {
	switch msgType {
	case openflow.OFPT_HELLO, openflow.OFPT_ERROR,
		openflow.OFPT_ECHO_REQUEST, openflow.OFPT_VENDOR:
		return openflow.DecodeImmutableMessage(msgType, data, df)
	case OFPT_FEATURES_REQUEST:
		msg = &FeaturesRequestMsg{}
	case OFPT_FEATURES_REPLY:
		msg = &FeaturesReplyMsg{}
	case OFPT_GET_CONFIG_REQUEST:
	case OFPT_GET_CONFIG_REPLY:
	case OFPT_SET_CONFIG:
	case OFPT_PACKET_IN:
	case OFPT_FLOW_REMOVED:
	case OFPT_PORT_STATUS:
	case OFPT_PACKET_OUT:
	case OFPT_FLOW_MOD:
	case OFPT_PORT_MOD:
	case OFPT_STATS_REQUEST:
		msg = &StatsRequestMsg{}
	case OFPT_STATS_REPLY:
		msg = &StatsReplyMsg{}
	case OFPT_BARRIER_REQUEST:
		msg = &BarrierRequestMsg{}
	case OFPT_BARRIER_REPLY:
		msg = &BarrierReplyMsg{}
	case OFPT_QUEUE_GET_CONFIG_REQUEST:
		// msg = &QueueGetConfigRequest{}
	case OFPT_QUEUE_GET_CONFIG_REPLY:
		// msg = &QueueGetConfigReply{}
	default:
		return nil, fmt.Errorf("Invalid OpenFlow 1.0 message type: %s", msgType)
	}

	if msg == nil {
		msg = &openflow.UnhandledMessage{Type: msgType}
	}

	if err := msg.DecodeFromBytes(data, df); err != nil {
		return nil, err
	}

	return msg, err
}
