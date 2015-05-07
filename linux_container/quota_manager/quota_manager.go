package quota_manager

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/garden-linux/old/logging"
	"github.com/cloudfoundry/gunk/command_runner"
	"github.com/pivotal-golang/lager"
)

type BtrfsQuotaManager struct {
	enabled bool

	binPath string
	runner  command_runner.CommandRunner

	mountPoint string
}

const QUOTA_BLOCK_SIZE = 1024

func New(runner command_runner.CommandRunner, mountPoint, binPath string) *BtrfsQuotaManager {
	return &BtrfsQuotaManager{
		enabled: true,

		binPath: binPath,
		runner:  runner,

		mountPoint: mountPoint,
	}
}

func (m *BtrfsQuotaManager) Disable() {
	m.enabled = false
}

func (m *BtrfsQuotaManager) SetLimits(logger lager.Logger, cid string, limits garden.DiskLimits) error {
	if !m.enabled {
		return nil
	}

	if limits.BlockSoft != 0 {
		limits.ByteSoft = limits.BlockSoft * QUOTA_BLOCK_SIZE
	}

	if limits.BlockHard != 0 {
		limits.ByteHard = limits.BlockHard * QUOTA_BLOCK_SIZE
	}

	runner := logging.Runner{
		Logger:        logger,
		CommandRunner: m.runner,
	}

	listCmdR, listCmdW, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("quota_manager: create OS pipe: %v", err)
	}
	defer listCmdR.Close()
	defer listCmdW.Close()

	listCmd := exec.Command("btrfs", "subvolume", "list", m.mountPoint)
	listCmd.Stdout = listCmdW

	if err = runner.Start(listCmd); err != nil {
		return fmt.Errorf("quota_manager: failed to list subvolumes: %v", err)
	}
	defer runner.Wait(listCmd)

	var path string
	var qgroupId, skip int
	found := false

	for {
		n, err := fmt.Fscanf(listCmdR, "ID %d gen %d top level %d path %s", &qgroupId, &skip, &skip, &path)
		if err != nil {
			return fmt.Errorf("quota_manager: failed to get subvolume qgroup id: %v", err)
		}
		if n != 4 {
			break
		}

		if strings.Contains(path, cid) {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("quota_manager: subvolume not found")
	}

	cmd := exec.Command("btrfs", "qgroup", "limit", fmt.Sprintf("%d", limits.ByteHard), fmt.Sprintf("0/%d", qgroupId), m.mountPoint)
	if err = runner.Run(cmd); err != nil {
		return fmt.Errorf("quota_manager: failed to apply limit: %v", err)
	}

	return nil
}

func (m *BtrfsQuotaManager) GetLimits(logger lager.Logger, uid int) (garden.DiskLimits, error) {
	if !m.enabled {
		return garden.DiskLimits{}, nil
	}

	repquota := exec.Command(path.Join(m.binPath, "repquota"), m.mountPoint, fmt.Sprintf("%d", uid))

	limits := garden.DiskLimits{}

	repR, repW, err := os.Pipe()
	if err != nil {
		return limits, err
	}

	defer repR.Close()
	defer repW.Close()

	repquota.Stdout = repW

	runner := logging.Runner{
		Logger:        logger,
		CommandRunner: m.runner,
	}

	err = runner.Start(repquota)
	if err != nil {
		return limits, err
	}

	defer runner.Wait(repquota)

	var skip uint32

	_, err = fmt.Fscanf(
		repR,
		"%d %d %d %d %d %d %d %d",
		&skip,
		&skip,
		&limits.BlockSoft,
		&limits.BlockHard,
		&skip,
		&skip,
		&limits.InodeSoft,
		&limits.InodeHard,
	)

	return limits, err
}

func (m *BtrfsQuotaManager) GetUsage(logger lager.Logger, uid int) (garden.ContainerDiskStat, error) {
	if !m.enabled {
		return garden.ContainerDiskStat{}, nil
	}

	repquota := exec.Command(path.Join(m.binPath, "repquota"), m.mountPoint, fmt.Sprintf("%d", uid))

	usage := garden.ContainerDiskStat{}

	out := new(bytes.Buffer)

	repquota.Stdout = out

	runner := logging.Runner{
		Logger:        logger,
		CommandRunner: m.runner,
	}

	err := runner.Run(repquota)
	if err != nil {
		return usage, err
	}

	var skip uint32

	_, err = fmt.Fscanf(
		out,
		"%d %d %d %d %d %d %d %d",
		&skip,
		&usage.BytesUsed,
		&skip,
		&skip,
		&skip,
		&usage.InodesUsed,
		&skip,
		&skip,
	)

	return usage, err
}

func (m *BtrfsQuotaManager) MountPoint() string {
	return m.mountPoint
}

func (m *BtrfsQuotaManager) IsEnabled() bool {
	return m.enabled
}
