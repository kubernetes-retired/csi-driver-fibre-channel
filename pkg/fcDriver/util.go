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
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/kubernetes/pkg/util/mount"
	volumeutil "k8s.io/kubernetes/pkg/volume/util"
)

//FCMounter struct holds required parameters to mount a Fibre Channel Disk
type FCMounter struct {
	ReadOnly     bool
	FsType       string
	MountOptions []string
	Mounter      *mount.SafeFormatAndMount
	Exec         mount.Exec
	DeviceUtil   volumeutil.DeviceUtil
	TargetPath   string
}

func RunNodePublishServer(d *CSIDriver, ns csi.NodeServer, ids csi.IdentityServer) {
	s := NewNonBlockingGRPCServer()
	s.Start(d.endpoint, ids, nil, ns)
	s.Wait()
}

func NewVolumeCapabilityAccessMode(mode csi.VolumeCapability_AccessMode_Mode) *csi.VolumeCapability_AccessMode {
	return &csi.VolumeCapability_AccessMode{Mode: mode}
}

func MountDisk(mnter FCMounter, devicePath string) error {
	mntPath := mnter.TargetPath
	notMnt, err := mnter.Mounter.IsLikelyNotMountPoint(mntPath)

	if err != nil {
		return fmt.Errorf("Heuristic determination of mount point failed: %v", err)
	}

	if !notMnt {
		fmt.Printf("fc: %s already mounted", mnter.TargetPath)
	}

	if err = mnter.Mounter.FormatAndMount(devicePath, mnter.TargetPath, mnter.FsType, nil); err != nil {
		return fmt.Errorf("fc: failed to mount fc volume %s [%s] to %s, error %v", devicePath, mnter.FsType, mnter.TargetPath, err)
	}

	return nil
}
