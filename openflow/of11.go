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
	"github.com/skydive-project/goloxi/of11"
)

var OpenFlow11 OpenFlow11Protocol

type OpenFlow11Protocol struct {
}

func (p OpenFlow11Protocol) String() string {
	return "OpenFlow 1.1"
}

func (p OpenFlow11Protocol) GetVersion() uint8 {
	return goloxi.VERSION_1_1
}

func (p OpenFlow11Protocol) NewHello(versionBitmap uint32) goloxi.Message {
	return of11.NewHello()
}

func (p OpenFlow11Protocol) NewEchoRequest() goloxi.Message {
	return of11.NewEchoRequest()
}

func (p OpenFlow11Protocol) NewEchoReply() goloxi.Message {
	return of11.NewEchoReply()
}

func (p OpenFlow11Protocol) NewBarrierRequest() goloxi.Message {
	return of11.NewBarrierRequest()
}

func (p OpenFlow11Protocol) DecodeMessage(data []byte) (goloxi.Message, error) {
	return of11.DecodeMessage(data)
}
