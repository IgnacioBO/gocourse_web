package enrollment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IgnacioBO/gocourse_web/pkg/meta"
)

type (
	//Controller sera una funcion que reciba REspone y Request
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
	}
	//Definiremos una struct para definir el request del Craete, con los campos que quiero recibir y los tags de json
	CreateRequest struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}
	MsgResponse struct {
		ID  string `json:"id"`
		Msg string `json:"msg"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"` //omitempty, asi cuando queremos enviamos la data cuando eta ok y cuando este eror se envie el campo error
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

// Funcion que se encargará de hacer los endopints
// Para eso necesitaremos una struct que se llamara endpoints
// Esta funcion va a DEVOLVER una struct de Endpoints, estos endpoints son los que vamos a poder utuaizlar en unestro dominio (course)
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("create course")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		//Variable con struct de request (datos usaurio)
		var reqStruct CreateRequest
		//r.Body tiene el body del request (se espera JSON) y lo decodifica al struct (reqStruct) (osea pasar el json enviado en el request a un struct)
		err := json.NewDecoder(r.Body).Decode(&reqStruct)
		if err != nil {
			//w.WriteHeader devuelve en el repsonse el CODE que se le indica
			w.WriteHeader(400)
			//Enviaremos la repsuesta con encode y creamos un Sruct ErrorRespone (Creado antes) con un texto
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format. " + err.Error()})
			return
		}

		//Validaciones
		if reqStruct.UserID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "user_id is required"})
			return
		}
		if reqStruct.CourseID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "course_id is required"})
			return
		}

		fmt.Println(reqStruct)
		reqStrucEnJson, _ := json.MarshalIndent(reqStruct, "", " ")
		fmt.Println(string(reqStrucEnJson))

		enrollNuevo, err := s.Create(reqStruct.UserID, reqStruct.CourseID)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: enrollNuevo})
	}
}
