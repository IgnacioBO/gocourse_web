package user

import "log"

//**Capa service o business layer**
//Parecido a la capa endpoint
//Crearemos una interface llamda Service
//En la capa controlador (endpoint) manejamos con struct
//Pero en capa sevicio y capa repositorio SE MANEJRA CON INTERFACE -> porque es mas facil mockearlo o utilizarlo de manera mas generica

//Aqui definiremos lo metodos que las struct deberan tener
type Service interface {
	Create(firstName, lastName, email, phone string) error //Metodo que recibira datos de creacion y devolvera un error
}

//Ahora crearemos un struct PRIVADA (pq desde afuera accederemoa a traves de Servivce)
type service struct {
}

//Crea un servicio que sera la interfaz (devovlerá una interface de tupo Service [creado arriba], PERO hara un RETURN especificamente del STRUCT service (con minusculas))
func NewService() Service {
	return &service{}
}

//Crearemos un metodo Create que será de la struct service (OJO NO CONFUNDIR CON EL INTERFACE)
func (s service) Create(firstName, lastName, email, phone string) error {
	log.Println("Create user service")
	return nil
}
