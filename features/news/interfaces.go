package news

import (
	"mime/multipart"
	"raihpeduli/features/news/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter) ([]News, error)
	Insert(newNews News) (*News, error)
	SelectByID(newsID int) (*News, error)
	Update(news News) error
	DeleteByID(newsID int) error
	UploadFile(file multipart.File) (string, error)
	DeleteFile(filename string) error
	SelectBookmarkedNewsID(ownerID int) (map[int]string, error)
	SelectBoockmarkByNewsAndOwnerID(newsID, ownerID int) (string, error)
	GetTotalData() int64
	GetTotalDataBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64
}

type Usecase interface {
	FindAll(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter, ownerID int) ([]dtos.ResNews, int64)
	FindByID(newsID, ownerID int) *dtos.ResNews
	Create(newNews dtos.InputNews, userID int, file multipart.File) (*dtos.ResNews, []string, error)
	Modify(newsData dtos.InputNews, file multipart.File, oldData dtos.ResNews) ([]string, error)
	Remove(newsID int, oldData dtos.ResNews) error
}

type Handler interface {
	GetNews() echo.HandlerFunc
	NewsDetails() echo.HandlerFunc
	CreateNews() echo.HandlerFunc
	UpdateNews() echo.HandlerFunc
	DeleteNews() echo.HandlerFunc
}
