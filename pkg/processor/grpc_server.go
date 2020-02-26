package processor

import (
	"context"
	"errors"
	"go-service-example/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	calculate grpctransport.Handler
}

func NewGRPCServer(endpoints Endpoints) pb.ProcessorServer {
	return &grpcServer{
		calculate: grpctransport.NewServer(
			endpoints.CalculateEndpoint,
			decodeGRPCCalculateRequest,
			encodeGRPCCalculateResponse,
		),
	}
}

func (s *grpcServer) Calculate(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	_, r, err := s.calculate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, ok := r.(*pb.ProcessResponse)
	if !ok {
		return nil, errors.New("Invalid response structure")
	}

	return resp, nil
}

func decodeGRPCCalculateRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.ProcessRequest)
	if !ok {
		return nil, errors.New("Invalid request structure")
	}

	return CalculateRequest{
		Value:      int(req.Value),
		Multiplier: int(req.Multiplier),
	}, nil
}

func encodeGRPCCalculateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(CalculateResponse)
	if !ok {
		return nil, errors.New("Invalid response structure")
	}

	return &pb.ProcessResponse{
		Result: int64(resp.Result),
		Err:    err2str(resp.Err),
	}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
