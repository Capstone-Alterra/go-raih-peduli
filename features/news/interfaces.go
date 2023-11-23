package news

import (
	"mime/multipart"
	"raihpeduli/features/news/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, keyword string) ([]News, error)
	Insert(newNews News) (int, error)
	SelectByID(newsID int) (*News, error)
	Update(news News) (int, error)
	DeleteByID(newsID int) (int, error)
	UploadFile(file multipart.File, objectName string) (string, error)
	SelectBookmarkedNewsID(ownerID int) (map[int]string, error)
	SelectBoockmarkByNewsAndOwnerID(newsID, ownerID int) (string, error)
}

type Usecase interface {
	FindAll(page, size int, keyword string, ownerID int) []dtos.ResNews
	FindByID(newsID, ownerID int) *dtos.ResNews
	Create(newNews dtos.InputNews, userID int, file multipart.File) (*dtos.ResNews, []string, error)
	Modify(newsData dtos.InputNews, file multipart.File, oldData dtos.ResNews) bool
	Remove(newsID int) bool
}

type Handler interface {
	GetNews() echo.HandlerFunc
	NewsDetails() echo.HandlerFunc
	CreateNews() echo.HandlerFunc
	UpdateNews() echo.HandlerFunc
	DeleteNews() echo.HandlerFunc
}
