package fc
import (
	"github.com/j-griffith/csi-connectors/fibrechannel"
	"golang.org/x/net/context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)
type fcNodeServer struct {
	Driver *CSIDriver
}
var fcDisk fcDevice

func (ns *fcNodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *fcNodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *fcNodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	fcDevice, err := getFCInfo(req)

	if err != nil {
		return nil, err
	}
	fcDisk = *fcDevice
	disk, connectError := fibrechannel.Connect(*fcDevice.connector)

	fcDisk.disk = disk
	fcmounter := getFCDiskMounter(req)
	fibrechannel.MountDisk(*fcmounter, disk)

	if connectError != nil {
		return nil, connectError
	}

	//Need to add mounting
	return nil, nil
}

func (ns *fcNodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {

	//Need to Add unmount

	if err := fibrechannel.Disconnect(*fcDisk.connector, fcDisk.disk); err != nil {
		return nil, err
	}

	return nil, nil
}