package proxy

import (
	"context"
	"go-service-example/pkg/processor"
)

type ProxyService interface {
	Multiply(context.Context, int, int) (int, error)
}

func NewProxyServer(grpcClient processor.ProcessorService) ProxyService {
	return proxyService{client: grpcClient}
}

type proxyService struct {
	client processor.ProcessorService
}

func (p proxyService) Multiply(ctx context.Context, value, multiplier int) (int, error) {
	result, err := p.client.Calculate(ctx, value, multiplier)
	return result, err
}
