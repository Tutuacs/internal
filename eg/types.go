package eg

type Eg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateEgDto struct {
	Name string `json:"name" validate:"required,min=1"` // To dont show on json: `json:"-"`
}
