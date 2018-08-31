package fc

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
)

type CSIDriver struct {
	name     string
	version  string
	nodeID   string
	endpoint string

	csCap []*csi.ControllerServiceCapability
	vcCap []*csi.VolumeCapability_AccessMode
}

const (
	driverVersion = "0.1.0"
	driverName    = "fibrechannel"
)

func NewDriver(nodeID string, endpoint string) *CSIDriver {
	glog.Infof("Driver: %v nodeID: %v endpoint: %v", driverName, nodeID, endpoint)

	if nodeID == "" {
		glog.Errorf("NodeID missing")
		return nil
	}

	if endpoint == "" {
		glog.Errorf("endpoint missing")
		return nil
	}

	driver := CSIDriver{
		name:     driverName,
		version:  driverVersion,
		nodeID:   nodeID,
		endpoint: endpoint,
	}

	driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	return &driver
}

func (d *CSIDriver) AddVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) []*csi.VolumeCapability_AccessMode {
	var vca []*csi.VolumeCapability_AccessMode

	for _, c := range vc {
		glog.Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}

	d.vcCap = vca
	return vca
}

func NewNodeServer(d *CSIDriver) *fcNodeServer {
	return &fcNodeServer{
		Driver: d,
	}
}

func NewIdentityServer(d *CSIDriver) *FcIdentityServer {
	return &FcIdentityServer{
		Driver: d,
	}
}

func (d *CSIDriver) Run() {
	RunNodePublishServer(d, NewNodeServer(d), NewIdentityServer(d))
}
