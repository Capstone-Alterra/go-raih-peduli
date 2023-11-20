package repository

import (
	"context"
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
	collection *mongo.Collection
}

func New(db *gorm.DB, collection *mongo.Collection) bookmark.Repository {
	return &model {
		db: db,
		collection: collection,
	}
}

func (mdl *model) Paginate(size, userID int) (*dtos.ResBookmark, error) {
	var bookmarkedPosts dtos.ResBookmark

	opts := options.Find().SetLimit(int64(size))

	cursor, err := mdl.collection.Find(context.Background(), bson.M{"user_id": userID, "post_type": "fundraise"}, opts)
	logrus.Error(err)
	
	var fundraises []dtos.ResFundraise
	if err = cursor.All(context.TODO(), &fundraises); err != nil {
		logrus.Error(err)
		return nil, err
	}
	
	cursor, err = mdl.collection.Find(context.Background(), bson.M{"user_id": userID, "post_type": "news"}, opts)
	logrus.Error(err)
	
	var news []dtos.ResNews
	if err = cursor.All(context.TODO(), &news); err != nil {
		logrus.Error(err)
		return nil, err
	}
	
	cursor, err = mdl.collection.Find(context.Background(), bson.M{"user_id": userID, "post_type": "vacancy"}, opts)
	logrus.Error(err)
	
	var vacancies []dtos.ResVolunteerVacancy
	if err = cursor.All(context.TODO(), &vacancies); err != nil {
		logrus.Error(err)
		return nil, err
	}

	bookmarkedPosts.Fundraise = fundraises
	bookmarkedPosts.News = news
	bookmarkedPosts.Vacancy = vacancies

	return &bookmarkedPosts, nil
}

func (mdl *model) Insert(document any) (bool, error) {
	if _, err := mdl.collection.InsertOne(context.Background(), document); err != nil {
		return false, err
	}

	return true, nil
}

func (mdl *model) SelectByPostAndUserID(postID int, userID int, postType string) (*bson.M, error) {
	var result bson.M

	if err := mdl.collection.FindOne(context.Background(), bson.M{"user_id": userID, "post_id": postID, "post_type": postType}).Decode(&result); err != nil {
		return nil, err
	}

	logrus.Info(result)

	return &result, nil
}

func (mdl *model) SelectByID(bookmarkID string) (*bson.M, error) {
	var result bson.M

	objectID, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return nil, err
	}

	if err := mdl.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
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

func (mdl *model) DeleteByID(bookmarkID string) (int, error) {
	objectID, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return 0, err
	}

	result := mdl.collection.FindOneAndDelete(context.Background(), bson.M{"_id": objectID})
	
	if result.Err() != nil {
		return 0, result.Err()
	}

	return 1, nil
}