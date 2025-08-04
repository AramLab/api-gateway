package clients

import (
	"context"
	"github.com/AramLab/api-gateway/internal/handlers"
	pb "github.com/AramLab/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

func NewAuthClient(addr string) (handlers.AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &authClient{conn: conn, client: pb.NewAuthServiceClient(conn)}, nil
}

func (a *authClient) Close() error {
	return a.conn.Close()
}

func (a *authClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return a.client.Login(ctx, &pb.LoginRequest{Username: req.Username, Password: req.Password})
}

func (a *authClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return a.client.Register(ctx, &pb.RegisterRequest{Username: req.Username, Password: req.Password, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName})
}
