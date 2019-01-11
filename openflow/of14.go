/*
 * Copyright (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy ofthe License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specificlanguage governing permissions and
 * limitations under the License.
 *
 */

package openflow

import (
	"github.com/skydive-project/goloxi"
	"github.com/skydive-project/goloxi/of14"
)

var OpenFlow14 OpenFlow14Protocol

type OpenFlow14Protocol struct {
}

func (p OpenFlow14Protocol) String() string {
	return "OpenFlow 1.4"
}

func (p OpenFlow14Protocol) GetVersion() uint8 {
	return goloxi.VERSION_1_4
}

func (p OpenFlow14Protocol) NewHello(versionBitmap uint32) goloxi.Message {
	msg := of14.NewHello()
	elem := of14.NewHelloElemVersionbitmap()
	elem.Length = 8
	bitmap := of14.NewUint32()
	bitmap.Value = versionBitmap
	elem.Bitmaps = append(elem.Bitmaps, bitmap)
	msg.Elements = append(msg.Elements, elem)
	return msg
}

func (p OpenFlow14Protocol) NewEchoRequest() goloxi.Message {
	return of14.NewEchoRequest()
}

func (p OpenFlow14Protocol) NewEchoReply() goloxi.Message {
	return of14.NewEchoReply()
}

func (p OpenFlow14Protocol) NewBarrierRequest() goloxi.Message {
	return of14.NewBarrierRequest()
}

func (p OpenFlow14Protocol) DecodeMessage(data []byte) (goloxi.Message, error) {
	return of14.DecodeMessage(data)
}
