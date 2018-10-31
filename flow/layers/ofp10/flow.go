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
	"net"
	"time"
)

// FlowModCommand is modeled after ofp_flow_mod_command
type FlowModCommand uint16

const (
	OFPFC_ADD FlowModCommand = iota
	OFPFC_MODIFY
	OFPFC_MODIFY_STRICT
	OFPFC_DELETE
	OFPFC_DELETE_STRICT
)

// ofp_flow_mod_flags
const (
	OFPFF_SEND_FLOW_REM = 1 << 0
	OFPFF_CHECK_OVERLAP = 1 << 1
	OFPFF_EMERG         = 1 << 2
)

// FlowMod is modeled after ofp_flow_mod
type FlowMod struct {
	Match       Match
	Cookie      uint64
	Command     uint16
	IdleTimeout uint16
	HardTimeout uint16
	Priority    uint16
	Buffer      uint32
	OutPort     uint16
	Flags       uint16
	Actions     []Action
}

// PortMod is modeled after ofp_port_mod
type PortMod struct {
	PortNumber uint16
	HWAddr     net.HardwareAddr
	Config     ConfigFlags
	Mask       ConfigFlags
	Advertise  ConfigFlags
}

// FlowRemovedReason is modeled after ofp_flow_removed_reason
type FlowRemovedReason uint8

const (
	OFPRR_IDLE_TIMEOUT FlowRemovedReason = iota
	OFPRR_HARD_TIMEOUT
	OFPRR_DELETE
)

// FlowRemoved is modeled after ofp_flow_removed
type FlowRemoved struct {
	Match       Match
	Cookie      uint64
	Priority    uint16
	Reason      FlowRemovedReason
	Duration    time.Duration
	IdleTimeout time.Duration
	PacketCount uint64
	ByteCount   uint64
}
