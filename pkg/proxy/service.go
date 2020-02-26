package proxy

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
)

type ProcessorService interface {
	Calculate(context.Context, int, int) (int, error)
}

func NewProxyServer(grpcClient ProcessorService) ProxyService {
	return ProxyService{client: grpcClient}
}

type ProxyService struct {
	client ProcessorService
}

func (p ProxyService) Multiply(ctx context.Context, value, multiplier int) (int, error) {
	result, err := p.client.Calculate(ctx, value, multiplier)
	if err != nil {
		logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger.Log("transport", "failed to connect to processor service")
		return result, ConnectionError{"failed to connect to processor service"}
	}
	return result, nil
}

type ConnectionError struct {
	Message string `json:"message"`
}

func (c ConnectionError) Error() string {
	return c.Message
}
