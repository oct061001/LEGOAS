package office

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Office struct {
	OfficeCode string `bson:"office_code"`
	OfficeName string `bson:"office_name"`
}

type OfficeRepository struct {
	collection *mongo.Collection
}

func NewOfficeRepository(db *mongo.Database) *OfficeRepository {
	return &OfficeRepository{
		collection: db.Collection("offices"),
	}
}

func (r *OfficeRepository) SearchOffices(ctx context.Context, query string) ([]Office, error) {
	filter := bson.M{}
	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"office_code": bson.M{"$regex": query, "$options": "i"}},
				{"office_name": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var offices []Office
	if err := cursor.All(ctx, &offices); err != nil {
		return nil, err
	}
	return offices, nil
}
