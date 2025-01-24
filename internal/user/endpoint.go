package user

//**Capa endpoint o controlador**

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Struct que tenga todos los endpoints que vayamos a utilizar
// Que teng una fucion que recibe un request y un response
type (
	//Controller sera una funcion que reciba REspone y Request
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create        Controller //Esto es lo mismo que decir Create func(w http.ResponseWriter, r *http.Request), pero como TODOS SON tipo Controller (Definido arriba) nos ahorramos ahcerlo
		Get           Controller
		GetAll        Controller
		Update        Controller
		Delete        Controller
		DeleteClassic Controller
	}
	//Definiremos una struct para definir el request del Craete, con los campos que quiero recibir y los tags de json
	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}
	//Definiremos una struct para definir el request del UPDATE, con los campos que quiero y SE PODRAN ACTUALIZAR y los tags de json
	//Seran de tipo puntero * para que puedan venir vacios y poder separar entre vacios "" y que no vengan
	UpdateRequest struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}

	MsgResponse struct {
		ID  string `json:"id"`
		Msg string `json:"msg"`
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
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		variablesPath := mux.Vars(r)
		id := variablesPath["id"]
		fmt.Println("id a eliminar es:", id)
		err := s.Delete(id)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"id": id, "msg": "success"})
	}
}

// w http.ResponseWriter -> Para enviar RESPUESTA AL CLIENTE (body, headers, statsu code)
// r *http.Request -> Contiende info de la SOLICITUD/REQUEST del cliente (aaceder al metodo http (GET,POST,ETC), paarametros url/query string, body, headers, etc
// *http.Request siempre como PUNTERO (*) por mas eficiencia, para poder modificar datos y es el estandar del package net/http
func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("create user")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

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
		usuarioNuevo, err := s.Create(reqStruct.FirstName, reqStruct.LastName, reqStruct.Email, reqStruct.Phone)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		//Para responder se usa el paquete json Encode (devolverá en el response(w) un JSON, este JSON será la transformacion del struct en json usando la funcion Encode)
		//Antes devolviemoas el reqStruct (que ERA LO MISMO QUE ENVIA EL CLIENTE)
		//Pero ahora devolveremos usuarioNuevo que seria el struct User (del dominio) que tiene como se inserto a la BBDD
		json.NewEncoder(w).Encode(usuarioNuevo)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update user")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		//Variable con struct de request (datos de atualizacion)
		var reqStruct UpdateRequest
		//r.Body tiene el body del request (se espera JSON) y lo decodifica al struct (reqStruct) (osea pasar el json enviado en el request a un struct)
		err := json.NewDecoder(r.Body).Decode(&reqStruct)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"invalid request format"})
			return
		}

		//Validaciones para que sean reqierod
		//Si first name es disinto de nil (osea el puntero NO VIENE VACIO) y le pone "first_name" como vacio (osea el cliene pone first_name = "") da error
		//PERO SI EL CLIENTE NO ENVIA first_namem reqStruct.Firstname sera igual a NIL! entonces no entra
		//OSea se permite NO ENVIAR ESTOS CAMPOS, PERO NO SE PERMITE ENVIARLSO VACIOS
		if reqStruct.FirstName != nil && *reqStruct.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"first_name can't be empty"})
			return
		}

		if reqStruct.LastName != nil && *reqStruct.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"last_name can't be empty"})
			return
		}
		variablesPath := mux.Vars(r)
		id := variablesPath["id"]

		err = s.Update(id, reqStruct.FirstName, reqStruct.LastName, reqStruct.Email, reqStruct.Phone)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()}) //Aqui devolvemo el posible erro
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"id": id, "msg": "success"})

	}
}
func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get user")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		//Aqui usamos MUX para extrar las variables del path (url)
		//Y con ["id"] obtenemos el valor del parametro {id} definido en el main.go (users/{id})
		//**¿ESTA BIEN ESTO?? USAMOS LIBRERIA EXTERNA EN "internal/users", no se deberia -> en notasAparte.txt tengo una solucion
		variablesPath := mux.Vars(r)
		id := variablesPath["id"]
		fmt.Println("id es:", id)
		usuario, err := s.Get(id)
		if err != nil {
			if usuario == nil { //Si usuario es vacio da 404
				w.WriteHeader(404)
				json.NewEncoder(w).Encode(ErrorResponse{err.Error() + ". user with id " + id + " doesn't exist"}) //Aqui devolvemo el posible erro
				return
			} else {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(ErrorResponse{err.Error()}) //Aqui devolvemo el posible erro
				return
			}
		}

		json.NewEncoder(w).Encode(usuario)

	}
}
func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getall user")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		allUsers, err := s.GetAll()
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		json.NewEncoder(w).Encode(allUsers)
	}
}
