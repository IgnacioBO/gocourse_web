package enrollment

import (
	"errors"
	"log"

	"github.com/IgnacioBO/gocourse_web/internal/course"
	"github.com/IgnacioBO/gocourse_web/internal/domain"
	"github.com/IgnacioBO/gocourse_web/internal/user"
)

type Service interface {
	Create(userID, courseID string) (*domain.Enrollment, error)
}

type service struct {
	log       *log.Logger
	repo      Repository
	userSrv   user.Service // Agregamos userSrv y CourseServ para poder validar que existan antes de asociarlos
	courseSrv course.Service
}

func NewService(log *log.Logger, userSrv user.Service, courseSrv course.Service, repo Repository) Service {
	return &service{
		log:       log,
		userSrv:   userSrv, // Agregamos userSrv y CourseServ para poder validar que existan antes de asociarlos
		courseSrv: courseSrv,
		repo:      repo,
	}
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	s.log.Println("Create enrollment service")

	enrollmentNuevo := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "P",
	}

	//Ahora llamamos al GET de curso y el de Course para OBTENER y ver SI EXISTEN
	//Sino existe dara ERROR, por lo que retornamos un error
	_, err := s.userSrv.Get(enrollmentNuevo.UserID)
	if err != nil {
		return nil, errors.New("user_id doesn't exist")
	}

	if _, err := s.courseSrv.Get(enrollmentNuevo.CourseID); err != nil {
		return nil, errors.New("course_id doesn't exist")

	}

	//Le pasamo al repo el domain.Course (del domain.go) a la capa repo a la funcion Create (que recibe puntero)
	err = s.repo.Create(enrollmentNuevo)
	//Si hay un error (por ejemplo al insertar, se devuelve el error y la capa endpoitn lo maneja con un status code y todo)
	if err != nil {
		return nil, err
	}
	return enrollmentNuevo, nil
}
