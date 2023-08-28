package sjaclient

import (
	"context"
	"fmt"
	"testing"

	"git.heb.com/provider-simplejsonapp/apis/records/v1alpha1"
	"github.com/kelseyhightower/envconfig"
)

var env_vars struct {
	Token string
}

func processEnvVars() string {
	envconfig.Process("simplejsonapp", &env_vars)
	return env_vars.Token
}

func TestGetRecords(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())
	got, err := sjaclient.GetRecords(context.TODO())
	if err != nil {
		t.Errorf("GetRecords() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestGetRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())
	got, err := sjaclient.GetRecord(context.TODO(), "chai")
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestPostRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())
	request := v1alpha1.RecordParameters{
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.PostRecord(context.TODO(), request)
	if err != nil {
		t.Errorf("PostRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestPutRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())
	request := v1alpha1.RecordParameters{
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.PutRecord(context.TODO(), request)
	if err != nil {
		t.Errorf("PutRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestDeleteRecord(t *testing.T) {
	sjaclient := CreateSjaClient(processEnvVars())
	request := v1alpha1.RecordParameters{
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.DeleteRecord(context.TODO(), request)
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}
