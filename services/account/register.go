package account

import (
	"context"
	"log"
	"time"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *AccountServiceServer) RegisterAccount(ctx context.Context, req *pb.RegisterAccountRequest) (*pb.RegisterAccountResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	accountDoc := AccountDoc{
		AccountName:  req.GetAccountName(),
		Password:     req.GetPassword(),
		UserInfo:     pbUserInfoToBSON(req.GetUserInfo()),
		RoleCodes:    req.GetRoleCodes(),
		OfficeCodes:  req.GetOfficeCodes(),
		AccessRights: pbAccessRightsToBSON(req.GetAccessRights()),
		CreatedAt:    time.Now(),
	}

	result, err := collection.InsertOne(ctx, accountDoc)
	if err != nil {
		log.Println("Failed to insert account:", err)
		return nil, err
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return &pb.RegisterAccountResponse{
		AccountId: oid.Hex(),
	}, nil
}
