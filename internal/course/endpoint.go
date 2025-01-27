package course

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IgnacioBO/gocourse_web/pkg/meta"
	"github.com/gorilla/mux"
)

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
		Name      string `json:"name"`
		StartDate string `json:"start_date"` //Sera de tipo STRING para convertir de string a date en capa servicio
		EndDate   string `json:"end_date"`
	}
	//Definiremos una struct para definir el request del UPDATE, con los campos que quiero y SE PODRAN ACTUALIZAR y los tags de json
	//Seran de tipo puntero * para que puedan venir vacios y poder separar entre vacios "" y que no vengan
	UpdateRequest struct {
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
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
		fmt.Println("delete course")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		variablesPath := mux.Vars(r)
		id := variablesPath["id"]
		fmt.Println("id a eliminar es:", id)
		err := s.Delete(id)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Err: err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		//Aqui le pasamos Response como struct para ahorrar memoria
		//Con puntero (&Response): Encode accede al struct original a través de la dirección de memoria. Esto evita copiar los datos.
		//Sin puntero (Response): Encode recibe una copia del struct completo, lo que puede ocupar más memoria si el struct es muy grande.
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: map[string]string{"id": id, "msg": "success"}})
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
		if reqStruct.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "name is required"})
			return
		}
		if reqStruct.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start_date is required"})
			return
		}
		if reqStruct.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "end_date is required"})
			return
		}

		fmt.Println(reqStruct)
		reqStrucEnJson, _ := json.MarshalIndent(reqStruct, "", " ")
		fmt.Println(string(reqStrucEnJson))

		//Usaremos la s recibida como parametro (de la capa Service y usaremos el metodo CREATE con lo que debe recibir)
		usuarioNuevo, err := s.Create(reqStruct.Name, reqStruct.StartDate, reqStruct.EndDate)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: usuarioNuevo})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update course")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		//Variable con struct de request (datos de atualizacion)
		var reqStruct UpdateRequest
		//r.Body tiene el body del request (se espera JSON) y lo decodifica al struct (reqStruct) (osea pasar el json enviado en el request a un struct)
		err := json.NewDecoder(r.Body).Decode(&reqStruct)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format." + err.Error()})
			return
		}

		if reqStruct.Name != nil && *reqStruct.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "name can't be empty"})
			return
		}
		if reqStruct.StartDate != nil && *reqStruct.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start_date can't be empty"})
			return
		}

		if reqStruct.EndDate != nil && *reqStruct.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "end_date can't be empty"})
			return
		}

		variablesPath := mux.Vars(r)
		id := variablesPath["id"]

		err = s.Update(id, reqStruct.Name, reqStruct.StartDate, reqStruct.EndDate)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 404, Err: err.Error()}) //Aqui devolvemo el posible erro
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: map[string]string{"id": id, "msg": "success"}})

	}
}
func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get course")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		//Aqui usamos MUX para extrar las variables del path (url)
		//Y con ["id"] obtenemos el valor del parametro {id} definido en el main.go (courses/{id})
		//**¿ESTA BIEN ESTO?? USAMOS LIBRERIA EXTERNA EN "internal/courses", no se deberia -> en notasAparte.txt tengo una solucion
		variablesPath := mux.Vars(r)
		id := variablesPath["id"]
		fmt.Println("id es:", id)
		usuario, err := s.Get(id)
		if err != nil {
			if usuario == nil { //Si usuario es vacio da 404
				w.WriteHeader(404)
				json.NewEncoder(w).Encode(&Response{Status: 404, Err: err.Error() + ". course with id " + id + " doesn't exist"}) //Aqui devolvemo el posible erro
				return
			} else {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()}) //Aqui devolvemo el posible erro
				return
			}
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: usuario})

	}
}
func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getall course")
		w.Header().Add("Content-Type", "application/json; charset=utf-8") //Linea miea para que se determine que respondera un json

		//URL.Query() que viene del request
		//Query() devielve un objeto que permite acceder a los parametros d la url (...?campo=123&campo2=hola)
		variablesURL := r.URL.Query()
		//Luego con podemos acceder a los parametos y guardarlos en el struct Filtro (creado en service.go)
		filtros := Filtros{
			Name: variablesURL.Get("name"),
		}
		//Ahora obtendremos el limit y la pagina desde los parametros
		limit, _ := strconv.Atoi(variablesURL.Get("limit"))
		page, _ := strconv.Atoi(variablesURL.Get("page"))

		//Ahora llamaremos al Count del service que creamos (antes de hacer la consulta completa)
		cantidad, err := s.Count(filtros)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()}) //Aqui devolvemo el posible erro
			return
		}
		//Luego crearemos un meta y le agregaremos la cantidad que consultamos, luego el meta lo ageregaremos a la respuesta
		meta, err := meta.New(page, limit, cantidad)

		allCourses, err := s.GetAll(filtros, meta.Offset(), meta.Limit()) //GetAll recibe el offset (desde q resultado mostrar) y el limit (cuantos desde el offset)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()}) //Aqui devolvemo el posible erro
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: allCourses, Meta: meta})
	}
}
