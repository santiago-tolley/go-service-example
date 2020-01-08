package processor

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CalculateEndpoint endpoint.Endpoint
}

func (e Endpoints) Calculate(ctx context.Context, value, multiplier int) (int, error) {
	resp, err := e.CalculateEndpoint(ctx, CalculateRequest{Value: value, Multiplier: multiplier})
	if err != nil {
		return 0, err
	}
	response, ok := resp.(CalculateResponse)
	if !ok {
		return 0, errors.New("Invalid response structure")
	}

	return response.Result, response.Err
}

func MakeEndpoints(p ProcessorService) Endpoints {
	return Endpoints{
		CalculateEndpoint: MakeCalculateEndpoint(p),
	}
}

func MakeCalculateEndpoint(p ProcessorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CalculateRequest)
		if !ok {
			return nil, errors.New("Invalid request structure")
		}

		v, err := p.Calculate(ctx, req.Value, req.Multiplier)
		return CalculateResponse{v, err}, nil
	}
}

type CalculateRequest struct {
	Value      int `json:"value"`
	Multiplier int `json:"multiplier"`
}

type CalculateResponse struct {
	Result int   `json:"result"`
	Err    error `json:"err"`
}
