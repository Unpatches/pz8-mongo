package notes

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("note not found")

type Repo struct {
	col *mongo.Collection
}

type Stats struct {
	Count         int64   `bson:"count"         json:"count"`
	AvgContentLen float64 `bson:"avgContentLen" json:"avgContentLen"`
}

func NewRepo(db *mongo.Database) (*Repo, error) {
	col := db.Collection("notes")
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	_, err = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "content", Value: "text"},
		},
		Options: options.Index().SetDefaultLanguage("russian"),
	})
	if err != nil {
		return nil, err
	}

	_, err = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil {
		return nil, err
	}

	return &Repo{col: col}, nil
}

func (r *Repo) Create(ctx context.Context, title, content string) (Note, error) {
	now := time.Now()
	n := Note{
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	res, err := r.col.InsertOne(ctx, n)
	if err != nil {
		return Note{}, err
	}
	n.ID = res.InsertedID.(primitive.ObjectID)
	return n, nil
}

func (r *Repo) ByID(ctx context.Context, idHex string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}
	var n Note
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return n, nil
}

func (r *Repo) List(ctx context.Context, q string, limit, skip int64) ([]Note, error) {
	filter := bson.M{}
	opts := options.Find().SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	if skip < 0 {
		skip = 0
	}
	opts.SetSkip(skip)

	if q != "" {
		filter["$text"] = bson.M{"$search": q}
		opts.SetProjection(bson.D{
			{Key: "score", Value: bson.M{"$meta": "textScore"}},
		})
		opts.SetSort(bson.D{
			{Key: "score", Value: bson.M{"$meta": "textScore"}},
		})
	}

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}

	return out, cur.Err()
}

func (r *Repo) ListAfter(ctx context.Context, q, afterID string, limit int64) ([]Note, error) {
	filter := bson.M{}

	if q != "" {
		filter["$text"] = bson.M{"$search": q}
	}

	if afterID != "" {
		if oid, err := primitive.ObjectIDFromHex(afterID); err == nil {
			filter["_id"] = bson.M{"$lt": oid}
		}
	}

	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, cur.Err()
}

func (r *Repo) Stats(ctx context.Context) (Stats, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$project", Value: bson.M{
			"contentLen": bson.M{"$strLenCP": "$content"},
		}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"count":         bson.M{"$sum": 1},
			"avgContentLen": bson.M{"$avg": "$contentLen"},
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return Stats{}, err
	}
	defer cur.Close(ctx)

	if cur.Next(ctx) {
		var s Stats
		if err := cur.Decode(&s); err != nil {
			return Stats{}, err
		}
		return s, nil
	}
	if err := cur.Err(); err != nil {
		return Stats{}, err
	}

	return Stats{Count: 0, AvgContentLen: 0}, nil
}

func (r *Repo) Update(ctx context.Context, idHex string, title, content *string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}

	set := bson.M{"updatedAt": time.Now()}
	if title != nil {
		set["title"] = *title
	}
	if content != nil {
		set["content"] = *content
	}

	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Note
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, after).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return updated, nil
}

func (r *Repo) Delete(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return ErrNotFound
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
