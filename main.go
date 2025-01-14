package main

//Primero hacemos git init y generamos el go.mod -> go mod init
//Usaremos el package net/http.
import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	port := ":3333"
	//Ahora setearemso que cuando entremos a users le pege a getusers y a courses a getcousres
	//Con handlefunc decimos que cuando valla a /users se ejecute la funcion getUsers
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/courses", getCourses)

	//Levantaremos un servidor y escucharemos el puerto deifnido
	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println(err)
	}
}

// Dos controladores iniciales
// w http.ResponseWriter -> Para enviar RESPUESTA AL CLIENTE (body, headers, statsu code)
// r *http.Request -> Contiende info de la SOLICITUD/REQUEST del cliente (aaceder al metodo hhttp (GET,POST,ETC), paarametros url/query string, body, headers, etc
// *http.Request siempre como PUNTERO (*) por mas eficiencia, para poder modificar datos y es el estandar del package net/http
func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get /users")
	//Modificarmeos el respone
	io.WriteString(w, "Endpoint users\n")

}

func getCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get /courses")
	io.WriteString(w, "Endpoint courses\n")

}
