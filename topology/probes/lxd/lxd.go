// +build linux

package lxd

import (
	"context"
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/logging"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/graph"
	ns "github.com/skydive-project/skydive/topology/probes/netns"
	"github.com/vishvananda/netns"
	"sync"
	"sync/atomic"
	"time"
)

type containerInfo struct {
	Pid  int
	Node *graph.Node
}

type LxdProbe struct {
	sync.RWMutex
	*ns.NetNSProbe
	state        int64
	wg           sync.WaitGroup
	connected    atomic.Value
	cancel       context.CancelFunc
	containerMap map[string]containerInfo
	hostNs       netns.NsHandle
	client       lxd.ContainerServer
}

func (probe *LxdProbe) containerNamespace(pid int) string {
	return fmt.Sprintf("/proc/%d/ns/net", pid)
}

func (probe *LxdProbe) registerContainer(id string) {
	probe.Lock()
	defer probe.Unlock()

	if _, ok := probe.containerMap[id]; ok {
		return
	}

	// state.Network[].HostName == host side interface
	// state.Pid for lookup of network namespace
	state, _, _ := probe.client.GetContainerState(id)

	if state.Status != "Running" {
		return
	}

	logging.GetLogger().Infof("Operating on %d", state.Pid)

	nsHandle, err := netns.GetFromPid(int(state.Pid))
	if err != nil {
		return
	}
	defer nsHandle.Close()

	namespace := probe.containerNamespace(int(state.Pid))

	var n *graph.Node
	if probe.hostNs.Equal(nsHandle) {
		n = probe.Root
	} else {
		n = probe.Register(namespace, id)

		probe.Graph.Lock()
		probe.Graph.AddMetadata(n, "Manager", "lxd")
		probe.Graph.Unlock()
	}

	probe.Graph.Lock()
	metadata := graph.Metadata{
		"Type":        "container",
		"Name":        id,
		"LXDSettings": map[string]interface{}{},
	}

	containerNode := probe.Graph.NewNode(graph.GenID(), metadata)
	topology.AddOwnershipLink(probe.Graph, n, containerNode, nil)
	probe.Graph.Unlock()

	probe.containerMap[id] = containerInfo{
		Pid:  int(state.Pid),
		Node: containerNode,
	}
}

func (probe *LxdProbe) unregisterContainer(id string) {
	state, _, _ := probe.client.GetContainerState(id)

	if state.Status == "Running" {
		return
	}

	probe.Lock()
	defer probe.Unlock()

	infos, ok := probe.containerMap[id]
	if !ok {
		return
	}

	probe.Graph.Lock()
	probe.Graph.DelNode(infos.Node)
	probe.Graph.Unlock()

	namespace := probe.containerNamespace(infos.Pid)
	probe.Unregister(namespace)

	delete(probe.containerMap, id)
}

func (probe *LxdProbe) connect() (err error) {
	if probe.hostNs, err = netns.Get(); err != nil {
		return err
	}
	defer probe.hostNs.Close()

	probe.wg.Add(1)
	defer probe.wg.Done()

	client, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		return err
	}
	probe.client = client

	containers, err := probe.client.GetContainers()
	if err != nil {
		return err
	}

	for _, n := range containers {
		probe.registerContainer(n.Name)
		probe.unregisterContainer(n.Name)
	}

	return nil
}

func (probe *LxdProbe) Start() {
	if !atomic.CompareAndSwapInt64(&probe.state, common.StoppedState, common.RunningState) {
		return
	}

	go func() {
		for {
			state := atomic.LoadInt64(&probe.state)
			if state == common.StoppingState || state == common.StoppedState {
				break
			}

			if probe.connect() != nil {
				logging.GetLogger().Debugf("Start")
				time.Sleep(1 * time.Second)
			}

			time.Sleep(30 * time.Second) // avoid polling too much
			probe.wg.Wait()
		}
	}()
}

func (probe *LxdProbe) Stop() {
	if !atomic.CompareAndSwapInt64(&probe.state, common.RunningState, common.StoppingState) {
		return
	}

	if probe.connected.Load() == true {
		probe.cancel()
		probe.wg.Wait()
	}

	atomic.StoreInt64(&probe.state, common.StoppedState)
}

func NewLxdProbe(nsProbe *ns.NetNSProbe, lxdURL string) (*LxdProbe, error) {
	probe := &LxdProbe{
		NetNSProbe:   nsProbe,
		state:        common.StoppedState,
		containerMap: make(map[string]containerInfo),
	}
	logging.GetLogger().Debugf("Probe created")

	return probe, nil
}
