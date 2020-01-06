package processor

import (
	"context"
)

type ProcessorService interface {
	Calculate(context.Context, int, int) (int, error)
}

func NewProcessorServer() ProcessorService {
	return processorService{}
}

type processorService struct{}

func (p processorService) Calculate(_ context.Context, value int, multiplier int) (int, error) {
	return value * multiplier, nil
}
