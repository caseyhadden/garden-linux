// This file was generated by counterfeiter
package fakes

import (
	"net"
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/fences/netfence/network/subnets"
)

type FakeBridgedSubnets struct {
	AllocateStub        func(subnets.SubnetSelector, subnets.IPSelector) (*net.IPNet, net.IP, string, error)
	allocateMutex       sync.RWMutex
	allocateArgsForCall []struct {
		arg1 subnets.SubnetSelector
		arg2 subnets.IPSelector
	}
	allocateReturns struct {
		result1 *net.IPNet
		result2 net.IP
		result3 string
		result4 error
	}
	ReleaseStub        func(*net.IPNet, net.IP) (bool, string, error)
	releaseMutex       sync.RWMutex
	releaseArgsForCall []struct {
		arg1 *net.IPNet
		arg2 net.IP
	}
	releaseReturns struct {
		result1 bool
		result2 string
		result3 error
	}
	RecoverStub        func(*net.IPNet, net.IP, string) error
	recoverMutex       sync.RWMutex
	recoverArgsForCall []struct {
		arg1 *net.IPNet
		arg2 net.IP
		arg3 string
	}
	recoverReturns struct {
		result1 error
	}
	CapacityStub        func() int
	capacityMutex       sync.RWMutex
	capacityArgsForCall []struct{}
	capacityReturns struct {
		result1 int
	}
}

func (fake *FakeBridgedSubnets) Allocate(arg1 subnets.SubnetSelector, arg2 subnets.IPSelector) (*net.IPNet, net.IP, string, error) {
	fake.allocateMutex.Lock()
	fake.allocateArgsForCall = append(fake.allocateArgsForCall, struct {
		arg1 subnets.SubnetSelector
		arg2 subnets.IPSelector
	}{arg1, arg2})
	fake.allocateMutex.Unlock()
	if fake.AllocateStub != nil {
		return fake.AllocateStub(arg1, arg2)
	} else {
		return fake.allocateReturns.result1, fake.allocateReturns.result2, fake.allocateReturns.result3, fake.allocateReturns.result4
	}
}

func (fake *FakeBridgedSubnets) AllocateCallCount() int {
	fake.allocateMutex.RLock()
	defer fake.allocateMutex.RUnlock()
	return len(fake.allocateArgsForCall)
}

func (fake *FakeBridgedSubnets) AllocateArgsForCall(i int) (subnets.SubnetSelector, subnets.IPSelector) {
	fake.allocateMutex.RLock()
	defer fake.allocateMutex.RUnlock()
	return fake.allocateArgsForCall[i].arg1, fake.allocateArgsForCall[i].arg2
}

func (fake *FakeBridgedSubnets) AllocateReturns(result1 *net.IPNet, result2 net.IP, result3 string, result4 error) {
	fake.AllocateStub = nil
	fake.allocateReturns = struct {
		result1 *net.IPNet
		result2 net.IP
		result3 string
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeBridgedSubnets) Release(arg1 *net.IPNet, arg2 net.IP) (bool, string, error) {
	fake.releaseMutex.Lock()
	fake.releaseArgsForCall = append(fake.releaseArgsForCall, struct {
		arg1 *net.IPNet
		arg2 net.IP
	}{arg1, arg2})
	fake.releaseMutex.Unlock()
	if fake.ReleaseStub != nil {
		return fake.ReleaseStub(arg1, arg2)
	} else {
		return fake.releaseReturns.result1, fake.releaseReturns.result2, fake.releaseReturns.result3
	}
}

func (fake *FakeBridgedSubnets) ReleaseCallCount() int {
	fake.releaseMutex.RLock()
	defer fake.releaseMutex.RUnlock()
	return len(fake.releaseArgsForCall)
}

func (fake *FakeBridgedSubnets) ReleaseArgsForCall(i int) (*net.IPNet, net.IP) {
	fake.releaseMutex.RLock()
	defer fake.releaseMutex.RUnlock()
	return fake.releaseArgsForCall[i].arg1, fake.releaseArgsForCall[i].arg2
}

func (fake *FakeBridgedSubnets) ReleaseReturns(result1 bool, result2 string, result3 error) {
	fake.ReleaseStub = nil
	fake.releaseReturns = struct {
		result1 bool
		result2 string
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBridgedSubnets) Recover(arg1 *net.IPNet, arg2 net.IP, arg3 string) error {
	fake.recoverMutex.Lock()
	fake.recoverArgsForCall = append(fake.recoverArgsForCall, struct {
		arg1 *net.IPNet
		arg2 net.IP
		arg3 string
	}{arg1, arg2, arg3})
	fake.recoverMutex.Unlock()
	if fake.RecoverStub != nil {
		return fake.RecoverStub(arg1, arg2, arg3)
	} else {
		return fake.recoverReturns.result1
	}
}

func (fake *FakeBridgedSubnets) RecoverCallCount() int {
	fake.recoverMutex.RLock()
	defer fake.recoverMutex.RUnlock()
	return len(fake.recoverArgsForCall)
}

func (fake *FakeBridgedSubnets) RecoverArgsForCall(i int) (*net.IPNet, net.IP, string) {
	fake.recoverMutex.RLock()
	defer fake.recoverMutex.RUnlock()
	return fake.recoverArgsForCall[i].arg1, fake.recoverArgsForCall[i].arg2, fake.recoverArgsForCall[i].arg3
}

func (fake *FakeBridgedSubnets) RecoverReturns(result1 error) {
	fake.RecoverStub = nil
	fake.recoverReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBridgedSubnets) Capacity() int {
	fake.capacityMutex.Lock()
	fake.capacityArgsForCall = append(fake.capacityArgsForCall, struct{}{})
	fake.capacityMutex.Unlock()
	if fake.CapacityStub != nil {
		return fake.CapacityStub()
	} else {
		return fake.capacityReturns.result1
	}
}

func (fake *FakeBridgedSubnets) CapacityCallCount() int {
	fake.capacityMutex.RLock()
	defer fake.capacityMutex.RUnlock()
	return len(fake.capacityArgsForCall)
}

func (fake *FakeBridgedSubnets) CapacityReturns(result1 int) {
	fake.CapacityStub = nil
	fake.capacityReturns = struct {
		result1 int
	}{result1}
}

var _ subnets.BridgedSubnets = new(FakeBridgedSubnets)