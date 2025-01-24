package main

//Primero hacemos git init y generamos el go.mod -> go mod init
//Usaremos el package net/http.
//Pero usaremos desde ahora gorilla/mux
//usar go get github.com/gorilla/mux
//Se usa similar que con net nativo (es decir usan w http.ResponseWriter, r *http.Request )
import (
	"log"
	"net/http"
	"time"

	"github.com/IgnacioBO/gocourse_web/internal/user"
	"github.com/IgnacioBO/gocourse_web/pkg/bootsrap"
	"github.com/gorilla/mux" //Manejar ruteo facilmente (paths y metodos)
	"github.com/joho/godotenv"
	//Driver mysql para gorm
	//Para manejar bbdd facilmente con strcut y funciones (en vez de querys directa)
)

func main() {
	//Generaremos un router usando gorilla/mux para generar un RUTEO (osea los paths y metodos)
	router := mux.NewRouter()

	//Definiremos un logger
	//Creamos un script bootsrap.go que INICIALIZA el log
	l := bootsrap.InitLogger()

	//Con godotenv cargamos las variables de entorn denrto de .env para usarlas en el DSN
	//Al usar godotenv.Load() cargar autoamticametne los valores en archivo .env de root
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	//DSN (Data Source Name) es una cadena de conexion de BBDD (tipo, servidor, nombre bbdd, user, pass)
	//Creamos un script bootsrap.go que INICIALIZA la conexion usando gorm y varaibles de entoero
	db, err := bootsrap.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	//Generaremos un objeto repo (que recibe la bbdd y logger) que luego le pasaremos a la capa servicio
	userRepo := user.NewRepo(l, db)
	//Crearemos un objeto de tipo servicio pasandole un objeto Repository (y logger) para luego pasarselo a la capa enpdoint
	userService := user.NewService(l, userRepo)
	//Crearemo un objeto de tipo endpoint y le pasamos el objeto creado (Service)
	userEndpoint := user.MakeEndpoints(userService)

	//Ahora setearemos que cuando le pegemos a /users le pege a las funciones definidas en el controlador user
	//Con handlefunc decimos que cuando valla a /users se ejecute la funcion correspondiente (userEnd.Create, Get, etc)
	//Podemos PONER y ESPECIFICAR EL METODO (si se quiere), si intento pegarle con otro no soportado tirará error 405
	router.HandleFunc("/users", userEndpoint.Create).Methods("POST")
	router.HandleFunc("/users", userEndpoint.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoint.Update).Methods("PATCH")  //Usamos {} para especifiar el NOMBRE del paramatro (que se obtiene con MUX dentro de endpoint.go) //Patch es ACT parcial (PUT es completa)
	router.HandleFunc("/users/{id}", userEndpoint.Delete).Methods("DELETE") //Delete o SoftDelete
	router.HandleFunc("/users/{id}", userEndpoint.Get).Methods("GET")       //Usamos {} para especifiar el NOMBRE del paramatro (que se obtiene con MUX dentro de endpoint.go)

	//Levantaremos un servidor pero de distina manera a antes
	//err := http.ListenAndServe(port, nil)
	//Crearemos un objeto server y lo configuraremos
	//Handler sera el router
	//Addr le ponemos la ip y puerto
	srv := &http.Server{
		Handler: router,
		//Handler:    http.TimeoutHandler(router, time.Second*3, "Timeeout!"), //Usnado TimeoutHandler permite manejar timeout con mensaje (diferente al read y writetomiiut)
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second, //Con estos SETEAMOS TIMEOUT DE ESCRITURA Y DE LECTURA (cuanto timepo maximo la api permite)
		WriteTimeout: 5 * time.Second, // Read es REQUEST, WRITE es RESPONE
	}
	//Y ahora iniciamos el servidor
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	} //Una manera de "ACORTAR todo esto es poner EN UNA LINEA" -> log.Fatal(srv.ListenAndServe())

}
