/*
 * Copyright (C) 2019 Red Hat, Inc.
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

package blockdev

import (
	"context"
	"encoding/json"
	"os/exec"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/graffiti/graph"
	"github.com/skydive-project/skydive/probe"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/probes"
	tp "github.com/skydive-project/skydive/topology/probes"
)

const blockdevGroupName string = "block devices"

// Types returned from lsblk JSON
const multipathType string = "mpath"
const lvmType string = "lvm"
const blockGroupType string = "cluster"

// Types for skydive nodes
const lvmNodeType string = "blockdevlvm"
const blockdevNodeType string = "blockdev"
const leafNodeType string = "blockdevleaf"

const managerType string = "blockdev"

// BlockDevice used for JSON parsing
// Prior to version 2.33 lsblk generated JSON that encodes all values
// as strings.  Version 2.33 changed how integers and booleans are
// encoded.  Using json.RawMessage allows the code deal with either
// encoding.
type BlockDevice struct {
	Children []BlockDevice `mapstructure:"children"`

	Alignment    int64  `mapstructure:"alignment"`
	DiscAln      int64  `mapstructure:"disc-aln"`
	DiscGran     string `mapstructure:"disc-gran"`
	DiscMax      string `mapstructure:"disc-max"`
	DiscZero     bool   `mapstructure:"disc-zero"`
	Fsavail      string `mapstructure:"fsavail"`
	Fssize       string `mapstructure:"fssize"`
	Fstype       string `mapstructure:"fstype"`
	FsusePercent string `mapstructure:"fsuse%"`
	Fsused       string `mapstructure:"fsused"`
	Group        string `mapstructure:"group"`
	Hctl         string `mapstructure:"hctl"`
	Hotplug      bool   `mapstructure:"hotplug"`
	Kname        string `mapstructure:"kname"`
	Label        string `mapstructure:"label"`
	LogSec       int64  `mapstructure:"log-sec"`
	MajMin       string `mapstructure:"maj:min"`
	MinIo        int64  `mapstructure:"min-io"`
	Mode         string `mapstructure:"mode"`
	Model        string `mapstructure:"model"`
	Mountpoint   string `mapstructure:"mountpoint"`
	Name         string `mapstructure:"name"`
	OptIo        int64  `mapstructure:"opt-io"`
	Owner        string `mapstructure:"owner"`
	Partflags    string `mapstructure:"partflags"`
	Partlabel    string `mapstructure:"partlabel"`
	Parttype     string `mapstructure:"parttype"`
	Partuuid     string `mapstructure:"partuuid"`
	Path         string `mapstructure:"path"`
	PhySec       int64  `mapstructure:"hpy-sec"`
	Pkname       string `mapstructure:"pkname"`
	Pttype       string `mapstructure:"pttype"`
	Ptuuid       string `mapstructure:"ptuuid"`
	Ra           int64  `mapstructure:"ra"`
	Rand         bool   `mapstructure:"rand"`
	Rev          string `mapstructure:"rev"`
	Rm           bool   `mapstructure:"rm"`
	Ro           bool   `mapstructure:"ro"`
	Rota         bool   `mapstructure:"rota"`
	RqSize       int64  `mapstructure:"rq-size"`
	Sched        string `mapstructure:"sched"`
	Serial       string `mapstructure:"serial"`
	Size         string `mapstructure:"size"`
	State        string `mapstructure:"state"`
	Subsystems   string `mapstructure:"subsystems"`
	Tran         string `mapstructure:"tran"`
	Type         string `mapstructure:"type"`
	UUID         string `mapstructure:"uuid"`
	Vendor       string `mapstructure:"vendor"`
	Wsame        string `mapstructure:"wsame"`
	WWN          string `mapstructure:"wwn"`
}

// Devices used for JSON parsing
type Devices struct {
	Blockdevices []BlockDevice `mapstructure:"blockdevices"`
}

type blockdevInfo struct {
	ID   string
	Node *graph.Node
}

// ProbeHandler describes a block device graph that enhances the graph
type ProbeHandler struct {
	common.RWMutex
	Ctx         tp.Context
	blockdevMap map[string]blockdevInfo
	wg          sync.WaitGroup
	Groups      map[string]*graph.Node
}

func (p *ProbeHandler) getPath(blockdev BlockDevice) string {
	if blockdev.Path != "" {
		return blockdev.Path
	}
	if blockdev.Name != "" {
		return blockdev.Name
	}
	return blockdev.Kname
}

func (p *ProbeHandler) getID(blockdev BlockDevice) string {
	if blockdev.WWN != "" {
		return blockdev.WWN
	}
	if blockdev.Serial != "" {
		return blockdev.Serial
	}
	return blockdev.Name
}

func (p *ProbeHandler) getName(blockdev BlockDevice) string {
	if blockdev.Mountpoint != "" {
		return blockdev.Mountpoint
	}
	return blockdev.Name
}

func (p *ProbeHandler) addGroupByName(name string, WWN string) *graph.Node {
	if p.Groups[name] != nil {
		return p.Groups[name]
	}
	g := p.Ctx.Graph
	g.Lock()
	defer g.Unlock()
	metadata := graph.Metadata{
		"Name":    name,
		"WWN":     WWN,
		"Type":    blockGroupType,
		"Manager": managerType,
	}
	groupNode, err := g.NewNode(graph.GenID(), metadata)
	if err != nil {
		p.Ctx.Logger.Error(err)
		return nil
	}
	p.Groups[name] = groupNode

	return groupNode
}

// addGroup adds a group to the graph
func (p *ProbeHandler) addGroup(blockdev BlockDevice, WWN string) *graph.Node {
	groupName := p.getID(blockdev)
	return p.addGroupByName(groupName, WWN)
}

func (p *ProbeHandler) getMetaData(blockdev BlockDevice, childCount int, parentWWN string) (metadata graph.Metadata) {
	var blockdevMetadata Metadata
	var nodeType string

	// The nodeType maps to the icon is used for display
	if childCount > 0 {
		if blockdev.Type == lvmType {
			nodeType = lvmNodeType
		} else {
			nodeType = blockdevNodeType
		}
	} else {
		nodeType = leafNodeType
	}

	// JSON for multipath devices doesn't include the WWN, so get it from
	// the parent for a multipathType
	var WWN string
	if blockdev.Type == multipathType {
		WWN = parentWWN
	} else {
		WWN = blockdev.WWN
	}

	blockdevMetadata = Metadata{
		Index:        p.getID(blockdev),
		Name:         p.getName(blockdev),
		Alignment:    blockdev.Alignment,
		DiscAln:      blockdev.DiscAln,
		DiscGran:     blockdev.DiscGran,
		DiscMax:      blockdev.DiscMax,
		DiscZero:     blockdev.DiscZero,
		Fsavail:      blockdev.Fsavail,
		Fssize:       blockdev.Fssize,
		Fstype:       blockdev.Fstype,
		FsusePercent: blockdev.FsusePercent,
		Fsused:       blockdev.Fsused,
		Group:        blockdev.Group,
		Hctl:         blockdev.Hctl,
		Hotplug:      blockdev.Hotplug,
		Kname:        blockdev.Kname,
		Label:        blockdev.Label,
		LogSec:       blockdev.LogSec,
		MajMin:       blockdev.MajMin,
		MinIo:        blockdev.MinIo,
		Mode:         blockdev.Mode,
		Model:        blockdev.Model,
		Mountpoint:   blockdev.Mountpoint,
		OptIo:        blockdev.OptIo,
		Owner:        blockdev.Owner,
		Partflags:    blockdev.Partflags,
		Partlabel:    blockdev.Partlabel,
		Parttype:     blockdev.Parttype,
		Partuuid:     blockdev.Partuuid,
		Path:         blockdev.Path,
		PhySec:       blockdev.PhySec,
		Pttype:       blockdev.Pttype,
		Ptuuid:       blockdev.Ptuuid,
		Ra:           blockdev.Ra,
		Rand:         blockdev.Rand,
		Rev:          blockdev.Rev,
		Rm:           blockdev.Rm,
		Ro:           blockdev.Ro,
		Rota:         blockdev.Rota,
		RqSize:       blockdev.RqSize,
		Sched:        blockdev.Sched,
		Serial:       blockdev.Serial,
		Size:         blockdev.Size,
		State:        blockdev.State,
		Subsystems:   blockdev.Subsystems,
		Tran:         blockdev.Tran,
		Type:         blockdev.Type,
		UUID:         blockdev.UUID,
		Vendor:       blockdev.Vendor,
		Wsame:        blockdev.Wsame,
		WWN:          WWN,
	}
	metadata = graph.Metadata{
		"Path":     p.getPath(blockdev),
		"Type":     nodeType,
		"Name":     p.getName(blockdev),
		"Manager":  managerType,
		"BlockDev": blockdevMetadata,
	}

	return metadata
}

func (p *ProbeHandler) findBlockDev(path string, id string, newDevInfo []BlockDevice) bool {
	for i := range newDevInfo {
		if id == p.getID(newDevInfo[i]) && path == p.getPath(newDevInfo[i]) {
			return true
		}
		if p.findBlockDev(path, id, newDevInfo[i].Children) == true {
			return true
		}
	}
	return false
}

func (p *ProbeHandler) deleteIfRemoved(currentDevInfo blockdevInfo, newDevInfo []BlockDevice) {
	blockID, err := currentDevInfo.Node.GetFieldString("BlockID")
	if err != nil {
		return
	}
	path, err := currentDevInfo.Node.GetFieldString("Path")
	if err != nil {
		return
	}
	if p.findBlockDev(path, blockID, newDevInfo) == false {
		p.unregisterBlockdev(path)
	}
}

func (p *ProbeHandler) registerBlockdev(blockdev BlockDevice, parentWWN string) *graph.Node {
	var groupNode *graph.Node
	childCount := len(blockdev.Children)
	if childCount > 0 {
		for i := range blockdev.Children {
			groupNode = p.registerBlockdev(blockdev.Children[i], blockdev.WWN)
		}
		if blockdev.Type == multipathType {
			groupNode = p.addGroup(blockdev, parentWWN)
		}
	} else {
		groupNode = p.addGroup(blockdev, parentWWN)
	}

	p.Lock()
	defer p.Unlock()
	var graphMetaData graph.Metadata

	if _, ok := p.blockdevMap[blockdev.Name]; ok {
		return groupNode
	}

	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()

	node := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"Path": p.getPath(blockdev)})

	if node == nil {
		graphMetaData = p.getMetaData(blockdev, childCount, parentWWN)
		var err error
		if node, err = p.Ctx.Graph.NewNode(graph.GenID(), graphMetaData); err != nil {
			p.Ctx.Logger.Error(err)
			return nil
		}

		// If there is a WWN - check to see if there are other paths to the same block device.
		if blockdev.WWN != "" {
			multiPathNodes := p.Ctx.Graph.GetNodes(graph.Metadata{"BlockID": p.getID(blockdev)})

			for _, multiPathNode := range multiPathNodes {
				topology.AddLink(p.Ctx.Graph, multiPathNode, node, "multipath", nil)
			}
		}

		topology.AddOwnershipLink(p.Ctx.Graph, groupNode, node, nil)

		// Link physical disks and DVD/rom devices to the blockdevGroupName
		if blockdev.Type == "disk" || blockdev.Type == "rom" {
			topology.AddLink(p.Ctx.Graph, p.Groups[blockdevGroupName], node, "connected", nil)
		}

	}

	for i := range blockdev.Children {
		childNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"Name": p.getName(blockdev.Children[i])})
		if childNode != nil {
			topology.AddLink(p.Ctx.Graph, childNode, node, "connected", nil)
		}
	}

	p.blockdevMap[blockdev.Name] = blockdevInfo{
		ID:   blockdev.Name,
		Node: node,
	}
	return groupNode
}

func (p *ProbeHandler) unregisterBlockdev(id string) {
	p.Lock()
	defer p.Unlock()

	info, ok := p.blockdevMap[id]
	if !ok {
		return
	}

	p.Ctx.Graph.Lock()
	if err := p.Ctx.Graph.DelNode(info.Node); err != nil {
		p.Ctx.Graph.Unlock()
		p.Ctx.Logger.Error(err)
		return
	}
	p.Ctx.Graph.Unlock()

	delete(p.blockdevMap, id)
}

func (p *ProbeHandler) connect() error {
	var (
		cmdOut []byte
		err    error
		intf   interface{}
		result Devices
	)

	lsblkPath := p.Ctx.Config.GetString("agent.topology.blockdev.lsblk_path")
	if cmdOut, err = exec.Command(lsblkPath, "-pO", "--json").Output(); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(cmdOut[:]), &intf); err != nil {
		return err
	}

	if err = mapstructure.WeakDecode(intf, &result); err != nil {
		return err
	}

	// loop through the devices in the current map to make sure they haven't
	// been removed
	for _, current := range p.blockdevMap {
		p.deleteIfRemoved(current, result.Blockdevices)
	}

	for i := range result.Blockdevices {
		p.registerBlockdev(result.Blockdevices[i], "")
	}

	return nil
}

// Do adds a group for block devices, then issues a lsblk command and parsers the
// JSON output as a basis for the blockdev links.
func (p *ProbeHandler) Do(ctx context.Context, wg *sync.WaitGroup) error {
	p.addGroupByName(blockdevGroupName, "")

	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()

	if !topology.HaveLink(p.Ctx.Graph, p.Ctx.RootNode, p.Groups[blockdevGroupName], "connected") {
		if _, err := topology.AddLink(p.Ctx.Graph, p.Ctx.RootNode, p.Groups[blockdevGroupName], "connected", nil); err != nil {
			p.Ctx.Logger.Error(err)
			return err
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		after := time.After(10 * time.Minute)
		for {
			if err := p.connect(); err != nil {
				p.Ctx.Logger.Error(err)
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-after:
			}
		}
	}()
	return nil
}

// NewProbe initializes a new topology blockdev probe
func NewProbe(ctx tp.Context, bundle *probe.Bundle) (probe.Handler, error) {

	p := &ProbeHandler{
		blockdevMap: make(map[string]blockdevInfo),
		Groups:      make(map[string]*graph.Node),
		Ctx:         ctx,
	}

	p.Ctx = ctx
	p.blockdevMap = make(map[string]blockdevInfo)
	p.Groups = make(map[string]*graph.Node)

	return probes.NewProbeWrapper(p), nil
}

// Register registers graph metadata decoders
func Register() {
	graph.NodeMetadataDecoders["BlockDev"] = MetadataDecoder
}
