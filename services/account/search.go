package account

import (
	"context"
	"log"
	"time"

	pb "legoas/legoas/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountServiceServer) SearchAccounts(ctx context.Context, req *pb.SearchAccountsRequest) (*pb.SearchAccountsResponse, error) {
	accountsColl := s.Mongo.Database("legoas").Collection("accounts")
	rolesColl := s.Mongo.Database("legoas").Collection("roles")
	officesColl := s.Mongo.Database("legoas").Collection("offices")
	menusColl := s.Mongo.Database("legoas").Collection("menus")

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

	cursor, err := accountsColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Database find error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []*pb.AccountData

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Decode error: %v", err)
			continue
		}

		account := &pb.AccountData{}

		if id, ok := doc["_id"].(primitive.ObjectID); ok {
			account.AccountId = id.Hex()
		}
		account.AccountName = getString(doc, "account_name")
		if createdAt, ok := doc["created_at"].(primitive.DateTime); ok {
			account.CreatedAt = createdAt.Time().Format(time.RFC3339)
		}

		if ui, ok := doc["user_info"].(bson.M); ok {
			account.UserInfo = &pb.UserInfo{
				Name:       getString(ui, "name"),
				Address:    getString(ui, "address"),
				PostalCode: getString(ui, "postal_code"),
				Province:   getString(ui, "province"),
				OfficeCode: getString(ui, "office_code"),
			}
		}

		account.Roles = fetchRoles(ctx, rolesColl, getStringArray(doc, "role_codes"))

		account.Offices = fetchOffices(ctx, officesColl, getStringArray(doc, "office_codes"))

		account.AccessRights = fetchAccessRights(ctx, menusColl, doc["access_rights"])

		results = append(results, account)
	}

	total, err := accountsColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Count error: %v", err)
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

func getStringArray(m bson.M, key string) []string {
	arr := []string{}
	if val, ok := m[key].(primitive.A); ok {
		for _, v := range val {
			if s, ok := v.(string); ok {
				arr = append(arr, s)
			}
		}
	}
	return arr
}

func fetchRoles(ctx context.Context, coll *mongo.Collection, codes []string) []*pb.Role {
	var roles []*pb.Role
	cursor, err := coll.Find(ctx, bson.M{"role_code": bson.M{"$in": codes}})
	if err != nil {
		return roles
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err == nil {
			roles = append(roles, &pb.Role{
				RoleCode: getString(doc, "role_code"),
				RoleName: getString(doc, "role_name"),
			})
		}
	}
	return roles
}

func fetchOffices(ctx context.Context, coll *mongo.Collection, codes []string) []*pb.Office {
	var offices []*pb.Office
	cursor, err := coll.Find(ctx, bson.M{"office_code": bson.M{"$in": codes}})
	if err != nil {
		return offices
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err == nil {
			offices = append(offices, &pb.Office{
				OfficeCode: getString(doc, "office_code"),
				OfficeName: getString(doc, "office_name"),
			})
		}
	}
	return offices
}

func fetchAccessRights(ctx context.Context, menusColl *mongo.Collection, value interface{}) []*pb.AccessRight {
	var accessRights []*pb.AccessRight
	if ars, ok := value.(primitive.A); ok {
		for _, a := range ars {
			if ar, ok := a.(bson.M); ok {
				menuCode := getString(ar, "menu_code")
				menuName := ""

				var menuDoc bson.M
				err := menusColl.FindOne(ctx, bson.M{"menu_code": menuCode}).Decode(&menuDoc)
				if err == nil {
					menuName = getString(menuDoc, "menu_name")
				}

				accessRights = append(accessRights, &pb.AccessRight{
					MenuCode:  menuCode,
					MenuName:  menuName,
					CanCreate: getBool(ar, "can_create"),
					CanRead:   getBool(ar, "can_read"),
					CanUpdate: getBool(ar, "can_update"),
					CanDelete: getBool(ar, "can_delete"),
				})
			}
		}
	}
	return accessRights
}
