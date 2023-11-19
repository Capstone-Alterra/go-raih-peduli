package usecase

import (
	"errors"
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
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

// func (svc *service) FindByID(bookmarkID int) *dtos.ResBookmark {
// 	var res dtos.ResBookmark

// 	bookmark := svc.model.SelectByID(bookmarkID)

// 	if bookmark == nil {
// 		return nil
// 	}

// 	if err := smapping.FillStruct(&res, smapping.MapFields(bookmark)); err != nil {
// 		logrus.Error(err)
// 		return nil
// 	}

// 	return &res
// }

func (svc *service) SetBookmark(input dtos.InputBookmarkPost, userID int) (bool, error) {
	var post any
	var bookmarkPost any

	switch(input.PostType) {
		case "news":
			post, _ = svc.model.SelectNewsByID(input.PostID)
			bookmarkPost = bookmark.NewsBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
			}
		case "fundraise":
			post, _ = svc.model.SelectFundraiseByID(input.PostID)
			bookmarkPost = bookmark.FundraiseBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
			}
		case "vacancy":
			post, _ = svc.model.SelectVolunteerByID(input.PostID)
			bookmarkPost = bookmark.VacancyBookmark{
				PostID: input.PostID,
				PostType: input.PostType,
			}
		default: 
			return false, errors.New("unknown post type. choose between 'news', 'fundraise' or 'vacancy'")
	}

	if post == nil {
		return false, errors.New("post not found")
	}

	if err := smapping.FillStruct(post, smapping.MapFields(post)); err != nil {
		return false, err
	}
	
	_, err := svc.model.Insert(bookmarkPost)
	
	if err != nil {
		return false, err
	}

	return true, nil
}

func (svc *service) UnsetBookmark(bookmarkID int) bool {
	_, err := svc.model.DeleteByID(bookmarkID)

	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}