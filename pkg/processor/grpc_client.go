package processor

import (
	"context"
	"errors"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"go-service-example/pb"
)

func NewGRPCClient(conn *grpc.ClientConn) ProcessorService {
	calculateEndpoint := grpctransport.NewClient(
		conn,
		"pb.Processor",
		"Calculate",
		encodeGRPCProcessRequest,
		decodeGRPCProcessResponse,
		pb.ProcessResponse{},
	).Endpoint()

	return Endpoints{
		CalculateEndpoint: calculateEndpoint,
	}
}

func encodeGRPCProcessRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(CalculateRequest)
	if !ok {
		return nil, errors.New("Invalid request structure")
	}

	return &pb.ProcessRequest{
		Value:      int64(req.Value),
		Multiplier: int64(req.Multiplier),
	}, nil
}

func decodeGRPCProcessResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.ProcessResponse)
	if !ok {
		return nil, errors.New("Invalid response structure")
	}

	return CalculateResponse{
		Result: int(reply.Result),
		Err:    str2err(reply.Err),
	}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}
