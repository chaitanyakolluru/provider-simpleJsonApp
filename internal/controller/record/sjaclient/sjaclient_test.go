package sjaclient

import (
	"context"
	"fmt"
	"testing"

	"github.com/chaitanyakolluru/provider-simplejsonapp/apis/records/v1alpha1"
)

func TestGetRecords(t *testing.T) {
	sjaclient := CreatSjaClient(context.Background())
	got, err := sjaclient.GetRecords()
	if err != nil {
		t.Errorf("GetRecords() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestGetRecord(t *testing.T) {
	sjaclient := CreatSjaClient(context.Background())
	got, err := sjaclient.GetRecord("chai")
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestPostRecord(t *testing.T) {
	sjaclient := CreatSjaClient(context.Background())
	request := v1alpha1.RecordParameters{
		Id:          2,
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.PostRecord(request)
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestPutRecord(t *testing.T) {
	sjaclient := CreatSjaClient(context.Background())
	request := v1alpha1.RecordParameters{
		Id:          2,
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.PutRecord(request)
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}

func TestDeleteRecord(t *testing.T) {
	sjaclient := CreatSjaClient(context.Background())
	request := v1alpha1.RecordParameters{
		Id:          2,
		Name:        "chai2",
		Age:         11,
		Designation: "happiness",
		Location:    "happiness",
		Todos:       []string{"gg"},
	}
	got, err := sjaclient.DeleteRecord(request)
	if err != nil {
		t.Errorf("GetRecord() failed with %s", err.Error())
	}
	fmt.Println(got)
}
