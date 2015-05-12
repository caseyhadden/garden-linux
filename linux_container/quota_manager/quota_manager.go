package quota_manager

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/garden-linux/old/logging"
	"github.com/cloudfoundry/gunk/command_runner"
	"github.com/pivotal-golang/lager"
)

type BtrfsQuotaManager struct {
	enabled    bool
	mountPoint string
	graphRoot  string
	runner     command_runner.CommandRunner
}

const QUOTA_BLOCK_SIZE = 1024

func New(runner command_runner.CommandRunner, mountPoint, graphRoot string) *BtrfsQuotaManager {
	return &BtrfsQuotaManager{
		enabled: true,

		graphRoot: graphRoot,
		runner:    runner,

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

	qgroupId, path, err := m.volumeInfo(logger, cid)

	if err != nil {
		return err
	}

	cmd := exec.Command("btrfs", "qgroup", "limit", fmt.Sprintf("%d", limits.ByteHard), fmt.Sprintf("0/%d", qgroupId), path)
	if err = runner.Run(cmd); err != nil {
		return fmt.Errorf("quota_manager: failed to apply limit: %v", err)
	}

	return nil
}

func (m *BtrfsQuotaManager) GetLimits(logger lager.Logger, cid string) (garden.DiskLimits, error) {
	var quotaOut bytes.Buffer
	var byteLimit uint64
	var err error

	if !m.enabled {
		return garden.DiskLimits{}, nil
	}

	runner := logging.Runner{
		Logger:        logger,
		CommandRunner: m.runner,
	}

	limits := garden.DiskLimits{}
	_, path, err := m.volumeInfo(logger, cid)
	if err != nil {
		return limits, err
	}

	quotaCmd := exec.Command("sh", "-c", fmt.Sprintf("btrfs qgroup show -rF --raw %s | tail -n 1 | awk '{ print $4 }'", path))
	quotaCmd.Stdout = &quotaOut

	if err = runner.Run(quotaCmd); err != nil {
		return limits, fmt.Errorf("quota_manager: failed to get limit: %s", err)
	}

	if byteLimit, err = strconv.ParseUint(strings.Trim(quotaOut.String(), "\n"), 10, 64); err != nil {
		return limits, fmt.Errorf("quota_manager: failed to parse result: %s", err)
	}

	limits.ByteHard = byteLimit
	limits.ByteSoft = byteLimit

	return limits, err
}

func (m *BtrfsQuotaManager) GetUsage(logger lager.Logger, cid string) (garden.ContainerDiskStat, error) {
	//func (m *BtrfsQuotaManager) GetUsage(logger lager.Logger, uid int) (garden.ContainerDiskStat, error) {
	// TODO properly move to cid.
	uid := 123

	if !m.enabled {
		return garden.ContainerDiskStat{}, nil
	}

	repquota := exec.Command("repquota", m.mountPoint, fmt.Sprintf("%d", uid))

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

func (m *BtrfsQuotaManager) volumeInfo(logger lager.Logger, cid string) (int, string, error) {
	runner := logging.Runner{
		Logger:        logger,
		CommandRunner: m.runner,
	}

	// graphpath!
	listCmd := exec.Command("btrfs", "subvolume", "list", m.graphRoot)
	var listOut bytes.Buffer
	listCmd.Stdout = &listOut

	if err := runner.Run(listCmd); err != nil {
		return -1, "", fmt.Errorf("quota_manager: failed to list subvolumes: %v", err)
	}

	var path string
	var qgroupId, skip int
	found := false

	var err error
	lines := strings.Split(listOut.String(), "\n")
	for _, line := range lines {
		var n int
		n, err = fmt.Sscanf(line, "ID %d gen %d top level %d path %s", &qgroupId, &skip, &skip, &path)

		if err != nil || n != 4 {
			break
		}

		if strings.Contains(path, cid) {
			found = true
			break
		}
	}

	if !found {
		return -1, "", fmt.Errorf("quota_manager: subvolume not found: %s", err)
	}

	volume := strings.Split(strings.Trim(m.graphRoot, "/"), "/")[0]
	path = filepath.Join("/", volume, path)

	return qgroupId, path, nil
}
