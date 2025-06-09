package account

import (
	"context"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *AccountServiceServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	update := bson.M{}

	if req.GetAccountName() != "" {
		update["account_name"] = req.GetAccountName()
	}
	if req.GetPassword() != "" {
		update["password"] = req.GetPassword()
	}

	if ui := req.GetUserInfo(); ui != nil {
		update["user_info"] = pbUserInfoToBSON(ui)
	}

	if len(req.GetRoleCodes()) > 0 {
		update["role_codes"] = req.GetRoleCodes()
	}

	if len(req.GetOfficeCodes()) > 0 {
		update["office_codes"] = req.GetOfficeCodes()
	}

	if len(req.GetAccessRights()) > 0 {
		update["access_rights"] = pbAccessRightsToBSON(req.GetAccessRights())
	}

	if len(update) == 0 {
		return &pb.UpdateAccountResponse{
			Success: false,
			Message: "No fields to update",
		}, nil
	}

	objID, err := primitive.ObjectIDFromHex(req.GetAccountId())
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return &pb.UpdateAccountResponse{
			Success: false,
			Message: "Account not found",
		}, nil
	}

	return &pb.UpdateAccountResponse{
		Success: true,
		Message: "Account updated successfully",
	}, nil
}
