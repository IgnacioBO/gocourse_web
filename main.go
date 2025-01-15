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
	"github.com/gorilla/mux"
)

func main() {
	//Generaremos un router usando gorilla/mux para generar un RUTEO (osea los paths y metodos)
	router := mux.NewRouter()

	userEnd := user.MakeEndpoints()

	//Ahora setearemos que cuando le pegemos a /users le pege a las funciones definidas en el controlador user
	//Con handlefunc decimos que cuando valla a /users se ejecute la funcion correspondiente (userEnd.Create, Get, etc)
	//Podemos PONER y ESPECIFICAR EL METODO (si se quiere), si intento pegarle con otro no soportado tirará error 405
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users", userEnd.Delete).Methods("DELETE")

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
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	} //Una manera de "ACORTAR todo esto es poner EN UNA LINEA" -> log.Fatal(srv.ListenAndServe())

}
