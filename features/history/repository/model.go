package repository

import (
	"context"
	"raihpeduli/features/history"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db         *gorm.DB
	collection *mongo.Collection
}

func New(db *gorm.DB, collection *mongo.Collection) history.Repository {
	return &model{
		db:         db,
		collection: collection,
	}
}

func (mdl *model) HistoryFundraiseCreatedByUser(userID int) ([]history.Fundraise, error) {
	var fundraise []history.Fundraise

	if err := mdl.db.Where("user_id = ? ", userID).Find(&fundraise).Error; err != nil {
		return nil, err
	}
	return fundraise, nil
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

func (mdl *model) HistoryVolunteerVacanciesCreatedByUser(userID int) ([]history.VolunteerVacancies, error) {
	var volunteer_vacancies []history.VolunteerVacancies

	if err := mdl.db.Where("user_id = ? ", userID).Find(&volunteer_vacancies).Error; err != nil {
		return nil, err
	}
	return volunteer_vacancies, nil
}

func (mdl *model) HistoryVolunteerVacanciesRegisterByUser(userID int) ([]history.Volunteer, error) {
	var volunteers []history.Volunteer

	if err := mdl.db.Table("volunteer_relations AS vr").
		Select("vr.id", "users.email", "users.fullname", "users.address", "users.phone_number", "users.gender", "users.nik", "vr.skills", "vr.reason", "vr.resume", "vr.status", "vr.photo").Joins("JOIN users ON users.id = vr.user_id").Where("vr.user_id = ? ", userID).Find(&volunteers).Error; err != nil {
		return nil, err
	}

	// if err := mdl.db.Joins("JOIN volunteer_relations ON volunteer_relations.volunteer_id = volunteer_vacancies.id").Where("volunteer_relations.user_id = ?", userID).Find(&VolunteerRelations).Error; err != nil {
	// 	return nil, err
	// }
	return volunteers, nil
}

func (mdl *model) SelectBookmarkedVacancyID(ownerID int) (map[int]string, error) {
	opts := options.Find().SetProjection(bson.M{"post_id": 1, "_id": 1})
	cursor, err := mdl.collection.Find(context.Background(), bson.M{"owner_id": ownerID, "post_type": "vacancy"}, opts)
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

func (mdl *model) GetTotalVolunteersByVacancyID(vacancyID int) int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_relations").Where("volunteer_id = ?", vacancyID).Count(&totalData)
	if result.Error != nil {
		logrus.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) HistoryUserTransaction(userID int) ([]history.Transaction, error) {
	var transaction []history.Transaction

	if err := mdl.db.Preload("User").Table("transactions").Joins("JOIN users ON transactions.user_id = users.id").Where("transactions.user_id = ?", userID).Find(&transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}
