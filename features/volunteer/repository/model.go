package repository

import (
	"context"
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"
	"raihpeduli/helpers"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db         *gorm.DB
	clStorage  helpers.CloudStorageInterface
	collection *mongo.Collection
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface, collection *mongo.Collection) volunteer.Repository {
	return &model{
		db:         db,
		clStorage:  clStorage,
		collection: collection,
	}
}

func (mdl *model) Paginate(page, size int, searchAndFilter dtos.SearchAndFilter) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).
		Where("title LIKE ?", "%"+searchAndFilter.Title+"%").
		Where("city LIKE ?", "%"+searchAndFilter.City+"%").
		Where("skills_required LIKE ?", "%"+searchAndFilter.Skill+"%").
		Where("number_of_vacancies >= ?", searchAndFilter.MinParticipant).
		Where("number_of_vacancies <= ?", searchAndFilter.MaxParticipant).
		Find(&volunteers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}

func (mdl *model) PaginateMobile(page, size int, searchAndFilter dtos.SearchAndFilter) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).
		Where("title LIKE ?", "%"+searchAndFilter.Title+"%").
		Where("city LIKE ?", "%"+searchAndFilter.City+"%").
		Where("skills_required LIKE ?", "%"+searchAndFilter.Skill+"%").
		Where("number_of_vacancies >= ?", searchAndFilter.MinParticipant).
		Where("number_of_vacancies <= ?", searchAndFilter.MaxParticipant).
		Where("status = ?", "accepted").
		Find(&volunteers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}

func (mdl *model) SelectVacancyByID(volunteerID int) *volunteer.VolunteerVacancies {
	var volunteer volunteer.VolunteerVacancies
	result := mdl.db.First(&volunteer, volunteerID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &volunteer
}

func (mdl *model) SelectBookmarkedVacancyID(ownerID int) (map[int]string, error) {
	opts := options.Find().SetProjection(bson.M{"post_id": 1, "_id": 1})
	cursor, err := mdl.collection.Find(context.Background(), bson.M{"owner_id": ownerID, "post_type": "vacancy"}, opts)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	var mapPostIDs = map[int]string{}

	for _, data := range results {
		postID := int(data["post_id"].(int32))
		mapPostIDs[postID] = data["_id"].(primitive.ObjectID).Hex()
	}

	return mapPostIDs, nil
}

func (mdl *model) SelectBookmarkByVacancyAndOwnerID(vacancyID, ownerID int) string {
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	var result bson.M
	if err := mdl.collection.FindOne(context.Background(), bson.M{"owner_id": ownerID, "post_id": vacancyID, "post_type": "vacancy"}, opts).Decode(&result); err != nil {
		return ""
	}

	objectIDString := result["_id"].(primitive.ObjectID).Hex()

	return objectIDString
}

func (mdl *model) UpdateVacancy(volunteer volunteer.VolunteerVacancies) int64 {
	result := mdl.db.Updates(&volunteer)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteVacancyByID(volunteerID int) int64 {
	result := mdl.db.Delete(&volunteer.VolunteerVacancies{}, volunteerID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}

func (mdl *model) InsertVacancy(newVolunteer *volunteer.VolunteerVacancies) (*volunteer.VolunteerVacancies, error) {
	result := mdl.db.Create(newVolunteer)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}
	return newVolunteer, nil
}

func (mdl *model) RegisterVacancy(registrar *volunteer.VolunteerRelations) error {
	result := mdl.db.Table("volunteer_relations").Create(registrar)

	if result.Error != nil {
		log.Error(result.Error)
		return result.Error
	}
	return nil
}

func (mdl *model) SelectRegistrarByID(registrarID int) *volunteer.VolunteerRelations {
	var registrar volunteer.VolunteerRelations
	result := mdl.db.First(&registrar, registrarID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &registrar
}

func (mdl *model) UpdateStatusRegistrar(registrar volunteer.VolunteerRelations) int64 {
	result := mdl.db.Table("volunteer_relations").Updates(&registrar)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) UploadFile(file multipart.File, oldFilename string) (string, error) {
	var config = config.LoadCloudStorageConfig()
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/vacancies/")
	var objectName string

	if file == nil {
		return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/vacancies/volunteer-vacancy.jpg", nil
	}

	if oldFilename != "" {
		objectName = oldFilename[urlLength:]

		if objectName == "volunteer-vacancy.jpg" {
			objectName = ""
		} else if err := mdl.clStorage.DeleteFile(objectName); err != nil {
			return "", err
		}
	}
	objectName = uuid.New().String()

	if err := mdl.clStorage.UploadFile(file, objectName); err != nil {
		return "", err
	}

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/vacancies/" + objectName, nil
}

func (mdl *model) GetTotalDataVacancies() int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_vacancies").Where("deleted_at IS NULL").Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataVacanciesMobile() int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_vacancies").
		Where("deleted_at IS NULL").
		Where("status = ?", "accepted").
		Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataVacanciesBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_vacancies").
		Where("title LIKE ?", "%"+searchAndFilter.Title+"%").
		Where("city LIKE ?", "%"+searchAndFilter.City+"%").
		Where("skills_required LIKE ?", "%"+searchAndFilter.Skill+"%").
		Where("number_of_vacancies >= ?", searchAndFilter.MinParticipant).
		Where("number_of_vacancies <= ?", searchAndFilter.MaxParticipant).
		Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataVacanciesBySearchAndFilterMobile(searchAndFilter dtos.SearchAndFilter) int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_vacancies").
		Where("title LIKE ?", "%"+searchAndFilter.Title+"%").
		Where("city LIKE ?", "%"+searchAndFilter.City+"%").
		Where("skills_required LIKE ?", "%"+searchAndFilter.Skill+"%").
		Where("number_of_vacancies >= ?", searchAndFilter.MinParticipant).
		Where("number_of_vacancies <= ?", searchAndFilter.MaxParticipant).
		Where("status = ?", "accepted").
		Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalVolunteersByVacancyID(vacancyID int) int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_relations").Where("volunteer_id = ?", vacancyID).Count(&totalData)
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) SelectVolunteersByVacancyID(vacancyID int, name string, page, size int) []volunteer.Volunteer {
	var volunteers []volunteer.Volunteer

	offset := (page - 1) * size

	result := mdl.db.Table("volunteer_relations AS vr").
		Select("vr.id", "users.fullname", "users.address", "users.nik", "vr.resume", "vr.status", "vr.photo").
		Joins("JOIN users ON users.id = vr.user_id").
		Where("vr.volunteer_id = ?", vacancyID).
		Where("users.fullname LIKE ?", "%"+name+"%").
		Offset(offset).Limit(size).Find(&volunteers)
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}

func (mdl *model) SelectVolunteerDetails(vacancyID int, volunteerID int) *volunteer.Volunteer {
	var volunteers volunteer.Volunteer

	result := mdl.db.Table("volunteer_relations AS vr").
		Select("vr.id", "users.fullname", "users.address", "users.nik", "vr.resume", "vr.status", "vr.photo").
		Joins("JOIN users ON users.id = vr.user_id").
		Where("vr.volunteer_id = ? AND vr.id = ?", vacancyID, volunteerID).
		Find(&volunteers)
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &volunteers
}

func (mdl *model) GetTotalVolunteers(vacancyID int, name string) int64 {
	var totalData int64

	result := mdl.db.Table("volunteer_relations AS vr").Select("users.fullname", "users.address", "users.nik", "vr.resume", "vr.status").
		Joins("JOIN users ON users.id = vr.user_id").
		Where("vr.volunteer_id = ?", vacancyID).
		Where("users.fullname LIKE ?", "%"+name+"%").Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) CheckUser(userID int) bool {
	var nik string

	result := mdl.db.Table("users").Select("nik").Where("id = ?", userID).Pluck("nik", &nik)
	if result.Error != nil {
		return false
	}

	if nik == "" {
		return false
	}

	return true
}
