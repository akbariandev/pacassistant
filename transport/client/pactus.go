package client

import (
	"fmt"

	"github.com/akbariandev/pacassistant/config"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pactusPB "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pactus struct {
	Blockchain  pactusPB.BlockchainClient
	Transaction pactusPB.TransactionClient
	Network     pactusPB.NetworkClient
}

func NewPactusClient(cfg *config.GrpcClient) (*Pactus, error) {
	address := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	dialOpts := make([]grpc.DialOption, 0)
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()))

	conn, err := grpc.Dial(address, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &Pactus{
		pactusPB.NewBlockchainClient(conn),
		pactusPB.NewTransactionClient(conn),
		pactusPB.NewNetworkClient(conn),
	}, nil
}
