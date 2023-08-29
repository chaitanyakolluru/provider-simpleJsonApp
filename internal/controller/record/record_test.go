/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package record

import (
	"context"
	"testing"

	"git.heb.com/provider-simplejsonapp/apis/records/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

// Unlike many Kubernetes projects Crossplane does not use third party testing
// libraries, per the common Go test review comments. Crossplane encourages the
// use of table driven unit tests. The tests of the crossplane-runtime project
// are representative of the testing style Crossplane encourages.
//
// https://github.com/golang/go/wiki/TestComments
// https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md#contributing-code

type fields struct {
	service *MockSjaClientInterface
}

type args struct {
	ctx context.Context
	mg  resource.Managed
}

var setupArgs = func(name, designation, location string, age int, todos []string) args {
	return args{
		ctx: context.TODO(),
		mg: &v1alpha1.Record{
			Spec: v1alpha1.RecordSpec{
				ForProvider: v1alpha1.RecordParameters{
					Name:        name,
					Age:         age,
					Designation: designation,
					Location:    location,
					Todos:       todos,
				},
			},
		},
	}
}

func returnMockObservationObject(mockName, mockLocation, mockDesignation string, mockAge int, mockTodos []string) v1alpha1.RecordObservation {
	return v1alpha1.RecordObservation{
		Name:        mockName,
		Age:         mockAge,
		Location:    mockLocation,
		Designation: mockDesignation,
		Todos:       mockTodos,
	}
}

func setupMocksForSjaClientInterface(t *testing.T, mockName, mockLocation, mockDesignation string, mockAge int, mockTodos []string, mockError error) *MockSjaClientInterface {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistryClient := NewMockSjaClientInterface(ctrl)
	mockRegistryClient.EXPECT().GetRecord(context.TODO(), gomock.Any()).AnyTimes().Return(returnMockObservationObject(mockName, mockLocation, mockDesignation, mockAge, mockTodos), mockError)
	mockRegistryClient.EXPECT().PostRecord(context.TODO(), gomock.Any()).AnyTimes().Return(returnMockObservationObject(mockName, mockLocation, mockDesignation, mockAge, mockTodos), mockError)
	mockRegistryClient.EXPECT().PutRecord(context.TODO(), gomock.Any()).AnyTimes().Return(returnMockObservationObject(mockName, mockLocation, mockDesignation, mockAge, mockTodos), mockError)
	mockRegistryClient.EXPECT().DeleteRecord(context.TODO(), gomock.Any()).AnyTimes().Return(returnMockObservationObject(mockName, mockLocation, mockDesignation, mockAge, mockTodos), mockError)

	return mockRegistryClient
}

// Observe() method test cases
func TestControllerObserve(t *testing.T) {

	type want struct {
		o   managed.ExternalObservation
		err error
	}

	type caseStructure struct {
		reason string
		fields fields
		args   args
		want   want
	}

	var setupWant = func(isExists, isUpdated bool, errorWant error) want {
		return want{
			o: managed.ExternalObservation{
				ResourceExists:    isExists,
				ResourceUpToDate:  isUpdated,
				ConnectionDetails: managed.ConnectionDetails{},
			},
			err: errorWant,
		}
	}

	type setupCaseStructureArgs struct {
		reason      string
		name        string
		age         int
		designation string
		location    string
		todos       []string

		mockName        string
		mockLocation    string
		mockDesignation string
		mockAge         int
		mockTodos       []string

		isExists  bool
		isUpdated bool
		errorGot  error
		errorWant error
	}

	var setupCaseStructure = func(t *testing.T, args setupCaseStructureArgs) caseStructure {
		return caseStructure{
			reason: args.reason,
			fields: fields{
				service: setupMocksForSjaClientInterface(t, args.mockName, args.mockLocation, args.mockDesignation, args.mockAge, args.mockTodos, args.errorGot),
			},
			args: setupArgs(args.name, args.designation, args.location, args.age, args.todos),
			want: setupWant(args.isExists, args.isUpdated, args.errorWant),
		}
	}

	cases := map[string]caseStructure{
		"reconciles as exists and updated": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as exists and updated",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1",
			mockAge:         0,
			mockDesignation: "desg",
			mockLocation:    "loc",
			mockTodos:       []string{"todo1", "todo2"},

			isExists:  true,
			isUpdated: true,
			errorGot:  nil,
			errorWant: nil,
		}),
		"reconciles as an not exists and not updated with an error when getrecord fails": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as not exists and not updated with an error",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1",
			mockAge:         0,
			mockDesignation: "desg",
			mockLocation:    "loc",
			mockTodos:       []string{"todo1", "todo2"},

			isExists:  false,
			isUpdated: false,
			errorGot:  errors.New("error"),
			errorWant: errors.New("error"),
		}),
		"reconciles as not exists and not updated when name from getrecord doesn't match ": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as not exists and not updated",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1-notsame",
			mockAge:         0,
			mockDesignation: "desg",
			mockLocation:    "loc",
			mockTodos:       []string{"todo1", "todo2"},

			isExists:  false,
			isUpdated: false,
			errorGot:  nil,
			errorWant: nil,
		}),
		"reconciles as exists and not updated when age from getrecord doesn't match ": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as exists and not updated",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1",
			mockAge:         1,
			mockDesignation: "desg",
			mockLocation:    "loc",
			mockTodos:       []string{"todo1", "todo2"},

			isExists:  true,
			isUpdated: false,
			errorGot:  nil,
			errorWant: nil,
		}),
		"reconciles as exists and not updated when location from getrecord doesn't match ": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as exists and not updated",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1",
			mockAge:         0,
			mockDesignation: "desg",
			mockLocation:    "loc-notsame",
			mockTodos:       []string{"todo1", "todo2"},

			isExists:  true,
			isUpdated: false,
			errorGot:  nil,
			errorWant: nil,
		}),
		"reconciles as exists and not updated when designation from getrecord doesn't match ": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as exists and not updated",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1",
			mockAge:         0,
			mockDesignation: "desg-notsame",
			mockLocation:    "loc",
			mockTodos:       []string{"todo1", "todo2"},

			isExists:  true,
			isUpdated: false,
			errorGot:  nil,
			errorWant: nil,
		}),
		"reconciles as exists and not updated when todos from getrecord doesn't match ": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not reconcile as exists and not updated",
			name:        "record1",
			age:         0,
			designation: "desg",
			location:    "loc",
			todos:       []string{"todo1", "todo2"},

			mockName:        "record1",
			mockAge:         0,
			mockDesignation: "desg",
			mockLocation:    "loc",
			mockTodos:       []string{"todo1", "todo2", "todo3"},

			isExists:  true,
			isUpdated: false,
			errorGot:  nil,
			errorWant: nil,
		}),
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{service: tc.fields.service}
			got, err := e.Observe(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, errors.Cause(err), test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}

func TestControllerCreate(t *testing.T) {
	type want struct {
		o   managed.ExternalCreation
		err error
	}
	type caseStructure struct {
		reason string
		fields fields
		args   args
		want   want
	}

	var setupWant = func(errorWant error) want {
		return want{
			o: managed.ExternalCreation{
				ConnectionDetails: managed.ConnectionDetails{},
			},
			err: errorWant,
		}
	}

	var setupWantNoConnectionDetails = func(errorWant error) want {
		return want{
			o:   managed.ExternalCreation{},
			err: errorWant,
		}
	}

	var setupTernaryAroundSetupWant = func(connDetails bool, errorWant error) want {
		if connDetails {
			return setupWant(errorWant)
		} else {
			return setupWantNoConnectionDetails(errorWant)
		}
	}

	type setupCaseStructureArgs struct {
		reason      string
		name        string
		mockName    string
		errorGot    error
		errorWant   error
		connDetails bool
	}

	var setupCaseStructure = func(t *testing.T, args setupCaseStructureArgs) caseStructure {
		return caseStructure{
			reason: args.reason,
			fields: fields{
				service: setupMocksForSjaClientInterface(t, args.mockName, "", "", 0, []string{}, args.errorGot),
			},
			args: setupArgs(args.name, "", "", 0, []string{}),
			want: setupTernaryAroundSetupWant(args.connDetails, args.errorWant),
		}
	}

	cases := map[string]caseStructure{
		"results as resource created": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not result as resource created",
			name:        "record1",
			mockName:    "record1",
			errorGot:    nil,
			errorWant:   nil,
			connDetails: true,
		}),
		"results an empty creation object with an error": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not return error",
			name:        "record1",
			mockName:    "record1",
			errorGot:    errors.New("create error"),
			errorWant:   errors.New("create error"),
			connDetails: false,
		}),
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{service: tc.fields.service}
			got, err := e.Create(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, errors.Cause(err), test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}

}

func TestControllerUpdate(t *testing.T) {
	type want struct {
		o   managed.ExternalUpdate
		err error
	}
	type caseStructure struct {
		reason string
		fields fields
		args   args
		want   want
	}

	var setupWant = func(errorWant error) want {
		return want{
			o: managed.ExternalUpdate{
				ConnectionDetails: managed.ConnectionDetails{},
			},
			err: errorWant,
		}
	}

	var setupWantNoConnectionDetails = func(errorWant error) want {
		return want{
			o:   managed.ExternalUpdate{},
			err: errorWant,
		}
	}

	var setupTernaryAroundSetupWant = func(connDetails bool, errorWant error) want {
		if connDetails {
			return setupWant(errorWant)
		} else {
			return setupWantNoConnectionDetails(errorWant)
		}
	}

	type setupCaseStructureArgs struct {
		reason      string
		name        string
		mockName    string
		errorGot    error
		errorWant   error
		connDetails bool
	}

	var setupCaseStructure = func(t *testing.T, args setupCaseStructureArgs) caseStructure {
		return caseStructure{
			reason: args.reason,
			fields: fields{
				service: setupMocksForSjaClientInterface(t, args.mockName, "", "", 0, []string{}, args.errorGot),
			},
			args: setupArgs(args.name, "", "", 0, []string{}),
			want: setupTernaryAroundSetupWant(args.connDetails, args.errorWant),
		}
	}

	cases := map[string]caseStructure{
		"results as resource updated": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not result as an updated resource",
			name:        "record1",
			mockName:    "record1",
			errorGot:    nil,
			errorWant:   nil,
			connDetails: true,
		}),
		"results an empty update object with an error": setupCaseStructure(t, setupCaseStructureArgs{
			reason:      "does not return error",
			name:        "record1",
			mockName:    "record1",
			errorGot:    errors.New("create error"),
			errorWant:   errors.New("create error"),
			connDetails: false,
		}),
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{service: tc.fields.service}
			got, err := e.Update(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, errors.Cause(err), test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}

}

func TestControllerDelete(t *testing.T) {
	type want struct {
		err error
	}
	type caseStructure struct {
		reason string
		fields fields
		args   args
		want   want
	}

	var setupWant = func(errorWant error) want {
		return want{
			err: errorWant,
		}
	}

	type setupCaseStructureArgs struct {
		reason    string
		name      string
		mockName  string
		errorGot  error
		errorWant error
	}

	var setupCaseStructure = func(t *testing.T, args setupCaseStructureArgs) caseStructure {
		return caseStructure{
			reason: args.reason,
			fields: fields{
				service: setupMocksForSjaClientInterface(t, args.mockName, "", "", 0, []string{}, args.errorGot),
			},
			args: setupArgs(args.name, "", "", 0, []string{}),
			want: setupWant(args.errorWant),
		}
	}

	cases := map[string]caseStructure{
		"results as resource deleted": setupCaseStructure(t, setupCaseStructureArgs{
			reason:    "does not result as a deleted resource",
			name:      "record1",
			mockName:  "record1",
			errorGot:  nil,
			errorWant: nil,
		}),
		"results an error": setupCaseStructure(t, setupCaseStructureArgs{
			reason:    "does not return error",
			name:      "record1",
			mockName:  "record1",
			errorGot:  errors.New("create error"),
			errorWant: errors.New("create error"),
		}),
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{service: tc.fields.service}
			_, err := e.Update(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, errors.Cause(err), test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
		})
	}

}
