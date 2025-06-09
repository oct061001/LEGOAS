package menu

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Menu struct {
	MenuCode string `bson:"menu_code"`
	MenuName string `bson:"menu_name"`
}

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository(db *mongo.Database) *MenuRepository {
	return &MenuRepository{
		collection: db.Collection("menus"),
	}
}

func (r *MenuRepository) SearchMenus(ctx context.Context, query string) ([]Menu, error) {
	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"menu_code": bson.M{"$regex": query, "$options": "i"}},
				{"menu_name": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []Menu
	if err := cursor.All(ctx, &menus); err != nil {
		return nil, err
	}
	return menus, nil
}
