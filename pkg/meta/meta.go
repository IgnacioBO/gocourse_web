package meta

import (
	"os"
	"strconv"
)

// Aqui IRA TODO lo relacionado con el META del request
// meta es info que enviaremos en el body y respuesta
// Info correspondiente al request
// Por ejemplo cantidad de registros, paginas, etc
type Meta struct {
	TotalCount int `json:"total_count"` //tendra total
	Page       int `json:"page"`        //especifica en que pagina estoy
	PerPage    int `json:"per_page"`    //cant de resultado por paginas
	PageCount  int `json:"page_count"`  //cantidad de paginas
}

// New sera para Crar un meta (cono un constucor)
// Le ponesmos solo New y no NewMeta para llamarlo como Meta.New() en vez de Meta.NewMeta()
func New(page, perPage, total int) (*Meta, error) {
	//Si perpage es 0 o negativo se pone un perpage por defecto (en variable de entorno)
	if perPage <= 0 {
		var err error //Se declara para no usar := abajo, ya que "crearia una nueva variable perPage" y da error pq no se ocupa
		perPage, err = strconv.Atoi(os.Getenv("PAGINATOR_LIMIT_DEFAULT"))

		if err != nil {
			return nil, err
		}
	}
	//Tengo que ver la cantidad de paginas que yo tengo en el GetAll, las paginas que vienen segun total y cantidad de resutlado pro apgina
	pageCount := 0
	if total >= 0 {
		// Esta formula es para asegurar que suempre el rsto o sobrente empjue el resultado al PROXIMO numero entero
		//Otra manera seria hacer un un math.Ceil(total/perPage), pero math.Ceu usa float64, por eso n oes muy eficient
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount { //Si la page es mayor que el pageCount, osea estoy en una pagina mayor que la paginacion maxima, page es pageCount. Por ejemplo partir en la page 300 pero hay solo 2
			page = pageCount
		}
	}

	if page < 1 { //Si envia page 0 o menor, devuelve la page 1
		page = 1
	}

	return &Meta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		PageCount:  pageCount,
	}, nil
}

// Metodo offset, dice a partir desde QUE NUMERO de filas comienza trae la info en la page
// Por ejemplo
// page1 1 2 3 offset = 0
// page2 4 5 6 offest = (2-1)*3 = 3 -> del 3, osea parte en 4
// page3 7 8 9 offset = (3-1)*3 = 6 -> del 6, osea parte en 7
func (p *Meta) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Hasta qye numero de fila devuelve
func (p *Meta) Limit() int {
	return p.PerPage
}
