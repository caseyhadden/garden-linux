package containerizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cloudfoundry/gunk/command_runner"
)

var timeout = time.Second * 3

//go:generate counterfeiter -o fake_container_execer/FakeContainerExecer.go . ContainerExecer
type ContainerExecer interface {
	Exec(binPath string, args ...string) (int, error)
}

//go:generate counterfeiter -o fake_rootfs_enterer/FakeRootFSEnterer.go . RootFSEnterer
type RootFSEnterer interface {
	Enter() error
}

//go:generate counterfeiter -o fake_container_initializer/FakeContainerInitializer.go . ContainerInitializer
type ContainerInitializer interface {
	Init() error
}

//go:generate counterfeiter -o fake_container_daemon/FakeContainerDaemon.go . ContainerDaemon
type ContainerDaemon interface {
	Init() error
	Run() error
}

//go:generate counterfeiter -o fake_signaller/FakeSignaller.go . Signaller
type Signaller interface {
	SignalError(err error) error
	SignalSuccess() error
}

//go:generate counterfeiter -o fake_waiter/FakeWaiter.go . Waiter
type Waiter interface {
	Wait(timeout time.Duration) error
	IsSignalError(err error) bool
}

type Containerizer struct {
	InitBinPath string
	InitArgs    []string
	Execer      ContainerExecer
	RootfsPath  string
	Initializer ContainerInitializer
	Daemon      ContainerDaemon
	Signaller   Signaller
	Waiter      Waiter
	// Temporary until we merge the hook scripts functionality in Golang
	CommandRunner command_runner.CommandRunner
	LibPath       string
}

func (c *Containerizer) Create() error {
	// Temporary until we merge the hook scripts functionality in Golang
	cmd := exec.Command(path.Join(c.LibPath, "hook"), "parent-before-clone")
	if err := c.CommandRunner.Run(cmd); err != nil {
		return fmt.Errorf("containerizer: run `parent-before-clone`: %s", err)
	}

	if err := setHardRlimits(); err != nil {
		return err
	}

	pid, err := c.Execer.Exec(c.InitBinPath, c.InitArgs...)
	if err != nil {
		return fmt.Errorf("containerizer: create container: %s", err)
	}

	// Temporary until we merge the hook scripts functionality in Golang
	err = os.Setenv("PID", strconv.Itoa(pid))
	if err != nil {
		return fmt.Errorf("containerizer: failed to set PID env var: %s", err)
	}

	cmd = exec.Command(path.Join(c.LibPath, "hook"), "parent-after-clone")
	if err := c.CommandRunner.Run(cmd); err != nil {
		return fmt.Errorf("containerizer: run `parent-after-clone`: %s", err)
	}

	pivotter := exec.Command(filepath.Join(c.LibPath, "pivotter"), "-rootfs", c.RootfsPath)
	pivotter.Env = append(pivotter.Env, fmt.Sprintf("TARGET_NS_PID=%d", pid))
	if err := c.CommandRunner.Run(pivotter); err != nil {
		return fmt.Errorf("containerizer: run pivotter: %s", err)
	}

	if err := c.Signaller.SignalSuccess(); err != nil {
		return fmt.Errorf("containerizer: send success singnal to the container: %s", err)
	}

	if err := c.Waiter.Wait(timeout); err != nil {
		return fmt.Errorf("containerizer: wait for container: %s", err)
	}

	return nil
}

func (c *Containerizer) Run() error {
	if err := c.Daemon.Init(); err != nil {
		return c.signalErrorf("containerizer: initialize daemon: %s", err)
	}

	if err := c.Waiter.Wait(timeout); err != nil {
		return c.signalErrorf("containerizer: wait for host: %s", err)
	}

	if err := c.Initializer.Init(); err != nil {
		return c.signalErrorf("containerizer: initializing the container: %s", err)
	}

	if err := c.Signaller.SignalSuccess(); err != nil {
		return c.signalErrorf("containerizer: signal host: %s", err)
	}

	if err := c.Daemon.Run(); err != nil {
		return c.signalErrorf("containerizer: run daemon: %s", err)
	}

	return nil
}

func (c *Containerizer) signalErrorf(format string, err error) error {
	err = fmt.Errorf(format, err)

	if signalErr := c.Signaller.SignalError(err); signalErr != nil {
		err = fmt.Errorf("containerizer: signal error: %s (while signalling %s)", signalErr, err)
	}
	return err
}

const RLIMIT_INFINITY = ^uint64(0)

type RLimitEntry struct {
	Id  int
	Max uint64
}

func setHardRlimits() error {
	maxNoFile, err := maxNrOpen()
	if err != nil {
		return err
	}

	rLimitsMap := map[string]*RLimitEntry{
		"cpu":        &RLimitEntry{Id: syscall.RLIMIT_CPU, Max: RLIMIT_INFINITY},
		"fsize":      &RLimitEntry{Id: syscall.RLIMIT_FSIZE, Max: RLIMIT_INFINITY},
		"data":       &RLimitEntry{Id: syscall.RLIMIT_DATA, Max: RLIMIT_INFINITY},
		"stack":      &RLimitEntry{Id: syscall.RLIMIT_STACK, Max: RLIMIT_INFINITY},
		"core":       &RLimitEntry{Id: syscall.RLIMIT_CORE, Max: RLIMIT_INFINITY},
		"rss":        &RLimitEntry{Id: 5, Max: RLIMIT_INFINITY},
		"nproc":      &RLimitEntry{Id: 6, Max: RLIMIT_INFINITY},
		"nofile":     &RLimitEntry{Id: syscall.RLIMIT_NOFILE, Max: maxNoFile},
		"memlock":    &RLimitEntry{Id: 8, Max: RLIMIT_INFINITY},
		"as":         &RLimitEntry{Id: syscall.RLIMIT_AS, Max: RLIMIT_INFINITY},
		"locks":      &RLimitEntry{Id: 10, Max: RLIMIT_INFINITY},
		"sigpending": &RLimitEntry{Id: 11, Max: RLIMIT_INFINITY},
		"msgqueue":   &RLimitEntry{Id: 12, Max: RLIMIT_INFINITY},
		"nice":       &RLimitEntry{Id: 13, Max: RLIMIT_INFINITY},
		"rtprio":     &RLimitEntry{Id: 14, Max: RLIMIT_INFINITY},
	}

	for label, entry := range rLimitsMap {
		if err := setHardRLimit(label, entry.Id, entry.Max); err != nil {
			return err
		}
	}

	return nil
}

func setHardRLimit(label string, rLimitId int, rLimitMax uint64) error {
	var rlimit syscall.Rlimit

	if err := syscall.Getrlimit(rLimitId, &rlimit); err != nil {
		return fmt.Errorf("containerizer: get system rlimit_%s: %s", label, err)
	}

	rlimit.Max = rLimitMax
	if err := syscall.Setrlimit(rLimitId, &rlimit); err != nil {
		return fmt.Errorf("containerizer: setting hard rlimit_%s: %s", label, err)
	}

	return nil
}

func maxNrOpen() (uint64, error) {
	contents, err := ioutil.ReadFile("/proc/sys/fs/nr_open")
	if err != nil {
		return 0, fmt.Errorf("containerizer: failed to read /proc/sys/fs/nr_open: %s", err)
	}

	contentStr := strings.TrimSpace(string(contents))
	maxFiles, err := strconv.ParseUint(contentStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("containerizer: failed to convert contents of /proc/sys/fs/nr_open: %s", err)
	}

	return maxFiles, nil
}
