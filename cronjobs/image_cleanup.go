package cronjobs

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const GracePeriod = -24 * time.Hour

func RunImageCleanupJob(dbClient *mongo.Client) {
	log.Println("--- [CRON] Starting image cleanup job ---")

	imageCollection := dbClient.Database(config.DBName).Collection("UploadedImage")

	cutoffTime := time.Now().Add(GracePeriod)
	log.Printf("[CRON] Cleanup threshold: Deleting 'pending' images created before %s\n", cutoffTime.Format(time.RFC3339))

	filter := bson.M{
		"status": models.ImageStatusPending,
		"createdAt": bson.M{
			"$lt": cutoffTime,
		},
	}

	cursor, err := imageCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("[CRON-ERROR] Failed to find images for cleanup: %v\n", err)
		return
	}
	defer cursor.Close(context.TODO())

	var imagesToDelete []models.UploadedImage
	if err = cursor.All(context.TODO(), &imagesToDelete); err != nil {
		log.Printf("[CRON-ERROR] Failed to decode images: %v\n", err)
		return
	}

	if len(imagesToDelete) == 0 {
		log.Println("[CRON] No old pending images found to delete. Job finished.")
		return
	}

	log.Printf("[CRON] Found %d images to delete.\n", len(imagesToDelete))
	deletedCount := 0
	errorCount := 0

	for _, image := range imagesToDelete {
		if err := os.Remove(image.Path); err != nil {
			if os.IsNotExist(err) {
				log.Printf("[CRON-WARN] File not found, will still remove DB record: %s\n", image.Path)
			} else {
				log.Printf("[CRON-ERROR] Failed deleting file %s: %v\n", image.Path, err)
				errorCount++
				continue
			}
		}

		_, err = imageCollection.DeleteOne(context.TODO(), bson.M{"_id": image.ID})
		if err != nil {
			log.Printf("[CRON-ERROR] Failed deleting DB record for image ID %s: %v\n", image.ID.Hex(), err)
			errorCount++
			continue
		}

		log.Printf("[CRON] Successfully deleted image and DB record for: %s\n", image.Filename)
		deletedCount++
	}

	log.Printf("--- [CRON] Cleanup job finished. Deleted: %d. Errors: %d. ---\n", deletedCount, errorCount)
}
