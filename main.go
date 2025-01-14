package main

//Primero hacemos git init y generamos el go.mod -> go mod init
//Usaremos el package net/http.
//Pero usaremos desde ahora gorilla/mux
//usar go get github.com/gorilla/mux
//Se usa similar que con net nativo (es decir usan w http.ResponseWriter, r *http.Request )
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	//Generaremos un router usando gorilla/mux para generar un RUTEO (osea los paths y metodos)
	//Ruteos lo haremos con gorilla/mux y no de manera nativa
	router := mux.NewRouter()

	//Ahora setearemos que cuando entremos a users le pege a getusers y a courses a getcousres
	//Con handlefunc decimos que cuando valla a /users se ejecute la funcion getUsers
	//Podemos PONER y ESPECIFICAR EL METODO (si se quiere), si intento pegarle con otro no soportado tirará error 405
	router.HandleFunc("/users", getUsers).Methods("GET", "POST")
	router.HandleFunc("/courses", getCourses).Methods("GET")

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

// Dos controladores iniciales
// w http.ResponseWriter -> Para enviar RESPUESTA AL CLIENTE (body, headers, statsu code)
// r *http.Request -> Contiende info de la SOLICITUD/REQUEST del cliente (aaceder al metodo hhttp (GET,POST,ETC), paarametros url/query string, body, headers, etc
// *http.Request siempre como PUNTERO (*) por mas eficiencia, para poder modificar datos y es el estandar del package net/http
func getUsers(w http.ResponseWriter, r *http.Request) {
	//Un ejemplo de delay
	time.Sleep(10 * time.Second)
	fmt.Println("get /users")
	//Para responder usaremos el paquete json
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})

}

func getCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get /courses")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})

}
