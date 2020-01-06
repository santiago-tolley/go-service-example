package proxy

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/log"
)

type Endpoints struct {
	MultiplyEndpoint endpoint.Endpoint
}

func (e Endpoints) Multiply(ctx context.Context, value int, multiplier int) (int, error) {
	resp, err := e.MultiplyEndpoint(ctx, MultiplyRequest{Value: value, Multiplier: multiplier})
	if err != nil {
		return 0, err
	}
	response := resp.(MultiplyResponse)

	return response.Result, response.Err
}

func MakeEndpoints(p ProxyService) Endpoints {
	return Endpoints{
		MultiplyEndpoint: MakeMultiplyEndpoint(p),
	}
}

func MakeMultiplyEndpoint(p ProxyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MultiplyRequest)
		v, err := p.Multiply(ctx, req.Value, req.Multiplier)
		if err != nil {
			return MultiplyResponse{0, err}, nil
		}
		return MultiplyResponse{v, nil}, nil
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
