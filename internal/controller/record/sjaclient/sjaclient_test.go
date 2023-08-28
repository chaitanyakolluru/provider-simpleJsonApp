package sjaclient

import (
	"context"
	"reflect"
	"testing"

	"git.heb.com/provider-simplejsonapp/apis/records/v1alpha1"
	"github.com/kelseyhightower/envconfig"
)

var env_vars struct {
	Token string
}

var testRecord v1alpha1.RecordParameters = v1alpha1.RecordParameters{
	Name:        "test-record",
	Age:         0,
	Designation: "desg",
	Location:    "loc",
	Todos:       []string{"todo1", "todo2"},
}

func processEnvVars() string {
	envconfig.Process("simplejsonapp", &env_vars)
	return env_vars.Token
}

func createRecord(sjaclient *SjaClient, t *testing.T) v1alpha1.RecordParameters {
	want, err := sjaclient.PostRecord(context.TODO(), testRecord)

	if err != nil {
		t.Errorf("error when creating test record: %s", err.Error())
	}

	return want
}

func deleteRecord(sjaclient *SjaClient, t *testing.T) v1alpha1.RecordParameters {
	want, err := sjaclient.DeleteRecord(context.TODO(), testRecord)

	if err != nil {
		t.Errorf("error when deleting test record: %s", err.Error())
	}

	return want
}

func TestGetRecords(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())

	want := createRecord(sjaclient, t)

	got, err := sjaclient.GetRecords(context.TODO())
	if err != nil {
		t.Errorf("GetRecords() failed with %s", err.Error())
	}

	for _, rp := range got {
		if rp.Name == want.Name {
			if !reflect.DeepEqual(rp, want) {
				t.Errorf("got: %v, want: %v", rp, want)
			}
		}
	}

	deleteRecord(sjaclient, t)
}

func TestGetRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())

	want := createRecord(sjaclient, t)
	got, err := sjaclient.GetRecord(context.TODO(), "test-record")

	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}

	deleteRecord(sjaclient, t)
}

func TestPostRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())
	got := createRecord(sjaclient, t)

	if !reflect.DeepEqual(got, testRecord) {
		t.Errorf("got: %v, want: %v", got, testRecord)
	}

	deleteRecord(sjaclient, t)
}

func TestPutRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())

	want := createRecord(sjaclient, t)

	request := want
	request.Designation = "designation changed"

	got, err := sjaclient.PutRecord(context.TODO(), request)
	if err != nil {
		t.Errorf("PutRecord() failed with %s", err.Error())
	}

	if !reflect.DeepEqual(got, request) {
		t.Errorf("got: %v, want: %v", got, request)
	}

	deleteRecord(sjaclient, t)
}

func TestDeleteRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())

	want := createRecord(sjaclient, t)

	got, err := sjaclient.DeleteRecord(context.TODO(), want)
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
