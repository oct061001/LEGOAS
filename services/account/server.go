package account

import (
	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/mongo"
)

type AccountServiceServer struct {
	pb.UnimplementedAccountServiceServer
	Mongo *mongo.Client
}

func NewAccountServiceServer(mongoClient *mongo.Client) *AccountServiceServer {
	return &AccountServiceServer{
		Mongo: mongoClient,
	}
}
