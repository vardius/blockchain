package server

import (
	"fmt"
	"log"
	"net"

	"github.com/vardius/blockchain/pkg/blockchain"
	"github.com/vardius/blockchain/pkg/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	blockchain blockchain.Chain
}

// Run runs gRPC server
func Run(host string, port int) error {
	srv := grpc.NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	proto.RegisterBlockchainServer(srv, &server{
		blockchain: blockchain.NewBlockchain(),
	})

	log.Println(fmt.Sprintf("Server listening on %s:%d", host, port))

	return srv.Serve(lis)
}

// AddBlock adds new block to blockchain
func (s *server) AddBlock(ctx context.Context, in *proto.AddBlockRequest) (*proto.AddBlockResponse, error) {
	parent := s.blockchain[len(s.blockchain)-1]
	newBlock, err := blockchain.NewBlock(parent, in.Data)

	if err != nil {
		return new(proto.AddBlockResponse), err
	}

	err = s.blockchain.Append(newBlock)

	return new(proto.AddBlockResponse), err
}

// GetBlockchain returns blockchain
func (s *server) GetBlockchain(ctx context.Context, in *proto.GetBlockchainRequest) (*proto.GetBlockchainResponse, error) {
	resp := new(proto.GetBlockchainResponse)

	for _, b := range s.blockchain {
		resp.Blocks = append(resp.Blocks, &proto.Block{
			Index:     b.Index,
			Timestamp: b.Time.Unix(),
			Hash:      b.Hash,
			PrevHash:  b.PrevHash,
			Data:      b.Data,
		})
	}

	return resp, nil
}
