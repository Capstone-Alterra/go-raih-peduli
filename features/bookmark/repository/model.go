package repository

import (
	"context"
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
	mongoDB *mongo.Database
}

func New(db *gorm.DB, mongoDB *mongo.Database) bookmark.Repository {
	return &model {
		db: db,
		mongoDB: mongoDB,
	}
}

func (mdl *model) Paginate(size, userID int) (*dtos.ResBookmark, error) {
	var bookmarkedPosts dtos.ResBookmark

	opts := options.Find().SetLimit(int64(size))

	cursor, err := mdl.mongoDB.Collection("bookmark").Find(context.Background(), bson.M{"user_id": userID, "post_type": "fundraise"}, opts)
	
	var fundraises []dtos.ResFundraise
	if err = cursor.All(context.TODO(), &fundraises); err != nil {
		return nil, err
	}
	
	cursor, err = mdl.mongoDB.Collection("bookmark").Find(context.Background(), bson.M{"user_id": userID, "post_type": "news"}, opts)
	
	var news []dtos.ResNews
	if err = cursor.All(context.TODO(), &news); err != nil {
		return nil, err
	}
	
	cursor, err = mdl.mongoDB.Collection("bookmark").Find(context.Background(), bson.M{"user_id": userID, "post_type": "vacancy"}, opts)

	var vacancies []dtos.ResVolunteerVacancy
	if err = cursor.All(context.TODO(), &vacancies); err != nil {
		return nil, err
	}

	bookmarkedPosts.Fundraise = fundraises
	bookmarkedPosts.News = news
	bookmarkedPosts.Vacancy = vacancies

	return &bookmarkedPosts, nil
}

func (mdl *model) Insert(document any) (bool, error) {
	if _, err := mdl.mongoDB.Collection("bookmark").InsertOne(context.Background(), document); err != nil {
		return false, err
	}

	return true, nil
}

func (mdl *model) SelectByID(bookmarkID int) (any, error) {
	var result any
	if err := mdl.mongoDB.Collection("bookmark").FindOne(context.TODO(), bson.M{"_id": bookmarkID}).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (mdl *model) SelectFundraiseByID(fundraiseID int) (*bookmark.Fundraise, error) {
	var fundraise bookmark.Fundraise

	if err := mdl.db.First(&fundraise, fundraiseID).Error; err != nil {
		return nil, err
	}

	return &fundraise, nil
}

func (mdl *model) SelectNewsByID(newsID int) (*bookmark.News, error) {
	var news bookmark.News

	if err := mdl.db.First(&news, newsID).Error; err != nil {
		return nil, err
	}

	return &news, nil
}

func (mdl *model) SelectVolunteerByID(volunteerID int) (*bookmark.VolunteerVacancy, error) {
	var volunteer bookmark.VolunteerVacancy

	if err := mdl.db.First(&volunteer, volunteerID).Error; err != nil {
		return nil, err
	}

	return &volunteer, nil
}

func (mdl *model) DeleteByID(bookmarkID int) (int, error) {
	result := mdl.mongoDB.Collection("bookmark").FindOneAndDelete(context.Background(), bson.M{"_id": bookmarkID})
	
	if result.Err() != nil {
		return 0, result.Err()
	}

	return 1, nil
}