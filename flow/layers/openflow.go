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
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/skydive-project/skydive/flow/layers/ofp10"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

// LayerTypeOpenflow registers a layer type for OpenFlow
var LayerTypeOpenflow = gopacket.RegisterLayerType(55557, gopacket.LayerTypeMetadata{Name: "LayerTypeOpenFlow", Decoder: gopacket.DecodeFunc(decodeOpenFlowLayer)})

type OpenFlowLayer struct {
	layers.BaseLayer
	openflow.Header
	openflow.Message
}

func (l *OpenFlowLayer) LayerType() gopacket.LayerType {
	return LayerTypeOpenflow
}

func (l *OpenFlowLayer) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	l.BaseLayer = layers.BaseLayer{Contents: data[:len(data)]}

	if err := l.Header.DecodeFromBytes(data, df); err != nil {
		return err
	}

	var msg openflow.Message
	var err error

	switch l.Version {
	case openflow.OPENFLOW_1_0:
		msg, err = ofp10.DecodeMessage(l.Type, data[8:], df)
	default:
		// Only supports OpenFlow 1.0 for now...
	}

	if err != nil {
		return err
	}

	l.Message = msg
	return nil
}

func (l *OpenFlowLayer) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	if err := l.Header.SerializeTo(b, opts); err != nil {
		return err
	}

	if err := l.Message.SerializeTo(b, opts); err != nil {
		return err
	}

	data := b.Bytes()
	binary.BigEndian.PutUint16(data[2:4], uint16(len(data)))
	l.Length = uint16(len(data))
	l.BaseLayer.Contents = data

	return nil
}

func decodeOpenFlowLayer(data []byte, p gopacket.PacketBuilder) error {
	l := &OpenFlowLayer{}

	err := l.DecodeFromBytes(data, p)
	if err != nil {
		return err
	}

	p.AddLayer(l)
	return nil
}
