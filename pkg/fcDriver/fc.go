/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package fc

import (
	"encoding/json"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-fc/fibrechannel"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/volume/util"
)

type fcDevice struct {
	connector *fibrechannel.Connector
	disk      string
}

func getFCInfo(req *csi.NodePublishVolumeRequest) (*fcDevice, error) {
	volName := req.GetVolumeId()
	lun := req.GetVolumeContext()["lun"]
	targetWWNs := req.GetVolumeContext()["targetWWNs"]
	wwids := req.GetVolumeContext()["WWIDs"]

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
		WWIDs:      wwidList,
		Lun:        lun,
	}

	//Only pass the connector
	return &fcDevice{
		connector: fcConnector,
	}, nil

}

func getFCDiskMounter(req *csi.NodePublishVolumeRequest) FCMounter {
	readOnly := req.GetReadonly()
	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	mountOptions := req.GetVolumeCapability().GetMount().GetMountFlags()
	return FCMounter{
		ReadOnly:     readOnly,
		FsType:       fsType,
		MountOptions: mountOptions,
		Mounter:      &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()},
		Exec:         mount.NewOsExec(),
		DeviceUtil:   util.NewDeviceHandler(util.NewIOHandler()),
		TargetPath:   req.GetTargetPath(),
	}
}
