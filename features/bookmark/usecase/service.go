package usecase

import (
	"errors"
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"
	"raihpeduli/helpers"

	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	model bookmark.Repository
	validation helpers.ValidationInterface
}

func New(model bookmark.Repository, validation helpers.ValidationInterface) bookmark.Usecase {
	return &service {
		model: model,
		validation: validation,
	}
}

func (svc *service) FindAll(size, userID int) *dtos.ResBookmark {
	bookmarks, err := svc.model.Paginate(size, userID)

	if err != nil {
		return nil
	}

	return bookmarks
}

func (svc *service) FindByID(bookmarkID string) *bson.M {
	bookmark, err := svc.model.SelectByID(bookmarkID)

	if err != nil {
		return nil
	}

	return bookmark
}

func (svc *service) SetBookmark(input dtos.InputBookmarkPost, ownerID int) (bool, []string, error) {
	var post any
	var document any
	var err error

	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		return false, errMap, errors.New("error")
	}

	bookmarked, _ := svc.model.SelectByPostAndOwnerID(input.PostID, ownerID, input.PostType)
	
	if bookmarked != nil {
		return false, nil, errors.New("this post has already been bookmarked by this user")
	}

	switch(input.PostType) {
		case "news":
			post, err = svc.model.SelectNewsByID(input.PostID)
			document = &bookmark.NewsBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
				OwnerID: ownerID,
			}
		case "fundraise":
			post, err = svc.model.SelectFundraiseByID(input.PostID)
			document = &bookmark.FundraiseBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
				OwnerID: ownerID,
			}
		case "vacancy":
			post, err = svc.model.SelectVolunteerByID(input.PostID)
			document = &bookmark.VacancyBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
				OwnerID: ownerID,
			}
		default: 
			return false, nil, errors.New("unknown post type. choose between 'news', 'fundraise' or 'vacancy'")
	}

	if err != nil {
		return false, nil, err
	}
	
	if err := smapping.FillStruct(document, smapping.MapFields(post)); err != nil {
		return false, nil, err
	}

	_, err = svc.model.Insert(document)
	
	if err != nil {
		return false, nil, err
	}

	return true, nil, nil
}

func (svc *service) UnsetBookmark(bookmarkID string, bookmark *primitive.M, ownerID int) (bool, error) {
	bsonData, _ := bson.Marshal(bookmark)

	var result dtos.ResOwnerID
	bson.Unmarshal(bsonData, &result)

	if result.OwnerID != ownerID {
		return false, errors.New("this user can't unbookmark this bookmarked post because of a mistmach of user id")
	}

	_, err := svc.model.DeleteByID(bookmarkID)

	if err != nil {
		logrus.Error(err)
		return false, err
	}

	return true, nil
}