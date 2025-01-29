package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//**DOMAIN** se usa para CREAR UN STRUCT que REPRESENTARÁ LA BBDD (ya que usaremos GORM que se encargará de hacer query en la BBDD)

// Definiremos el struc de tupo User
// Si nos fijamos es parecido al struct "CreateRequest" de endpoint.go pero tiene ID y creatdAt y UpdatedAt
// GORM AUTOMATICAMENTE cuando se incerte un registro a la bbdd, de setear el create y el updated
// Existe label para definir propiedades en gorm (bbdd) -> https://gorm.io/docs/models.html
// Tenemos NUESTRO DOMINIO LISTO PARA GORM
type User struct {
	ID        string         `json:"id" gorm:"type:char(36);not null;primaryKey;uniqueIndex"` //Char 36 (pq uuid), no nulo, PK, unico
	FirstName string         `json:"first_name" gorm:"type:char(50);not null"`                //Tambien definiremos maximo 50 de char
	LastName  string         `json:"last_name" gorm:"type:char(50);not null"`
	Email     string         `json:"email" gorm:"type:char(50);not null"`
	Phone     string         `json:"phone" gorm:"type:char(30);not null"`
	CreatedAt time.Time      `json:"-"` // `json:"-"` PARA QUE NO se incluya este campo en las respuestas JSON
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"` //Este campo es para que se haga un SOFTDELETE en vez de un hard delete
}

// Usaremos los hooks de GORM
// Estos se ejecutaran DE MANERA AUTOAMTICA ANTES DE CADA ACCION, segun corresponda
// Por ejemplo BeforeCreate SE EJECUTARA SIEMPRE ANTES DE una CRAECION DE USARIO de manera automatica (sin llaamrla)
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
