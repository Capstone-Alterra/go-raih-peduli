package usecase

import (
	"errors"
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type service struct {
	model bookmark.Repository
}

func New(model bookmark.Repository) bookmark.Usecase {
	return &service {
		model: model,
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

func (svc *service) SetBookmark(input dtos.InputBookmarkPost, userID int) (bool, error) {
	var post any
	var document any
	var err error

	bookmarked, _ := svc.model.SelectByPostAndUserID(input.PostID, userID, input.PostType)
	
	if bookmarked != nil {
		return false, errors.New("this post has already been bookmarked by this user")
	}

	switch(input.PostType) {
		case "news":
			post, err = svc.model.SelectNewsByID(input.PostID)
			document = &bookmark.NewsBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
				UserID: userID,
			}
		case "fundraise":
			post, err = svc.model.SelectFundraiseByID(input.PostID)
			document = &bookmark.FundraiseBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
				UserID: userID,
			}
		case "vacancy":
			post, err = svc.model.SelectVolunteerByID(input.PostID)
			document = &bookmark.VacancyBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
				UserID: userID,
			}
		default: 
			return false, errors.New("unknown post type. choose between 'news', 'fundraise' or 'vacancy'")
	}

	if err != nil {
		return false, err
	}
	
	if err := smapping.FillStruct(document, smapping.MapFields(post)); err != nil {
		return false, err
	}

	_, err = svc.model.Insert(document)
	
	if err != nil {
		return false, err
	}

	return true, nil
}

func (svc *service) UnsetBookmark(bookmarkID string) bool {
	_, err := svc.model.DeleteByID(bookmarkID)

	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}