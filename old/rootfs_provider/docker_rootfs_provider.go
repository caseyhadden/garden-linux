package rootfs_provider

import (
	"errors"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/garden-linux/old/repository_fetcher"
	"github.com/cloudfoundry-incubator/garden-linux/process"
)

type dockerRootFSProvider struct {
	graphDriver   graphdriver.Driver
	volumeCreator VolumeCreator
	repoFetcher   repository_fetcher.RepositoryFetcher
	namespacer    Namespacer
	copier        Copier
	clock         clock.Clock

	fallback RootFSProvider
}

var ErrInvalidDockerURL = errors.New("invalid docker url")

//go:generate counterfeiter -o fake_graph_driver/fake_graph_driver.go . GraphDriver
type GraphDriver interface {
	graphdriver.Driver
}

//go:generate counterfeiter -o fake_copier/fake_copier.go . Copier
type Copier interface {
	Copy(src, dest string) error
}

func NewDocker(
	repoFetcher repository_fetcher.RepositoryFetcher,
	graphDriver GraphDriver,
	volumeCreator VolumeCreator,
	namespacer Namespacer,
	copier Copier,
	clock clock.Clock,
) (RootFSProvider, error) {
	return &dockerRootFSProvider{
		repoFetcher:   repoFetcher,
		graphDriver:   graphDriver,
		volumeCreator: volumeCreator,
		namespacer:    namespacer,
		copier:        copier,
		clock:         clock,
	}, nil
}

func (provider *dockerRootFSProvider) ProvideRootFS(logger lager.Logger, id string, url *url.URL, shouldNamespace bool) (string, process.Env, error) {
	if len(url.Path) == 0 {
		return "", nil, ErrInvalidDockerURL
	}

	tag := "latest"
	if len(url.Fragment) > 0 {
		tag = url.Fragment
	}

	imageID, envvars, volumes, err := provider.repoFetcher.Fetch(logger, url, tag)
	if err != nil {
		return "", nil, err
	}

	if shouldNamespace {
		if imageID, err = provider.namespace(imageID); err != nil {
			return "", nil, err
		}
	}

	err = provider.graphDriver.Create(id, imageID)
	if err != nil {
		return "", nil, err
	}

	rootPath, err := provider.graphDriver.Get(id, "")
	if err != nil {
		return "", nil, err
	}

	for _, v := range volumes {
		if err = provider.volumeCreator.Create(rootPath, v); err != nil {
			return "", nil, err
		}
	}

	return rootPath, envvars, nil
}

func (provider *dockerRootFSProvider) namespace(imageID string) (string, error) {
	namespacedImageID := imageID + "@namespaced"
	if !provider.graphDriver.Exists(namespacedImageID) {
		if err := provider.createNamespacedLayer(namespacedImageID, imageID); err != nil {
			return "", err
		}
	}

	return namespacedImageID, nil
}

func (provider *dockerRootFSProvider) createNamespacedLayer(id string, parentId string) error {
	var err error
	var path string
	if path, err = provider.createLayer(id, parentId); err != nil {
		return err
	}

	return provider.namespacer.Namespace(path)
}

func (provider *dockerRootFSProvider) createLayer(id, parentId string) (string, error) {
	errs := func(err error) (string, error) {
		return "", err
	}

	if err := provider.graphDriver.Create(id, parentId); err != nil {
		return errs(err)
	}

	namespacedRootfs, err := provider.graphDriver.Get(id, "")
	if err != nil {
		return errs(err)
	}

	return namespacedRootfs, nil
}

func (provider *dockerRootFSProvider) CleanupRootFS(logger lager.Logger, id string) error {
	provider.graphDriver.Put(id)

	var err error
	maxAttempts := 10

	for errorCount := 0; errorCount < maxAttempts; errorCount++ {
		err = provider.graphDriver.Remove(id)
		if err == nil {
			break
		}

		logger.Error("cleanup-rootfs", err, lager.Data{
			"current-attempts": errorCount + 1,
			"max-attempts":     maxAttempts,
		})

		provider.clock.Sleep(200 * time.Millisecond)
	}

	return err
}

type ShellOutCp struct {
	WorkDir string
}

func (s ShellOutCp) Copy(src, dest string) error {
	if err := os.Remove(dest); err != nil {
		return err
	}

	return exec.Command("cp", "-a", src, dest).Run()
}
