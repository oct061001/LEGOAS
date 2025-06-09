package account

import (
	"context"
	"time"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountServiceServer) GetAccountById(ctx context.Context, req *pb.GetAccountByIdRequest) (*pb.GetAccountByIdResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	objID, err := primitive.ObjectIDFromHex(req.GetAccountId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid account ID: %v", err)
	}

	var accountDoc AccountDoc

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&accountDoc)
	if err == mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get account: %v", err)
	}

	return &pb.GetAccountByIdResponse{
		AccountId:    accountDoc.ID.Hex(),
		AccountName:  accountDoc.AccountName,
		UserInfo:     bsonUserInfoToPB(accountDoc.UserInfo),
		RoleCodes:    accountDoc.RoleCodes,
		OfficeCodes:  accountDoc.OfficeCodes,
		AccessRights: bsonAccessRightsToPB(accountDoc.AccessRights),
		CreatedAt:    accountDoc.CreatedAt.Format(time.RFC3339),
	}, nil
}
