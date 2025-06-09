package role

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Role struct {
	RoleCode string `bson:"role_code"`
	RoleName string `bson:"role_name"`
}

type RoleRepository struct {
	collection *mongo.Collection
}

func NewRoleRepository(db *mongo.Database) *RoleRepository {
	return &RoleRepository{
		collection: db.Collection("roles"),
	}
}

func (r *RoleRepository) SearchRoles(ctx context.Context, query string) ([]Role, error) {
	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"role_code": bson.M{"$regex": query, "$options": "i"}},
				{"role_name": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}
