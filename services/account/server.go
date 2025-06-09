package account

import (
	"context"
	"log"
	"time"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	accessRights := []bson.M{}
	for _, ar := range req.GetAccessRights() {
		accessRights = append(accessRights, bson.M{
			"menu_code":  ar.GetMenuCode(),
			"can_create": ar.GetCanCreate(),
			"can_read":   ar.GetCanRead(),
			"can_update": ar.GetCanUpdate(),
			"can_delete": ar.GetCanDelete(),
		})
	}

	accountDoc := bson.M{
		"account_name": req.GetAccountName(),
		"password":     req.GetPassword(),
		"user_info": bson.M{
			"name":        req.GetUserInfo().GetName(),
			"address":     req.GetUserInfo().GetAddress(),
			"postal_code": req.GetUserInfo().GetPostalCode(),
			"province":    req.GetUserInfo().GetProvince(),
			"office_code": req.GetUserInfo().GetOfficeCode(),
		},
		"role_codes":    req.GetRoleCodes(),
		"office_codes":  req.GetOfficeCodes(),
		"access_rights": accessRights,
		"created_at":    time.Now(),
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

func (s *AccountServiceServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	update := bson.M{}

	if req.GetAccountName() != "" {
		update["account_name"] = req.GetAccountName()
	}
	if req.GetPassword() != "" {
		update["password"] = req.GetPassword()
	}

	// Update user_info if provided
	if ui := req.GetUserInfo(); ui != nil {
		update["user_info"] = bson.M{
			"name":        ui.GetName(),
			"address":     ui.GetAddress(),
			"postal_code": ui.GetPostalCode(),
			"province":    ui.GetProvince(),
			"office_code": ui.GetOfficeCode(),
		}
	}

	// Update role_codes if any
	if len(req.GetRoleCodes()) > 0 {
		update["role_codes"] = req.GetRoleCodes()
	}

	// Update office_codes if any
	if len(req.GetOfficeCodes()) > 0 {
		update["office_codes"] = req.GetOfficeCodes()
	}

	// Update access_rights if any
	if len(req.GetAccessRights()) > 0 {
		accessRights := []bson.M{}
		for _, ar := range req.GetAccessRights() {
			accessRights = append(accessRights, bson.M{
				"menu_code":  ar.GetMenuCode(),
				"can_create": ar.GetCanCreate(),
				"can_read":   ar.GetCanRead(),
				"can_update": ar.GetCanUpdate(),
				"can_delete": ar.GetCanDelete(),
			})
		}
		update["access_rights"] = accessRights
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
