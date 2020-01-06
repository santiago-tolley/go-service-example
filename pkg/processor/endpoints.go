package processor

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/log"
)

type Endpoints struct {
	CalculateEndpoint endpoint.Endpoint
}

func (e Endpoints) Calculate(ctx context.Context, value int, multiplier int) (int, error) {
	resp, err := e.CalculateEndpoint(ctx, CalculateRequest{Value: value, Multiplier: multiplier})
	if err != nil {
		return 0, err
	}
	response := resp.(CalculateResponse)

	return response.Result, response.Err
}

func MakeEndpoints(p ProcessorService) Endpoints {
	return Endpoints{
		CalculateEndpoint: MakeCalculateEndpoint(p),
	}
}

func MakeCalculateEndpoint(p ProcessorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CalculateRequest)
		v, err := p.Calculate(ctx, req.Value, req.Multiplier)
		if err != nil {
			return CalculateResponse{0, err}, nil
		}
		return CalculateResponse{v, nil}, nil
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
