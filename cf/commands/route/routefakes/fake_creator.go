// This file was generated by counterfeiter
package routefakes

import (
	"sync"

	"github.com/cloudfoundry/cli/cf/commands/route"
	"github.com/cloudfoundry/cli/cf/models"
)

type FakeCreator struct {
	CreateRouteStub        func(hostName string, path string, port int, randomPort bool, domain models.DomainFields, space models.SpaceFields) (route models.Route, apiErr error)
	createRouteMutex       sync.RWMutex
	createRouteArgsForCall []struct {
		hostName   string
		path       string
		port       int
		randomPort bool
		domain     models.DomainFields
		space      models.SpaceFields
	}
	createRouteReturns struct {
		result1 models.Route
		result2 error
	}
}

func (fake *FakeCreator) CreateRoute(hostName string, path string, port int, randomPort bool, domain models.DomainFields, space models.SpaceFields) (route models.Route, apiErr error) {
	fake.createRouteMutex.Lock()
	fake.createRouteArgsForCall = append(fake.createRouteArgsForCall, struct {
		hostName   string
		path       string
		port       int
		randomPort bool
		domain     models.DomainFields
		space      models.SpaceFields
	}{hostName, path, port, randomPort, domain, space})
	fake.createRouteMutex.Unlock()
	if fake.CreateRouteStub != nil {
		return fake.CreateRouteStub(hostName, path, port, randomPort, domain, space)
	} else {
		return fake.createRouteReturns.result1, fake.createRouteReturns.result2
	}
}

func (fake *FakeCreator) CreateRouteCallCount() int {
	fake.createRouteMutex.RLock()
	defer fake.createRouteMutex.RUnlock()
	return len(fake.createRouteArgsForCall)
}

func (fake *FakeCreator) CreateRouteArgsForCall(i int) (string, string, int, bool, models.DomainFields, models.SpaceFields) {
	fake.createRouteMutex.RLock()
	defer fake.createRouteMutex.RUnlock()
	return fake.createRouteArgsForCall[i].hostName, fake.createRouteArgsForCall[i].path, fake.createRouteArgsForCall[i].port, fake.createRouteArgsForCall[i].randomPort, fake.createRouteArgsForCall[i].domain, fake.createRouteArgsForCall[i].space
}

func (fake *FakeCreator) CreateRouteReturns(result1 models.Route, result2 error) {
	fake.CreateRouteStub = nil
	fake.createRouteReturns = struct {
		result1 models.Route
		result2 error
	}{result1, result2}
}

var _ route.Creator = new(FakeCreator)
