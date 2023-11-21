package repository

import (
	"context"
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/news"
	"raihpeduli/helpers"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db         *gorm.DB
	clStorage  helpers.CloudStorageInterface
	collection *mongo.Collection
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface, collection *mongo.Collection) news.Repository {
	return &model{
		db:         db,
		clStorage:  clStorage,
		collection: collection,
	}
}

func (mdl *model) Paginate(page, size int, keyword string) ([]news.News, error) {
	var news []news.News

	offset := (page - 1) * size
	searching := "%" + keyword + "%"

	if err := mdl.db.Offset(offset).Limit(size).Where("title LIKE ?", searching).Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

func (mdl *model) Insert(newNews news.News) (int, error) {
	if err := mdl.db.Create(&newNews).Error; err != nil {
		return 0, err
	}

	return newNews.ID, nil
}

func (mdl *model) SelectByID(newsID int) (*news.News, error) {
	var news news.News

	if err := mdl.db.First(&news, newsID).Error; err != nil {
		return nil, err
	}

	return &news, nil
}

func (mdl *model) Update(news news.News) (int, error) {
	result := mdl.db.Save(&news)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (mdl *model) DeleteByID(newsID int) (int, error) {
	result := mdl.db.Delete(&news.News{}, newsID)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (mdl *model) UploadFile(file multipart.File, objectName string) (string, error) {
	config := config.LoadCloudStorageConfig()
	randomChar := uuid.New().String()
	if objectName == "" {
		objectName = randomChar
	}

	if err := mdl.clStorage.UploadFile(file, objectName); err != nil {
		return "", err
	}

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/news/" + objectName, nil
}

func (mdl *model) SelectBookmarkedNewsID(ownerID int) (map[int]string, error) {
	opts := options.Find().SetProjection(bson.M{"post_id": 1, "_id": 1})
	cursor, err := mdl.collection.Find(context.Background(), bson.M{"owner_id": ownerID, "post_type": "news"}, opts)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	var mapPostIDs = map[int]string{}

	for _, data := range results {
		postID := int(data["post_id"].(int32))
		mapPostIDs[postID] = data["_id"].(primitive.ObjectID).Hex()
	}

	return mapPostIDs, nil
}

func (mdl *model) SelectBoockmarkByNewsAndOwnerID(newsID, ownerID int) (string, error) {
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	var result bson.M
	if err := mdl.collection.FindOne(context.Background(), bson.M{"owner_id": ownerID, "post_id": newsID, "post_type": "news"}, opts).Decode(&result); err != nil {
		return "", err
	}

	objectIDString := result["_id"].(primitive.ObjectID).Hex()

	return objectIDString, nil
}
