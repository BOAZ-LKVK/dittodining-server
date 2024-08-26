package sample

import "github.com/BOAZ-LKVK/LKVK-server/pkg/domain/sample"

type CreateSampleResponse struct {
	Sample *sample.Sample `json:"sample"`
}

type ListSamplesResponse struct {
	Samples []*sample.Sample `json:"samples"`
}

type GetSampleResponse struct {
	Sample *sample.Sample `json:"sample"`
}

type UpdateSampleResponse struct {
	Sample *sample.Sample `json:"sample"`
}

type DeleteSampleResponse struct{}
