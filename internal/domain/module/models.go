package module

// Module contains metadata of module and list of its Parts (below)
type Module struct {
	Id   int        `json:"id" bson:"_id"`
	Meta ModuleMeta `json:"meta"`
	Data []Part     `json:"parts"`
}

type ModuleMeta struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Logo        string   `json:"logo"`
	Completed   bool     `json:"completed"` // empty field, well be fulfilled for rendering
}

// Part contains all valuable Data. Probably it will be MD or HTML
type Part struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}
