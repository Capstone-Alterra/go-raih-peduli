package usecase

import (
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"github.com/labstack/gommon/log"
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

func (svc *service) FindAll(page, size int) []dtos.ResBookmark {
	var bookmarks []dtos.ResBookmark

	bookmarksEnt := svc.model.Paginate(page, size)

	for _, bookmark := range bookmarksEnt {
		var data dtos.ResBookmark

		if err := smapping.FillStruct(&data, smapping.MapFields(bookmark)); err != nil {
			log.Error(err.Error())
		} 
		
		bookmarks = append(bookmarks, data)
	}

	return bookmarks
}

func (svc *service) FindByID(bookmarkID int) *dtos.ResBookmark {
	var res dtos.ResBookmark

	bookmark := svc.model.SelectByID(bookmarkID)

	if bookmark == nil {
		return nil
	}

	if err := smapping.FillStruct(&res, smapping.MapFields(bookmark)); err != nil {
		logrus.Error(err)
		return nil
	}

	return &res

}

func (svc *service) SetBookmark(newBookmark dtos.InputBookmark) *dtos.ResBookmark {
	bookmark := bookmark.Bookmark{}
	
	err := smapping.FillStruct(&bookmark, smapping.MapFields(newBookmark))
	if err != nil {
		log.Error(err)
		return nil
	}

	bookmarkID := svc.model.Insert(bookmark)

	if bookmarkID == -1 {
		return nil
	}

	resBookmark := dtos.ResBookmark{}
	errRes := smapping.FillStruct(&resBookmark, smapping.MapFields(newBookmark))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resBookmark
}

func (svc *service) UnsetBookmark(bookmarkID int) bool {
	rowsAffected := svc.model.DeleteByID(bookmarkID)

	if rowsAffected <= 0 {
		log.Error("There is No Bookmark Deleted!")
		return false
	}

	return true
}