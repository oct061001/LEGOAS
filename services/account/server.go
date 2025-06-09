package account

import (
	"context"
	"log"
	"time"

	pb "legoas/proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountServiceServer struct {
	pb.UnimplementedAccountServiceServer
	Mongo *mongo.Client
}

type Account struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AccountName string             `bson:"account_name"`
	Password    string             `bson:"password"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func NewAccountServiceServer(mongoClient *mongo.Client) *AccountServiceServer {
	return &AccountServiceServer{
		Mongo: mongoClient,
	}
}

func (s *AccountServiceServer) RegisterAccount(ctx context.Context, req *pb.RegisterAccountRequest) (*pb.RegisterAccountResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	account := Account{
		AccountName: req.AccountName,
		Password:    req.Password,
		CreatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, account)
	if err != nil {
		log.Println("Failed to insert account:", err)
		return nil, err
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return &pb.RegisterAccountResponse{
		AccountId: oid.Hex(),
	}, nil
}
