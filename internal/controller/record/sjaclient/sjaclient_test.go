package sjaclient

import (
	"context"
	"fmt"
	"testing"

	"github.com/chaitanyakolluru/provider-simplejsonapp/apis/records/v1alpha1"
)

func TestGetRecords(t *testing.T) {
	sjaclient := CreateSjaClient()
	got, err := sjaclient.GetRecords(context.Background())
	if err != nil {
		t.Errorf("GetRecords() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestGetRecord(t *testing.T) {
	sjaclient := CreateSjaClient()
	got, err := sjaclient.GetRecord(context.Background(), "chai")
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestPostRecord(t *testing.T) {
	sjaclient := CreateSjaClient()
	request := v1alpha1.RecordParameters{
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.PostRecord(context.Background(), request)
	if err != nil {
		t.Errorf("PostRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestPutRecord(t *testing.T) {
	sjaclient := CreateSjaClient()
	request := v1alpha1.RecordParameters{
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.PutRecord(context.Background(), request)
	if err != nil {
		t.Errorf("PutRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestDeleteRecord(t *testing.T) {
	sjaclient := CreateSjaClient()
	request := v1alpha1.RecordParameters{
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.DeleteRecord(context.Background(), request)
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}
