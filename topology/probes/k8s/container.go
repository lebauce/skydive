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

package k8s

import (
	"github.com/skydive-project/skydive/filters"
	"github.com/skydive-project/skydive/probe"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/graph"

	"k8s.io/api/core/v1"
)

const (
	dockerContainerNameField = "Docker.Labels.io.kubernetes.container.name"
	dockerPodNameField       = "Docker.Labels.io.kubernetes.pod.name"
	dockerPodNamespaceField  = "Docker.Labels.io.kubernetes.pod.namespace"
)

type containerProbe struct {
	graph.EventHandler
	graph.DefaultGraphListener
	podCache *ResourceCache
	graph    *graph.Graph
}

func (p *containerProbe) getContainerMetadata(pod *v1.Pod, container *v1.Container) graph.Metadata {
	m := NewMetadata(Manager, "container", container, container.Name, pod.Namespace)
	m.SetField("Pod", pod.Name)
	m.SetFieldAndNormalize("Labels", pod.Labels)
	m.SetField("Image", container.Image)
	return m
}

func (p *containerProbe) handleContainers(podNode *graph.Node) map[graph.Identifier]*graph.Node {
	pod := p.podCache.getByNode(podNode).(*v1.Pod)
	containers := make(map[graph.Identifier]*graph.Node)
	for _, container := range pod.Spec.Containers {
		uid := graph.GenID(string(pod.GetUID()), container.Name)
		m := p.getContainerMetadata(pod, &container)
		node := p.graph.GetNode(uid)
		if node == nil {
			node = p.graph.NewNode(uid, m)
			p.graph.NewEdge(graph.GenID(), podNode, node, topology.OwnershipMetadata(), "")
		} else {
			p.graph.SetMetadata(node, m)
		}
		containers[uid] = node
	}
	return containers
}

func (p *containerProbe) OnNodeAdded(node *graph.Node) {
	p.handleContainers(node)
}

func (p *containerProbe) OnNodeUpdated(node *graph.Node) {
	previousContainers := p.graph.LookupChildren(node, graph.Metadata{"Type": "container"}, topology.OwnershipMetadata())
	containers := p.handleContainers(node)

	for _, container := range previousContainers {
		if _, found := containers[container.ID]; !found {
			p.graph.DelNode(container)
		}
	}
}

func (p *containerProbe) OnNodeDeleted(node *graph.Node) {
	containers := p.graph.LookupChildren(node, graph.Metadata{"Type": "container"}, topology.OwnershipMetadata())
	for _, container := range containers {
		p.graph.DelNode(container)
	}
}

func (p *containerProbe) Start() {
	p.podCache.AddEventListener(p)
}

func (p *containerProbe) Stop() {
	p.podCache.RemoveEventListener(p)
}

func newContainerProbe(podProbe Subprobe, g *graph.Graph) Subprobe {
	return &containerProbe{graph: g, podCache: podProbe.(*ResourceCache)}
}

func newDockerIndexer(g *graph.Graph) *graph.MetadataIndexer {
	m := graph.NewElementFilter(filters.NewAndFilter(
		filters.NewTermStringFilter("Manager", "docker"),
		filters.NewTermStringFilter("Type", "container"),
		filters.NewNotNullFilter(dockerPodNamespaceField),
		filters.NewNotNullFilter(dockerPodNameField),
	))

	return graph.NewMetadataIndexer(g, g, m, dockerPodNamespaceField, dockerPodNameField, dockerContainerNameField)
}

func newContainerLinker(g *graph.Graph, subprobes map[string]Subprobe) probe.Probe {
	podProbe := subprobes["pod"]
	if podProbe == nil {
		return nil
	}

	k8sIndexer := newObjectIndexer(g, g, "container", "Namespace", "Pod", "Name")
	k8sIndexer.Start()

	dockerIndexer := newDockerIndexer(g)
	dockerIndexer.Start()

	return graph.NewMetadataIndexerLinker(g, k8sIndexer, dockerIndexer, newEdgeMetadata())
}
