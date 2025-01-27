package meta

// Aqui IRA TODO lo relacionado con el META del request
// meta es info que enviaremos en el body y respuesta
// Info correspondiente al request
// Por ejemplo cantidad de registros, paginas, etc
type Meta struct {
	TotalCount int `json:"total_count"` //tendra total
}

// New sera para Crar un meta (cono un constucor)
// Le ponesmos solo New y no NewMeta para llamarlo como Meta.New() en vez de Meta.NewMeta()
func New(total int) (*Meta, error) {

	return &Meta{
		TotalCount: total,
	}, nil
}
