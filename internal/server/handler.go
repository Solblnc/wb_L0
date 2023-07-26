package server

import (
	"L0/internal/nats"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	NatsServer nats.Nats
}

func NewHandler(NatsServer nats.Nats) *Handler {
	return &Handler{NatsServer: NatsServer}
}

func (h *Handler) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/order/{id}", h.GetOrderByID)

	return r
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect input format"))
	}

	idInt, err := strconv.Atoi(id)
	fmt.Println(idInt)
	if err != nil {
		w.Write([]byte("\ncannot convert id to integer"))
	}

	order, err := h.NatsServer.GetOrderFromCache(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\nno order specified with that id\n"))
	}

	if err = json.NewEncoder(w).Encode(order); err != nil {
		panic(err)
	}

}
