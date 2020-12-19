package data

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
)

type Set struct{
	Type string `json:"type" validate:"required,eq=bloom-filter"'`
	Config Config `json:"config"`
}

type Config struct{
	Size uint `json:"size" validate:"required"`
	Functions uint `json:"functions" validate:"required,eq=3"`
}

var Filters = map[string]*Filter{} //this can be a redis implementation


func (s *Set) FromJSON (r io.Reader) error{
	e:= json.NewDecoder(r)
	return e.Decode(s)
}

func (s *Set) ToJSON (w io.Writer) error{
	e:= json.NewEncoder(w)
	return e.Encode(s)
}

func (s *Set) Validate() error{
	validate := validator.New()
	return validate.Struct(s)
}