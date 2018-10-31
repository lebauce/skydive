/*
 * Copyright (C) 2018 Red Hat
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

package ovsdb

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/skydive-project/skydive/flow/layers/ofp10"
	"github.com/skydive-project/skydive/flow/layers/openflow"
)

type testMonitorListener struct {
	featuresChan chan openflow.Message
	statsChan    chan openflow.Message
}

func (l *testMonitorListener) OnMessage(msg openflow.Message) {
	switch msg.MessageType() {
	case ofp10.OFPT_FEATURES_REPLY:
		l.featuresChan <- msg
	case ofp10.OFPT_STATS_REPLY:
		l.statsChan <- msg
	}
}

func TestOVSMonitorConnect(t *testing.T) {
	monitor, err := NewMonitor("tcp:127.0.0.1:6633")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := monitor.Start(ctx); err != nil {
		t.Error(err)
	}

	monitor.sendMessage(&ofp10.FeaturesRequestMsg{})

	listener := &testMonitorListener{
		featuresChan: make(chan openflow.Message),
		statsChan:    make(chan openflow.Message),
	}
	monitor.RegisterListener(listener)

	select {
	case <-time.After(3 * time.Second):
		t.Error("No features reply message received")
	case features := <-listener.featuresChan:
		t.Logf("Features: %+v", features)
	}

	monitor.sendMessage(&ofp10.StatsRequestMsg{
		Type: ofp10.OFPST_VENDOR,
		Content: &ofp10.VendorStatsMsg{
			Vendor: ofp10.NX_VENDOR_ID,
			Content: &ofp10.NiciraStatsMsg{
				Subtype: ofp10.NXST_FLOW,
				Content: &ofp10.NxFlowMonitorRequest{
					Flags:   0x3f,
					OutPort: ofp10.OFPP_NONE,
					TableID: ofp10.OF_ALL_TABLES,
				},
			},
		},
	})

	select {
	case <-time.After(3 * time.Second):
		t.Error("No stats reply message received")
	case stats := <-listener.statsChan:
		if statsMessage, ok := stats.(*ofp10.StatsReplyMsg); ok {
			if vendorStatsMsg, ok := statsMessage.Content.(*ofp10.VendorStatsMsg); ok {
				if niciraStatsMsg, ok := vendorStatsMsg.Content.(*ofp10.NiciraStatsMsg); ok {
					t.Logf("Stats: %+v (%s)", niciraStatsMsg.Content, spew.Sdump(niciraStatsMsg))
				}
			}
		}
	}

	cancel()
}
