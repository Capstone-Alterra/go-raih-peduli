package repository

import (
	"context"
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"
	"raihpeduli/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db           *gorm.DB
	clStorage    helpers.CloudStorageInterface
	collection   *mongo.Collection
	nsCollection *mongo.Collection
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface, collection *mongo.Collection, nsCollection *mongo.Collection) fundraise.Repository {
	return &model{
		db:           db,
		clStorage:    clStorage,
		collection:   collection,
		nsCollection: nsCollection,
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
		Order("created_at desc").
		Find(&fundraises).Error; err != nil {
		return nil, err
	}

	return fundraises, nil
}

func (mdl *model) PaginateMobile(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter) ([]fundraise.Fundraise, error) {
	var fundraises []fundraise.Fundraise

	offset := (pagination.Page - 1) * pagination.PageSize
	title := "%" + searchAndFilter.Title + "%"

	currentTimeUTC := time.Now()
	wibLocation, _ := time.LoadLocation("Asia/Jakarta")
	currentTimeWIB := currentTimeUTC.In(wibLocation)

	if err := mdl.db.Offset(offset).Limit(pagination.PageSize).
		Where("title LIKE ?", title).
		Where("target >= ?", searchAndFilter.MinTarget).
		Where("target <= ?", searchAndFilter.MaxTarget).
		Where("end_date > ?", currentTimeWIB.Format("2006-01-02 15:04:05")).
		Where("status = ?", "accepted").
		Find(&fundraises).Error; err != nil {
		return nil, err
	}

	return fundraises, nil
}

func (mdl *model) Insert(newFundraise fundraise.Fundraise) (*fundraise.Fundraise, error) {
	if err := mdl.db.Create(&newFundraise).Error; err != nil {
		return nil, err
	}

	return &newFundraise, nil
}

func (mdl *model) SelectByID(fundraiseID int) (*dtos.FundraiseDetails, error) {
	var fundraise dtos.FundraiseDetails

	if err := mdl.db.Table("fundraises").Select("fundraises.*, users.fullname as user_fullname, users.profile_picture as user_photo").Joins("LEFT JOIN users ON users.id = fundraises.user_id").
		First(&fundraise, fundraiseID).Error; err != nil {
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

func (mdl *model) Update(fundraise fundraise.Fundraise) error {
	result := mdl.db.Table("fundraises").Updates(&fundraise)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (mdl *model) DeleteByID(fundraiseID int) error {
	result := mdl.db.Delete(&fundraise.Fundraise{}, fundraiseID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (mdl *model) UploadFile(file multipart.File) (string, error) {
	config := config.LoadCloudStorageConfig()
	randomChar := uuid.New().String()
	filename := randomChar

	if err := mdl.clStorage.UploadFile(file, filename); err != nil {
		return "", err
	}

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/" + filename, nil
}

func (mdl *model) DeleteFile(filename string) error {
	if err := mdl.clStorage.DeleteFile(filename); err != nil {
		return err
	}

	return nil
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

	currentTimeUTC := time.Now()
	wibLocation, _ := time.LoadLocation("Asia/Jakarta")
	currentTimeWIB := currentTimeUTC.In(wibLocation)

	if err := mdl.db.Table("fundraises").
		Where("deleted_at IS NULL").
		Where("title LIKE ?", title).
		Where("target >= ?", searchAndFilter.MinTarget).
		Where("target <= ?", searchAndFilter.MaxTarget).
		Where("end_date > ?", currentTimeWIB.Format("2006-01-02 15:04:05")).
		Where("status = ?", "accepted").
		Count(&totalData).Error; err != nil {
		logrus.Error(err)
		return 0
	}

	return totalData
}

func (mdl *model) SelectByTitle(title string) (*fundraise.Fundraise, error) {
	var fundraise fundraise.Fundraise

	if err := mdl.db.Where("title = ?", title).
		Where("status = ?", "accepted").First(&fundraise).Error; err != nil {
		return nil, err
	}

	return &fundraise, nil
}

func (mdl *model) GetDeviceToken(userID int) string {
	var result fundraise.NotificationToken

	if err := mdl.nsCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&result); err != nil {
		logrus.Error(err)
		return ""
	}

	return result.DeviceToken
}
