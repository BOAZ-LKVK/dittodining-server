package repository

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/customerrors"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/domain/sample"
	"github.com/google/uuid"
)

type SampleRepository interface {
	FindOneSample(sampleID string) (*sample.Sample, error)
	CreateSample(sample *sample.Sample) (*sample.Sample, error)
	UpdateSample(sample *sample.Sample) (*sample.Sample, error)
	DeleteSample(sampleID string) error
	FindAllSamples() ([]*sample.Sample, error)
}

type sampleRepository struct {
	samples []*sample.Sample
}

func NewSampleRepository() SampleRepository {
	return &sampleRepository{
		samples: make([]*sample.Sample, 0),
	}
}

func (r *sampleRepository) FindOneSample(sampleID string) (*sample.Sample, error) {
	for _, s := range r.samples {
		if s.ID == sampleID {
			return s, nil
		}
	}

	return nil, customerrors.ErrorSampleNotFound
}

func (r *sampleRepository) CreateSample(sample *sample.Sample) (*sample.Sample, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	sample.ID = u.String()
	r.samples = append(r.samples, sample)
	return sample, nil
}

func (r *sampleRepository) UpdateSample(sample *sample.Sample) (*sample.Sample, error) {
	for i, s := range r.samples {
		if s.ID == sample.ID {
			r.samples[i] = sample
			return sample, nil
		}
	}

	return nil, customerrors.ErrorSampleNotFound
}

func (r *sampleRepository) DeleteSample(sampleID string) error {
	for i, s := range r.samples {
		if s.ID == sampleID {
			r.samples = append(r.samples[:i], r.samples[i+1:]...)
			return nil
		}
	}

	return customerrors.ErrorSampleNotFound
}

func (r *sampleRepository) FindAllSamples() ([]*sample.Sample, error) {
	return r.samples, nil
}
