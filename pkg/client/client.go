package client

import (
	"fmt"
	"time"

	"github.com/vardius/blockchain/pkg/blockchain"
	"github.com/vardius/blockchain/pkg/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// BlockchainClient interface
type BlockchainClient interface {
	AddBlock(string) error
	GetBlockchain() (blockchain.Chain, error)
}

type blockchainClient struct {
	host string
	port int
}

// New creates new blockchain client
func New(host string, port int) BlockchainClient {
	return &blockchainClient{host, port}
}

func (c *blockchainClient) AddBlock(data string) error {
	conn, err := c.getConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := proto.NewBlockchainClient(conn)

	_, err = client.AddBlock(context.Background(), &proto.AddBlockRequest{
		Data: data,
	})

	return err
}

func (c *blockchainClient) GetBlockchain() (blockchain.Chain, error) {
	var bchain blockchain.Chain

	conn, err := c.getConnection()
	if err != nil {
		return bchain, err
	}
	defer conn.Close()

	client := proto.NewBlockchainClient(conn)

	chain, err := client.GetBlockchain(context.Background(), &proto.GetBlockchainRequest{})
	if err != nil {
		return bchain, err
	}

	for _, b := range chain.Blocks {
		bchain.Append(&blockchain.Block{
			Index:    b.Index,
			Time:     time.Unix(b.Timestamp, 0),
			Data:     b.Data,
			PrevHash: b.Hash,
		})
	}

	return bchain, nil
}

func (c *blockchainClient) getConnection() (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("%s:%d", c.host, c.port), grpc.WithInsecure())
}
