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

	"github.com/chaitanyakolluru/provider-simplejsonapp/apis/records/v1alpha1"
	"github.com/chaitanyakolluru/provider-simplejsonapp/internal/controller/record/sjaclient"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
)

// Unlike many Kubernetes projects Crossplane does not use third party testing
// libraries, per the common Go test review comments. Crossplane encourages the
// use of table driven unit tests. The tests of the crossplane-runtime project
// are representative of the testing style Crossplane encourages.
//
// https://github.com/golang/go/wiki/TestComments
// https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md#contributing-code

type fields struct {
	service *sjaclient.SjaClient
}

type args struct {
	ctx context.Context
	mg  resource.Managed
}

type want struct {
	o   managed.ExternalObservation
	err error
}

type caseStructure map[string]struct {
	reason string
	fields fields
	args   args
	want   want
}

var setupArgs = func(testRecord v1alpha1.RecordParameters) args {
	return args{
		ctx: context.Background(),
		mg:  &v1alpha1.Record{Spec: v1alpha1.RecordSpec{ForProvider: testRecord}}}
}

var setupWant = func(resourceExists, resouceUpToDate bool) want {
	return want{
		o: managed.ExternalObservation{
			ResourceExists:    resourceExists,
			ResourceUpToDate:  resouceUpToDate,
			ConnectionDetails: managed.ConnectionDetails{},
		},
		err: nil}
}

var setupTestCase = func(name, reason string, resourceExists, resouceUpToDate bool, testRecord v1alpha1.RecordParameters) caseStructure {
	return caseStructure{name: {
		reason: reason,
		fields: fields{service: sjaclient.CreateSjaClient("token")},
		args:   setupArgs(testRecord),
		want:   setupWant(resourceExists, resouceUpToDate),
	},
	}
}

func runTestCases(t *testing.T, cases caseStructure) {
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{service: tc.fields.service}
			got, err := e.Observe(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}

// Observe() method test cases
func TestObserveNoChange(t *testing.T) {

	cases := setupTestCase("returns as object exists and upto date", "doesn't exist", true, true, v1alpha1.RecordParameters{
		Name:        "chai3",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	})

	runTestCases(t, cases)

}

func TestObserveDoesntExist(t *testing.T) {

	cases := setupTestCase("returns as object doesn't exist", "doesn't exist", false, false, v1alpha1.RecordParameters{
		Name:        "RECORD DOESN'T EXIST",
		Age:         11,
		Designation: "happiness CHANGED",
		Location:    "happiness",
		Todos:       []string{"gg"},
	})

	runTestCases(t, cases)
}

func TestObserveExistAndNotUpToDate(t *testing.T) {

	cases := setupTestCase("returns as object exist but not upto date", "not upto date", true, false, v1alpha1.RecordParameters{
		Name:        "chai3",
		Age:         11,
		Designation: "happiness CHANGED",
		Location:    "happiness",
		Todos:       []string{"gg"},
	})

	runTestCases(t, cases)
}
