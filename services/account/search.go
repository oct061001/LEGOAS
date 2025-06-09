package account

import (
	"context"
	"log"
	"time"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountServiceServer) SearchAccounts(ctx context.Context, req *pb.SearchAccountsRequest) (*pb.SearchAccountsResponse, error) {
	collection := s.Mongo.Database("legoas").Collection("accounts")

	// Build filter
	filter := bson.M{}
	if req.GetQuery() != "" {
		filter["$or"] = []bson.M{
			{"account_name": bson.M{"$regex": req.GetQuery(), "$options": "i"}},
			{"user_info.name": bson.M{"$regex": req.GetQuery(), "$options": "i"}},
		}
	}
	if req.GetRoleCode() != "" {
		filter["role_codes"] = req.GetRoleCode()
	}
	if req.GetOfficeCode() != "" {
		filter["office_codes"] = req.GetOfficeCode()
	}

	// Pagination
	page := req.GetPage()
	if page < 1 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize < 1 {
		pageSize = 10
	}
	skip := int64((page - 1) * pageSize)

	opts := options.Find().SetSkip(skip).SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Database find error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []*pb.AccountData
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Failed to decode document: %v", err)
			continue
		}

		account := &pb.AccountData{}

		if id, ok := doc["_id"].(primitive.ObjectID); ok {
			account.AccountId = id.Hex()
		}
		if name, ok := doc["account_name"].(string); ok {
			account.AccountName = name
		}
		if createdAt, ok := doc["created_at"].(primitive.DateTime); ok {
			account.CreatedAt = createdAt.Time().Format(time.RFC3339)
		}

		// Parse user_info
		if ui, ok := doc["user_info"].(bson.M); ok {
			account.UserInfo = &pb.UserInfo{
				Name:       getString(ui, "name"),
				Address:    getString(ui, "address"),
				PostalCode: getString(ui, "postal_code"),
				Province:   getString(ui, "province"),
				OfficeCode: getString(ui, "office_code"),
			}
		}

		// Parse role codes
		if rc, ok := doc["role_codes"].(primitive.A); ok {
			for _, r := range rc {
				if role, ok := r.(string); ok {
					account.RoleCodes = append(account.RoleCodes, role)
				}
			}
		}

		// Parse office codes
		if oc, ok := doc["office_codes"].(primitive.A); ok {
			for _, o := range oc {
				if office, ok := o.(string); ok {
					account.OfficeCodes = append(account.OfficeCodes, office)
				}
			}
		}

		// Parse access rights
		if ars, ok := doc["access_rights"].(primitive.A); ok {
			for _, a := range ars {
				if ar, ok := a.(bson.M); ok {
					account.AccessRights = append(account.AccessRights, &pb.AccessRight{
						MenuCode:  getString(ar, "menu_code"),
						CanCreate: getBool(ar, "can_create"),
						CanRead:   getBool(ar, "can_read"),
						CanUpdate: getBool(ar, "can_update"),
						CanDelete: getBool(ar, "can_delete"),
					})
				}
			}
		}

		results = append(results, account)
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count documents: %v", err)
	}

	return &pb.SearchAccountsResponse{
		Accounts:   results,
		TotalCount: int32(total),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func getString(m bson.M, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func getBool(m bson.M, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
