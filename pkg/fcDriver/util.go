package fc

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
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
