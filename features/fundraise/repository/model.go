package repository

import (
	"context"
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"
	"raihpeduli/helpers"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db        *gorm.DB
	clStorage helpers.CloudStorageInterface
	collection *mongo.Collection
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface, collection *mongo.Collection) fundraise.Repository {
	return &model{
		db:        db,
		clStorage: clStorage,
		collection: collection,
	}
}

func (mdl *model) Paginate(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter) ([]fundraise.Fundraise, error) {
	var fundraises []fundraise.Fundraise

	offset := (pagination.Page - 1) * pagination.PageSize
	title := "%" + searchAndFilter.Title + "%"

	if err := mdl.db.Offset(offset).Limit(pagination.PageSize).
		Where("title LIKE ?", title).
		Where("target >= ?", searchAndFilter.MinTarget).
		Where("target <= ?", searchAndFilter.MaxTarget).
		Find(&fundraises).Error; err != nil {
		return nil, err
	}

	return fundraises, nil
}

func (mdl *model) PaginateMobile(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter) ([]fundraise.Fundraise, error) {
	var fundraises []fundraise.Fundraise

	offset := (pagination.Page - 1) * pagination.PageSize
	title := "%" + searchAndFilter.Title + "%"

	if err := mdl.db.Offset(offset).Limit(pagination.PageSize).
		Where("title LIKE ?", title).
		Where("target >= ?", searchAndFilter.MinTarget).
		Where("target <= ?", searchAndFilter.MaxTarget).
		Where("status = ?", "accepted").
		Find(&fundraises).Error; err != nil {
		return nil, err
	}

	return fundraises, nil
}

func (mdl *model) Insert(newFundraise fundraise.Fundraise) (int, error) {
	if err := mdl.db.Create(&newFundraise).Error; err != nil {
		return 0, err
	}

	return newFundraise.ID, nil
}

func (mdl *model) SelectByID(fundraiseID int) (*fundraise.Fundraise, error) {
	var fundraise fundraise.Fundraise

	if err := mdl.db.First(&fundraise, fundraiseID).Error; err != nil {
		return nil, err
	}

	return &fundraise, nil
}

func (mdl *model) TotalFundAcquired(fundraiseID int) (int32, error) {
	var totalAcquired int32

	mdl.db.Table("transactions").Select("sum(amount)").
		Where("fundraise_id = ?", fundraiseID).
		Where("status = 5").
		Row().
		Scan(&totalAcquired)

	return totalAcquired, nil
}

func (mdl *model) Update(fundraise fundraise.Fundraise) (int, error) {
	result := mdl.db.Table("fundraises").Updates(&fundraise)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (mdl *model) DeleteByID(fundraiseID int) (int, error) {
	result := mdl.db.Delete(&fundraise.Fundraise{}, fundraiseID)

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

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/" + objectName, nil
}

func (mdl *model) SelectBookmarkedFundraiseID(ownerID int) (map[int]string, error) {
	opts := options.Find().SetProjection(bson.M{"post_id": 1, "_id": 1})
	cursor, err := mdl.collection.Find(context.Background(), bson.M{"owner_id": ownerID, "post_type": "fundraise"}, opts)
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

func (mdl *model) SelectBookmarkByFundraiseAndOwnerID(fundraiseID, ownerID int) (string, error) {
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	var result bson.M
	if err := mdl.collection.FindOne(context.Background(), bson.M{"owner_id": ownerID, "post_id": fundraiseID, "post_type": "fundraise"}, opts).Decode(&result); err != nil {
		return "", err
	}

	objectIDString := result["_id"].(primitive.ObjectID).Hex()

	return objectIDString, nil
}

func (mdl *model) GetTotalData() int64 {
	var totalData int64

	result := mdl.db.Table("fundraises").Where("deleted_at IS NULL").Count(&totalData)

	if result.Error != nil {
		logrus.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataMobile() int64 {
	var totalData int64

	if err := mdl.db.Table("fundraises").
		Where("deleted_at IS NULL").
		Where("status = ?", "accepted").
		Count(&totalData).Error; err != nil {
			logrus.Error(err)
			return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64 {
	var totalData int64

	title := "%" + searchAndFilter.Title + "%"

	if err := mdl.db.Table("fundraises").
		Where("deleted_at IS NULL").
		Where("title LIKE ?", title).
		Where("target >= ?", searchAndFilter.MinTarget).
		Where("target <= ?", searchAndFilter.MaxTarget).
		Count(&totalData).Error; err != nil {
		logrus.Error(err)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataBySearchAndFilterMobile(searchAndFilter dtos.SearchAndFilter) int64 {
	var totalData int64

	title := "%" + searchAndFilter.Title + "%"

	if err := mdl.db.Table("fundraises").
		Where("deleted_at IS NULL").
		Where("title LIKE ?", title).
		Where("target >= ?", searchAndFilter.MinTarget).
		Where("target <= ?", searchAndFilter.MaxTarget).
		Where("status = ?", "accepted").
		Count(&totalData).Error; err != nil {
		logrus.Error(err)
		return 0
	}

	return totalData
}