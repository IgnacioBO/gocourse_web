package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID        string         `json:"id" gorm:"type:char(36);not null;primaryKey;uniqueIndex"` //Char 36 (pq uuid), no nulo, PK, unico
	Name      string         `json:"name" gorm:"type:char(50);not null"`                      //Tambien definiremos maximo 50 de char
	StartDate time.Time      `json:"start_date"`                                              //no es necesario poner el tipo? lo identifia solo?
	EndDate   time.Time      `json:"end_date"`
	CreatedAt *time.Time     `json:"-"` // `json:"-"` PARA QUE NO se incluya este campo en las respuestas JSON
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"` //Este campo es para que se haga un SOFTDELETE en vez de un hard delete
}

//Que pasa si quiero agregar dentro del Struct Course poner un User user.User? Y dentro del struct User poner un course.Course?
//No nos dejara importar, no podemos importar package b en package a y a la vez importar package a en package b
//Por eso debemos oensar en PATRONES de DISEÑOS pensado en package.
//Pensar evitando estos problemas
//Idealmente mejor generar un package Domain con TODOS LOS DOMINIOS DE nuestra aplicacion
//ENTONCES COMO User y Course, ambos estan dentro del package domain se pueden llamar una a la otra

func (c *Course) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
