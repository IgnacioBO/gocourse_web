package domain

//Enrollent sera una especia de tabla intermedia entre course y user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID        string     `json:"id" gorm:"type:char(36);not null;primaryKey;uniqueIndex"` //Char 36 (pq uuid), no nulo, PK, unico
	UserID    string     `json:"user_id,omitempty" gorm:"type:char(36)"`                  //Este será el forignKey, como se pone "NombreEntidadID" y luego abajo la *Entidad, Gorm de manera automatica sabe que UserID es clave foranea de User
	User      *User      `json:"user,omitempty"`                                          //Puntero de user
	CourseID  string     `json:"course_id,omitempty" gorm:"type:char(36)"`                //Este será el forignKey
	Course    *Course    `json:"course,omitempty"`                                        //Puntero de course
	Status    string     `json:"status" gorm:"type:char(2)"`                              //Estado del curso
	CreatedAt *time.Time `json:"-"`                                                       // `json:"-"` PARA QUE NO se incluya este campo en las respuestas JSON
	UpdatedAt *time.Time `json:"-"`
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return
}
