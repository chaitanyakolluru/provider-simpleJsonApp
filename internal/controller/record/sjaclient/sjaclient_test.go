package sjaclient

import (
	"context"
	"fmt"
	"testing"
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
