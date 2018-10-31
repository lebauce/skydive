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

package layers

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/google/gopacket"
	"github.com/skydive-project/skydive/flow/layers/ofp10"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

func TestHelloMsg(t *testing.T) {
	b := gopacket.NewSerializeBuffer()
	ofLayer := &OpenFlowLayer{Version: openflow.OPENFLOW_1_0, Type: ofp10.OFPT_HELLO, Xid: 1, Message: ofp10.NewHelloMsg([]byte{})}
	if err := ofLayer.SerializeTo(b, gopacket.SerializeOptions{}); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(b.Bytes(), []byte{0x01, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01}) {
		t.Errorf("Wrong content: %+v", b.Bytes())
	}

	ofLayer2 := &OpenFlowLayer{}
	if err := ofLayer2.DecodeFromBytes(b.Bytes(), gopacket.NilDecodeFeedback); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ofLayer, ofLayer2) {
		t.Errorf("Both layer should be equal: %+v != %+v", ofLayer, ofLayer2)
	}
}

func TestFeatureRequest(t *testing.T) {
	// Request
	// 0105000800000003

	// Reply
	reply := []byte{
		0x01, 0x06, 0x00, 0x50, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x8a, 0xdd, 0x71, 0x39, 0xc8, 0x40,
		0x00, 0x00, 0x00, 0x00, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x0f, 0xff,
		0xff, 0xfe, 0x8a, 0xdd, 0x71, 0x39, 0xc8, 0x40, 0x62, 0x72, 0x2d, 0x69, 0x6e, 0x74, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	ofLayer := &OpenFlowLayer{}
	if err := ofLayer.DecodeFromBytes(reply, gopacket.NilDecodeFeedback); err != nil {
		t.Error(err)
	}

	msg := &ofp10.FeaturesReplyMsg{
		DatapathID:      152683692017728,
		BufferedPackets: 0,
		TableCount:      254,
		Capabilities:    199,
		Actions:         4095,
		Ports: []ofp10.PhysicalPort{
			ofp10.PhysicalPort{
				PortNumber: 65534,
				MAC:        net.HardwareAddr([]byte{0x8a, 0xdd, 0x71, 0x39, 0xc8, 0x40}),
				Name:       "br-int",
				Config:     1,
				State:      1,
				Current:    0,
				Advertised: 0,
				Supported:  0,
				Peer:       0,
			},
		},
	}

	if !reflect.DeepEqual(ofLayer.Message, msg) {
		t.Errorf("Both message should be equal: %+v != %+v", ofLayer.Message, msg)
	}
}
