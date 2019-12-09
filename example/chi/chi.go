package chirest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Person struct {
	Age  int
	Name string
}

func startChi() {
	r := SetupRouter()
	fmt.Println("Server listening on port 3333...")
	http.ListenAndServe(":3333", r)
}

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	//r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", Index)
	//r.Get("/hello", Hello)
	r.Get("/param/{key}", ParamTest)
	r.Get("/getjson", ReturnJSON)
	r.Post("/postjson", PostJSON)

	return r
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bare minimum API server in go with chi router")
}

func ParamTest(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	//fmt.Println("Key ", key)
	//render.PlainText(w, r, fmt.Sprintf("key = %s", key))
	render.JSON(w, r, fmt.Sprintf("key = %s", key))
}

func ReturnJSON(w http.ResponseWriter, r *http.Request) {
	p := Person{}
	p.Age = 16
	p.Name = "Mohan"
	render.JSON(w, r, p)
}

func PostJSON(w http.ResponseWriter, r *http.Request) {
	p := Person{}
	defer r.Body.Close()
	err := render.DecodeJSON(r.Body, &p)
	if err != nil {
		http.Error(w, "JSON bind failed", http.StatusBadRequest)
		return
	}
	render.JSON(w, r, "OK")
}
