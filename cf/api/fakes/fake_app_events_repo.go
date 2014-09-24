// This file was generated by counterfeiter
package fakes

import (
	"github.com/cloudfoundry/cli/cf/models"
	"sync"
)

type FakeAppEventsRepository struct {
	RecentEventsStub        func(appGuid string, limit int64) ([]models.EventFields, error)
	recentEventsMutex       sync.RWMutex
	recentEventsArgsForCall []struct {
		arg1 string
		arg2 int64
	}
	recentEventsReturns struct {
		result1 []models.EventFields
		result2 error
	}
}

func (fake *FakeAppEventsRepository) RecentEvents(arg1 string, arg2 int64) ([]models.EventFields, error) {
	fake.recentEventsMutex.Lock()
	defer fake.recentEventsMutex.Unlock()
	fake.recentEventsArgsForCall = append(fake.recentEventsArgsForCall, struct {
		arg1 string
		arg2 int64
	}{arg1, arg2})
	if fake.RecentEventsStub != nil {
		return fake.RecentEventsStub(arg1, arg2)
	} else {
		return fake.recentEventsReturns.result1, fake.recentEventsReturns.result2
	}
}

func (fake *FakeAppEventsRepository) RecentEventsCallCount() int {
	fake.recentEventsMutex.RLock()
	defer fake.recentEventsMutex.RUnlock()
	return len(fake.recentEventsArgsForCall)
}

func (fake *FakeAppEventsRepository) RecentEventsArgsForCall(i int) (string, int64) {
	fake.recentEventsMutex.RLock()
	defer fake.recentEventsMutex.RUnlock()
	return fake.recentEventsArgsForCall[i].arg1, fake.recentEventsArgsForCall[i].arg2
}

func (fake *FakeAppEventsRepository) RecentEventsReturns(result1 []models.EventFields, result2 error) {
	fake.recentEventsReturns = struct {
		result1 []models.EventFields
		result2 error
	}{result1, result2}
}
