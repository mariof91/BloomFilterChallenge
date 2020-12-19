package handlers

import (
	"BloomFilter/data"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Sets struct{
	l *log.Logger
}


func NewSets (l *log.Logger) *Sets{
	return &Sets{l}
}

func (s *Sets) AddSet( rw http.ResponseWriter, r *http.Request  ){
	s.l.Println("Handle POST Set")
	vars := mux.Vars(r)
	setName:=vars["set-name"]
	set := r.Context().Value(KeySet{}).(*data.Set)
	data.Filters[setName]=data.New(set.Config.Size)
}


func (s *Sets) PutItem( rw http.ResponseWriter, r *http.Request ){
	s.l.Println("Handle Put item on Set")
	vars := mux.Vars(r)
	setName:=vars["set-name"]
	itemName:=vars["item-name"]

	filter,ok :=data.Filters[setName]
	if !ok{
		s.l.Println("[ERROR] set-name does not exist")
		http.Error(rw, "Error putting Item -- set-name does not exist", http.StatusBadRequest)
		return
	}
	filter.Add([]byte(itemName))


}


type KeySet struct{}

func (s Sets) GetStats(rw http.ResponseWriter, r *http.Request ){
	s.l.Println("Handle GET Stats")
	vars := mux.Vars(r)
	setName:=vars["set-name"]
	filter,ok :=data.Filters[setName]
	if !ok{
		s.l.Println("[ERROR] set-name does not exist")
		http.Error(rw, "Error getting Stat! set-name does not exist", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(rw)
	err:= e.Encode(filter.GetStat())

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (s Sets) GetItem(rw http.ResponseWriter, r *http.Request){
	s.l.Println("Handle GET exists Item")
	vars := mux.Vars(r)
	setName:=vars["set-name"]
	itemName:=vars["item-name"]
	filter,ok :=data.Filters[setName]
	if !ok{
		s.l.Println("[ERROR] set-name does not exist")
		http.Error(rw, "Error getting Item! set-name does not exist", http.StatusBadRequest)
		return
	}
	mayExist:=filter.Test([]byte(itemName))

	rw.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(rw)
	err:= e.Encode(&struct {Exists bool `json:"exist"`}{mayExist})
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (s Sets) MiddlewareValidateSet (next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		set := &data.Set{}

		err := set.FromJSON(r.Body)
		if err != nil {
			s.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading Set", http.StatusBadRequest)
			return
		}

		err = set.Validate()
		if err != nil {
			s.l.Println("[ERROR] validating Set-Name", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating Set-Name: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		// add the set to the context
		ctx := context.WithValue(r.Context(), KeySet{}, set)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
		})
}




