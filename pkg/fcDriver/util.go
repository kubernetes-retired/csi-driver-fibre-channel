package fc

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
)

func RunNodePublishServer(d *CSIDriver, ns csi.NodeServer, ids csi.IdentityServer) {
	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.endpoint, ids, nil, ns)
	s.Wait()
}

func NewVolumeCapabilityAccessMode(mode csi.VolumeCapability_AccessMode_Mode) *csi.VolumeCapability_AccessMode {
	return &csi.VolumeCapability_AccessMode{Mode: mode}
}
