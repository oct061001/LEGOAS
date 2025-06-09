package account

import (
	"context"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountServiceServer) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	objectID, err := primitive.ObjectIDFromHex(req.GetAccountId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account_id format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete account: %v", err)
	}

	return &pb.DeleteAccountResponse{
		Success: res.DeletedCount > 0,
	}, nil
}
