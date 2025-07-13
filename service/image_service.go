package service

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

func MarkImagesAsUsed(db *mongo.Client, content string) error {
	re := regexp.MustCompile(regexp.QuoteMeta(config.BaseURL) + `/uploads/([a-f0-9\-]+\.[a-zA-Z]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return nil // Không có ảnh nào để xử lý
	}

	var filenames []string
	for _, match := range matches {
		if len(match) > 1 {
			filenames = append(filenames, match[1])
		}
	}

	if len(filenames) == 0 {
		return nil
	}

	collection := db.Database(config.DBName).Collection("UploadedImage")
	filter := bson.M{"filename": bson.M{"$in": filenames}}
	update := bson.M{"$set": bson.M{"status": models.ImageStatusUsed}}

	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}
