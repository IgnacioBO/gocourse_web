package enrollment

import (
	"log"

	"gorm.io/gorm"

	"github.com/IgnacioBO/gocourse_web/internal/domain"
)

type Repository interface {
	Create(e *domain.Enrollment) error
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db, //Devuevle un struct repo con la bbdd
	}

}

func (r *repo) Create(enrollment *domain.Enrollment) error {
	r.log.Println("repository Create:", enrollment)

	result := r.db.Create(enrollment)

	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}
	r.log.Printf("enrollment created with id: %s, rows affected: %d\n", enrollment.ID, result.RowsAffected)
	return nil
}
