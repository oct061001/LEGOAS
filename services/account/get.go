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
		Roles:        bsonRolesToPB(accountDoc.RoleCodes),
		Offices:      bsonOfficesToPB(accountDoc.OfficeCodes),
		AccessRights: bsonAccessRightsToPB(accountDoc.AccessRights),
		CreatedAt:    accountDoc.CreatedAt.Format(time.RFC3339),
	}, nil
}

func bsonRolesToPB(codes []string) []*pb.Role {
	var roles []*pb.Role
	for _, code := range codes {
		roles = append(roles, &pb.Role{
			RoleCode: code,
			// RoleName can be filled from DB or omitted if not available
		})
	}
	return roles
}

func bsonOfficesToPB(codes []string) []*pb.Office {
	var offices []*pb.Office
	for _, code := range codes {
		offices = append(offices, &pb.Office{
			OfficeCode: code,
			// OfficeName can be filled from DB or omitted if not available
		})
	}
	return offices
}
