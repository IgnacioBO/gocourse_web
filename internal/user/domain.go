package user

import "time"

//**DOMAIN** se usa para CREAR UN STRUCT que REPRESENTARÁ LA BBDD (ya que usaremos GORM que se encargará de hacer query en la BBDD)

//Definiremos el struc de tupo User
//Si nos fijamos es parecido al struct "CreateRequest" de endpoint.go pero tiene ID y creatdAt y UpdatedAt
//GORM AUTOMATICAMENTE cuando se incerte un registro a la bbdd, de setear el create y el updated
//Existe label para definir propiedades en gorm (bbdd) -> https://gorm.io/docs/models.html
//Tenemos NUESTRO DOMINIO LISTO PARA GORM
type User struct {
	ID        string    `json:"id" gorm:"type:char(36);not null;primaryKey;uniqueIndex"` //Char 36 (pq uuid), no nulo, PK, unico
	FirstName string    `json:"first_name" gorm:"type:char(50);not null"`                //Tambien definiremos maximo 50 de char
	LastName  string    `json:"last_name" gorm:"type:char(50);not null"`
	Email     string    `json:"email" gorm:"type:char(50);not null"`
	Phone     string    `json:"phone" gorm:"type:char(30);not null"`
	CreatedAt time.Time `json:"-"` // `json:"-"` PARA QUE NO se incluya este campo en las respuestas JSON
	UpdatedAt time.Time `json:"-"`
}
