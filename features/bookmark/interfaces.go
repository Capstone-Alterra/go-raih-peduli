package bookmark

import (
	"raihpeduli/features/bookmark/dtos"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Paginate(size, ownerID int) (*dtos.ResBookmark, error)
	Insert(document any) (bool, error)
	SelectByID(bookmarkID string) (*bson.M, error)
	SelectByPostAndOwnerID(postID int, ownerID int, postType string) (*bson.M, error) 
	SelectFundraiseByID(fundraiseID int) (*Fundraise, error)
	SelectNewsByID(newsID int) (*News, error)
	SelectVolunteerByID(volunteerID int) (*VolunteerVacancy, error)
	DeleteByID(bookmarkID string) (int, error)
}

type Usecase interface {
	FindAll(size, userID int) *dtos.ResBookmark
	FindByID(bookmarkID string) *bson.M
	SetBookmark(input dtos.InputBookmarkPost, ownerID int) (bool, []string, error)
	UnsetBookmark(bookmarkID string, bookmark *primitive.M, ownerID int) (bool, error)
}

type Handler interface {
	GetBookmarksByUserID() echo.HandlerFunc
	BookmarkAPost() echo.HandlerFunc
	UnBookmarkAPost() echo.HandlerFunc
}
