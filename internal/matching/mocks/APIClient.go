// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	apihelper "github.com/s3f4/locationmatcher/pkg/apihelper"

	mock "github.com/stretchr/testify/mock"

	models "github.com/s3f4/locationmatcher/internal/matching/models"
)

// APIClient is an autogenerated mock type for the APIClient type
type APIClient struct {
	mock.Mock
}

// FindNearest provides a mock function with given fields: url, query
func (_m *APIClient) FindNearest(url string, query *models.Query) (*apihelper.Response, error) {
	ret := _m.Called(url, query)

	var r0 *apihelper.Response
	if rf, ok := ret.Get(0).(func(string, *models.Query) *apihelper.Response); ok {
		r0 = rf(url, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apihelper.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *models.Query) error); ok {
		r1 = rf(url, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}