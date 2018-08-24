package fibrechannel_kubernetes_csi_driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
)

type CSIDriver struct {
	name string
	nodeID string
	version string
	csCap     []*csi.ControllerServiceCapability
	vcCap     []*csi.VolumeCapability_AccessMode
}

func NewCSIDriver(name string, v string, nodeID string) *CSIDriver {
	if name == "" {
		glog.Errorf("Driver name missing")
		return nil
	}

	if nodeID == "" {
		glog.Errorf("NodeID missing")
		return nil
	}

	if len(v) == 0 {
		glog.Errorf("Version argument missing")
		return nil
	}

	driver := CSIDriver{
		name:    name,
		version: v,
		nodeID:  nodeID,
	}

	return &driver
}
