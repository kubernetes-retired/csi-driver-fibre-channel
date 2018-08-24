package fc

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"fmt"
)

type CSIDriver struct {
	name string
	nodeID string
	version string
	csCap     []*csi.ControllerServiceCapability
	vcCap     []*csi.VolumeCapability_AccessMode
}


func NewCSIDriver(name string, v string, nodeID string) *CSIDriver {
	glog.Infof("Driver: %v version: %v nodeID: %v", name, v, nodeID)

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

func (d *CSIDriver) ValidateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, cap := range d.csCap {
		if c == cap.GetRpc().GetType() {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf("%s", c))
}
