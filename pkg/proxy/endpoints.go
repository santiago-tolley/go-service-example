package proxy

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	MultiplyEndpoint endpoint.Endpoint
}

func (e Endpoints) Multiply(ctx context.Context, value, multiplier int) (int, error) {
	resp, err := e.MultiplyEndpoint(ctx, MultiplyRequest{Value: value, Multiplier: multiplier})
	if err != nil {
		return 0, err
	}
	response, ok := resp.(MultiplyResponse)
	if !ok {
		return 0, errors.New("Invalid response structure")
	}

	return response.Result, response.Err
}

func MakeEndpoints(p ProxyService) Endpoints {
	return Endpoints{
		MultiplyEndpoint: MakeMultiplyEndpoint(p),
	}
}

func MakeMultiplyEndpoint(p ProxyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(MultiplyRequest)
		if !ok {
			return nil, errors.New("Invalid request structure")
		}

		v, err := p.Multiply(ctx, req.Value, req.Multiplier)
		return MultiplyResponse{v, err}, nil
	}
}

type MultiplyRequest struct {
	Value      int `json:"value"`
	Multiplier int `json:"multiplier"`
}

type MultiplyResponse struct {
	Result int   `json:"result"`
	Err    error `json:"err"`
}
