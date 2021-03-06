// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Circuit is an autogenerated mock type for the Circuit type
type Circuit struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, url, reader
func (_m *Circuit) Execute(ctx context.Context, url string, reader io.Reader) (*http.Response, error) {
	ret := _m.Called(ctx, url, reader)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, string, io.Reader) *http.Response); ok {
		r0 = rf(ctx, url, reader)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, io.Reader) error); ok {
		r1 = rf(ctx, url, reader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
