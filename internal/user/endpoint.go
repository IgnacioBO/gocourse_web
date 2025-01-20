package user

//**Capa endpoint o controlador**

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct que tenga todos los endpoints que vayamos a utilizar
// Que teng una fucion que recibe un request y un response
type (
	//Controller sera una funcion que reciba REspone y Request
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller //Esto es lo mismo que decir Create func(w http.ResponseWriter, r *http.Request), pero como TODOS SON tipo Controller (Definido arriba) nos ahorramos ahcerlo
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}
	//Definiremos una struct para definir el request del Craete, con los campos que quiero recibir y los tags de json
	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

// Funcion que se encargará de hacer los endopints
// Para eso necesitaremos una struct que se llamara endpoints
// Esta funcion va a DEVOLVER una struct de Endpoints, estos endpoints son los que vamos a poder utuaizlar en unestro dominio (user)
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}

// Este devolver un Controller, retora una función de tipo Controller (que definimos arriba) con esta caractesitica
// Es privado porque se llamar solo de este dominio
func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("delete user")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

// w http.ResponseWriter -> Para enviar RESPUESTA AL CLIENTE (body, headers, statsu code)
// r *http.Request -> Contiende info de la SOLICITUD/REQUEST del cliente (aaceder al metodo http (GET,POST,ETC), paarametros url/query string, body, headers, etc
// *http.Request siempre como PUNTERO (*) por mas eficiencia, para poder modificar datos y es el estandar del package net/http
func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("create user")

		//Variable con struct de request (datos usaurio)
		var reqStruct CreateRequest
		//r.Body tiene el body del request (se espera JSON) y lo decodifica al struct (reqStruct) (osea pasar el json enviado en el request a un struct)
		err := json.NewDecoder(r.Body).Decode(&reqStruct)
		if err != nil {
			//w.WriteHeader devuelve en el repsonse el CODE que se le indica
			w.WriteHeader(400)
			//Enviaremos la repsuesta con encode y creamos un Sruct ErrorRespone (Creado antes) con un texto
			json.NewEncoder(w).Encode(ErrorResponse{"invalid request format"})
			return
		}

		//Validaciones
		if reqStruct.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"first_name is required"})
			return
		}
		if reqStruct.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"last_name is required"})
			return
		}
		fmt.Println(reqStruct)
		reqStrucEnJson, _ := json.MarshalIndent(reqStruct, "", " ")
		fmt.Println(string(reqStrucEnJson))

		//Usaremos la s recibida como parametro (de la capa Service y usaremos el metodo CREATE con lo que debe recibir)
		err = s.Create(reqStruct.FirstName, reqStruct.LastName, reqStruct.Email, reqStruct.Phone)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		//Para responder se usa el paquete json Encode (devolverá en el response(w) un JSON, este JSON será la transformacion del struct en json usando la funcion Encode)
		json.NewEncoder(w).Encode(reqStruct)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update user")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get user")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getall user")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
