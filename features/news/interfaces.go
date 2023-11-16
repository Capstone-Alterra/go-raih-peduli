package news

import (
	"mime/multipart"
	"raihpeduli/features/news/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, keyword string) ([]News, error)
	Insert(newNews News) (int, error)
	SelectByID(newsID int) *News
	Update(news News) int64
	DeleteByID(newsID int) int64
	UploadFile(file multipart.File, objectName string) (string, error)
}

type Usecase interface {
	FindAll(page, size int, keyword string) []dtos.ResNews
	FindByID(newsID int) *dtos.ResNews
	Create(newNews dtos.InputNews, userID int, file multipart.File) (*dtos.ResNews, error)
	Modify(newsData dtos.InputNews, newsID int) bool
	Remove(newsID int) bool
}

type Handler interface {
	GetNews() echo.HandlerFunc
	NewsDetails() echo.HandlerFunc
	CreateNews() echo.HandlerFunc
	UpdateNews() echo.HandlerFunc
	DeleteNews() echo.HandlerFunc
}
