package repository

import (
	"context"
	"fmt"
	"time"

	"url-shortener/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"url-shortener/internal/model"
)

type Repository interface {
	Save(shortURL model.ShortURL) (model.ShortURL, error)
	Get(shortCode string) (string, error)
}

type mongoRepository struct {
	collection *mongo.Collection
	redisDB 	db.RedisDB
}

func NewMongoRepository(client *mongo.Client,redisDB db.RedisDB) Repository {
	collection := client.Database("urlshortener").Collection("urls")
	return &mongoRepository{
		collection: collection,
		redisDB:    redisDB,
	}
}

func (r *mongoRepository) Save(shortURL model.ShortURL) (model.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inserted, err := r.collection.InsertOne(ctx, shortURL)
	if err != nil {
		return shortURL, err
	}
	shortURL.ID = inserted.InsertedID.(string)
	r.redisDB.Save(shortURL.ShortCode, shortURL.OriginalURL)
	return shortURL, nil
}

func (r *mongoRepository) Get(shortCode string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	originalURL, err := r.redisDB.Get(shortCode)
	if err == nil {
		return *originalURL, nil
	}
	var shortURL model.ShortURL
	filter := bson.M{"short_code": shortCode} 

	err = r.collection.FindOne(ctx, filter).Decode(&shortURL)
	if err != nil {
		return "", err
	}
	err = r.redisDB.Save(shortURL.ShortCode, shortURL.OriginalURL)
	if err != nil {
		fmt.Println(err)
	}
	return shortURL.OriginalURL, nil
}
