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
	"github.com/skydive-project/goloxi/of12"
)

var OpenFlow12 OpenFlow12Protocol

type OpenFlow12Protocol struct {
}

func (p OpenFlow12Protocol) String() string {
	return "OpenFlow 1.2"
}

func (p OpenFlow12Protocol) GetVersion() uint8 {
	return goloxi.VERSION_1_2
}

func (p OpenFlow12Protocol) NewHello(versionBitmap uint32) goloxi.Message {
	return of12.NewHello()
}

func (p OpenFlow12Protocol) NewEchoRequest() goloxi.Message {
	return of12.NewEchoRequest()
}

func (p OpenFlow12Protocol) NewEchoReply() goloxi.Message {
	return of12.NewEchoReply()
}

func (p OpenFlow12Protocol) NewBarrierRequest() goloxi.Message {
	return of12.NewBarrierRequest()
}

func (p OpenFlow12Protocol) DecodeMessage(data []byte) (goloxi.Message, error) {
	return of12.DecodeMessage(data)
}
