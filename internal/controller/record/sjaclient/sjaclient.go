package sjaclient

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/carlmjohnson/requests"
	v1alpha1 "github.com/chaitanyakolluru/provider-simplejsonapp/apis/records/v1alpha1"
)

const SIMPLE_JSON_APP_BASEURL = "http://localhost:8080/json"
const (
	errGetError = "/records error"
)

type SjaClient struct {
	ctx context.Context
}

func CreatSjaClient(ctx context.Context) *SjaClient {
	return &SjaClient{ctx: ctx}
}

func (s *SjaClient) GetRecords() ([]v1alpha1.RecordParameters, error) {
	var response []v1alpha1.RecordParameters

	err := requests.URL(fmt.Sprintf("%s/records", SIMPLE_JSON_APP_BASEURL)).ToJSON(&response).Fetch(s.ctx)
	if err != nil {
		return []v1alpha1.RecordParameters{}, errors.New(errGetError)
	}

	return response, nil
}

func (s *SjaClient) GetRecord(name string) (v1alpha1.RecordParameters, error) {
	var response v1alpha1.RecordParameters

	err := requests.URL(fmt.Sprintf("%s/records/%s", SIMPLE_JSON_APP_BASEURL, name)).ToJSON(&response).Fetch(s.ctx)
	if err != nil {
		return v1alpha1.RecordParameters{}, errors.New(errGetError)
	}

	return response, nil
}
