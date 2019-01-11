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
	"github.com/skydive-project/goloxi/of10"
)

var OpenFlow10 OpenFlow10Protocol

type OpenFlow10Protocol struct {
}

func (p OpenFlow10Protocol) String() string {
	return "OpenFlow 1.0"
}

func (p OpenFlow10Protocol) GetVersion() uint8 {
	return goloxi.VERSION_1_0
}

func (p OpenFlow10Protocol) NewHello(versionBitmap uint32) goloxi.Message {
	return of10.NewHello()
}

func (p OpenFlow10Protocol) NewEchoRequest() goloxi.Message {
	return of10.NewEchoRequest()
}

func (p OpenFlow10Protocol) NewEchoReply() goloxi.Message {
	return of10.NewEchoReply()
}

func (p OpenFlow10Protocol) NewBarrierRequest() goloxi.Message {
	return of10.NewBarrierRequest()
}

func (p OpenFlow10Protocol) DecodeMessage(data []byte) (goloxi.Message, error) {
	return of10.DecodeMessage(data)
}
