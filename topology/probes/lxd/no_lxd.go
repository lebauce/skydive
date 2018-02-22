// +build !linux

package lxd

import (
	"github.com/skydive-project/skydive/common"
	ns "github.com/skydive-project/skydive/topology/probes/netns"
)

// LxdProbe describes a Lxd topology graph that enhance the graph
type LxdProbe struct{}

// Start the probe
func (u *LxdProbe) Start() {}

// Stop the probe
func (u *LxdProbe) Stop() {}

// NewLxdProbe creates a new topology Lxd probe
func NewLxdProbe(nsProbe *ns.NetNSProbe, lxdURL string) (*LxdProbe, error) {
	return nil, common.ErrNotImplemented
}
