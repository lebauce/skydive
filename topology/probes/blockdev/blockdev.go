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
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

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
	Children []BlockDevice `json:"children"`

	Alignment    json.RawMessage `json:"alignment"` //int64
	DiscAln      json.RawMessage `json:"disc-aln"`  // int64
	DiscGran     string          `json:"disc-gran"`
	DiscMax      string          `json:"disc-max"`
	DiscZero     json.RawMessage `json:"disc-zero"` // bool
	Fsavail      string          `json:"fsavail"`
	Fssize       string          `json:"fssize"`
	Fstype       string          `json:"fstype"`
	FsusePercent string          `json:"fsuse%"`
	Fsused       string          `json:"fsused"`
	Group        string          `json:"group"`
	Hctl         string          `json:"hctl"`
	Hotplug      json.RawMessage `json:"hotplug"` // bool
	Kname        string          `json:"kname"`
	Label        string          `json:"label"`
	LogSec       json.RawMessage `json:"log-sec"` // int64
	MajMin       string          `json:"maj:min"`
	MinIo        json.RawMessage `json:"min-io"` // int64
	Mode         string          `json:"mode"`
	Model        string          `json:"model"`
	Mountpoint   string          `json:"mountpoint"`
	Name         string          `json:"name"`
	OptIo        json.RawMessage `json:"opt-io"` // int64
	Owner        string          `json:"owner"`
	Partflags    string          `json:"partflags"`
	Partlabel    string          `json:"partlabel"`
	Parttype     string          `json:"parttype"`
	Partuuid     string          `json:"partuuid"`
	Path         string          `json:"path"`
	PhySec       json.RawMessage `json:"hpy-sec"` // int64
	Pkname       string          `json:"pkname"`
	Pttype       string          `json:"pttype"`
	Ptuuid       string          `json:"ptuuid"`
	Ra           json.RawMessage `json:"ra"`   // int64
	Rand         json.RawMessage `json:"rand"` // bool
	Rev          string          `json:"rev"`
	Rm           json.RawMessage `json:"rm"`      // bool
	Ro           json.RawMessage `json:"ro"`      // bool
	Rota         json.RawMessage `json:"rota"`    // bool
	RqSize       json.RawMessage `json:"rq-size"` // int64
	Sched        string          `json:"sched"`
	Serial       string          `json:"serial"`
	Size         string          `json:"size"`
	State        string          `json:"state"`
	Subsystems   string          `json:"subsystems"`
	Tran         string          `json:"tran"`
	Type         string          `json:"type"`
	UUID         string          `json:"uuid"`
	Vendor       string          `json:"vendor"`
	Wsame        string          `json:"wsame"`
	WWN          string          `json:"wwn"`
}

// Devices used for JSON parsing
type Devices struct {
	Blockdevices []BlockDevice `json:"blockdevices"`
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
		"Name": name,
		"WWN":  WWN,
		"Type": blockGroupType,
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

func (p *ProbeHandler) getInt(field json.RawMessage) int64 {
	// If the JSON has an integer value of "field": null field parameter
	// to this function will be set to nil
	if field == nil {
		return 0
	}
	if utf8.Valid(field) {
		var i int64
		i, err := strconv.ParseInt(strings.Replace(string(field), "\"", "", -1), 10, 64)

		if err == nil {
			return i
		}
	}
	p.Ctx.Logger.Error(errors.New("blockdev probe failed to parse integer JSON field from lsblk"))
	return -1
}

func (p *ProbeHandler) getBool(field json.RawMessage) bool {

	if utf8.Valid(field) {
		// ParseBool will deal with:
		// 1, t, T, true, TRUE, True along with the variations for false
		//
		// If the JSON is encoded as a string it will include the quotes.  Use
		// strings.Replace() to remove them
		b, err := strconv.ParseBool(strings.Replace(string(field), "\"", "", -1))
		if err != nil {
			p.Ctx.Logger.Error(err)
		}
		return b
	}
	p.Ctx.Logger.Error(errors.New("blockdev probe failed to parse boolean JSON field from lsblk"))
	return false
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
		Alignment:    p.getInt(blockdev.Alignment),
		DiscAln:      p.getInt(blockdev.DiscAln),
		DiscGran:     blockdev.DiscGran,
		DiscMax:      blockdev.DiscMax,
		DiscZero:     p.getBool(blockdev.DiscZero),
		Fsavail:      blockdev.Fsavail,
		Fssize:       blockdev.Fssize,
		Fstype:       blockdev.Fstype,
		FsusePercent: blockdev.FsusePercent,
		Fsused:       blockdev.Fsused,
		Group:        blockdev.Group,
		Hctl:         blockdev.Hctl,
		Hotplug:      p.getBool(blockdev.Hotplug),
		Kname:        blockdev.Kname,
		Label:        blockdev.Label,
		LogSec:       p.getInt(blockdev.LogSec),
		MajMin:       blockdev.MajMin,
		MinIo:        p.getInt(blockdev.MinIo),
		Mode:         blockdev.Mode,
		Model:        blockdev.Model,
		Mountpoint:   blockdev.Mountpoint,
		OptIo:        p.getInt(blockdev.OptIo),
		Owner:        blockdev.Owner,
		Partflags:    blockdev.Partflags,
		Partlabel:    blockdev.Partlabel,
		Parttype:     blockdev.Parttype,
		Partuuid:     blockdev.Partuuid,
		Path:         blockdev.Path,
		PhySec:       p.getInt(blockdev.PhySec),
		Pttype:       blockdev.Pttype,
		Ptuuid:       blockdev.Ptuuid,
		Ra:           p.getInt(blockdev.Ra),
		Rand:         p.getBool(blockdev.Rand),
		Rev:          blockdev.Rev,
		Rm:           p.getBool(blockdev.Rm),
		Ro:           p.getBool(blockdev.Ro),
		Rota:         p.getBool(blockdev.Rota),
		RqSize:       p.getInt(blockdev.RqSize),
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
		result Devices
	)

	testFile := os.Getenv("SKYDIVE_BLOCKDEV_TEST_FILE")

	if testFile != "" {
		if cmdOut, err = exec.Command("cat", testFile).Output(); err != nil {
			p.Ctx.Logger.Error(err)
			os.Exit(1)
		}
	} else {
		if cmdOut, err = exec.Command("lsblk", "-pO", "--json").Output(); err != nil {
			p.Ctx.Logger.Error(err)
			os.Exit(1)
		}
	}

	if err = json.Unmarshal([]byte(cmdOut[:]), &result); err != nil {
		p.Ctx.Logger.Error(err)
		os.Exit(1)
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

	if _, err := topology.AddLink(p.Ctx.Graph, p.Ctx.RootNode, p.Groups[blockdevGroupName], "connected", nil); err != nil {
		p.Ctx.Logger.Error(err)
		return err
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
