package fc

import (
	"github.com/j-griffith/csi-connectors/fibrechannel"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"fmt"
	"encoding/json"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/volume/util"
)

type fcDevice struct {
	connector *fibrechannel.Connector
	disk string
}

func getFCInfo(req *csi.NodePublishVolumeRequest) (*fcDevice, error) {
	volName := req.GetVolumeId()
	lun := req.GetVolumeAttributes()["lun"]
	targetWWNs := req.GetVolumeAttributes()["targetWWNs"]
	wwids := req.GetVolumeAttributes()["WWIDs"]

	if lun == "" || (targetWWNs == "" && wwids == "") {
		return nil, fmt.Errorf("FC target information is missing")
	}

	targetList := []string{}
	wwidList := []string{}

	if err := json.Unmarshal([]byte(targetWWNs), &targetList); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(wwids), &wwidList); err != nil {
		return nil, err
	}

	fcConnector := &fibrechannel.Connector{
		VolumeName: volName,
		TargetWWNs: targetList,
		WWIDs: wwidList,
		Lun: lun,
	}

	//Only pass the connector
	return &fcDevice{
		connector: fcConnector,
	}, nil

}

func getFCDiskMounter(req *csi.NodePublishVolumeRequest) *fibrechannel.FCMounter {
	readOnly := req.GetReadonly()
	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	mountOptions := req.GetVolumeCapability().GetMount().GetMountFlags()

	return &fibrechannel.FCMounter{
		fsType:       fsType,
		readOnly:     readOnly,
		mountOptions: mountOptions,
		mounter:      &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()},
		exec:         mount.NewOsExec(),
		targetPath:   req.GetTargetPath(),
		deviceUtil:   util.NewDeviceHandler(util.NewIOHandler()),
	}
}