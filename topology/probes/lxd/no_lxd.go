// +build !linux

package lxd

import (
	"github.com/skydive-project/skydive/common"
	ns "github.com/skydive-project/skydive/topology/probes/netns"
)

type LxdProbe struct{}

func (u *LxdProbe) Start() {}
func (u *LxdProbe) Stop()  {}

func NewLxdProbe(nsProbe *ns.NetNSProbe, lxdURL string) (*LxdProbe, error) {
	return nil, common.ErrNotImplemented
}
