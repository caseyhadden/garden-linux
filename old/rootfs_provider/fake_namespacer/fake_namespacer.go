// This file was generated by counterfeiter
package fake_namespacer

import (
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/old/rootfs_provider"
)

type FakeNamespacer struct {
	NamespaceStub        func(rootfsPath string) error
	namespaceMutex       sync.RWMutex
	namespaceArgsForCall []struct {
		rootfsPath string
	}
	namespaceReturns struct {
		result1 error
	}
}

func (fake *FakeNamespacer) Namespace(rootfsPath string) error {
	fake.namespaceMutex.Lock()
	fake.namespaceArgsForCall = append(fake.namespaceArgsForCall, struct {
		rootfsPath string
	}{rootfsPath})
	fake.namespaceMutex.Unlock()
	if fake.NamespaceStub != nil {
		return fake.NamespaceStub(rootfsPath)
	} else {
		return fake.namespaceReturns.result1
	}
}

func (fake *FakeNamespacer) NamespaceCallCount() int {
	fake.namespaceMutex.RLock()
	defer fake.namespaceMutex.RUnlock()
	return len(fake.namespaceArgsForCall)
}

func (fake *FakeNamespacer) NamespaceArgsForCall(i int) string {
	fake.namespaceMutex.RLock()
	defer fake.namespaceMutex.RUnlock()
	return fake.namespaceArgsForCall[i].rootfsPath
}

func (fake *FakeNamespacer) NamespaceReturns(result1 error) {
	fake.NamespaceStub = nil
	fake.namespaceReturns = struct {
		result1 error
	}{result1}
}

var _ rootfs_provider.Namespacer = new(FakeNamespacer)
