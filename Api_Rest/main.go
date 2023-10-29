package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Struct base de parametros a pedir
type registro struct {
	ID     int    `json:"ID"`
	Nombre string `json:"Nombre"`
	Genero string `json:"Genero"`
}

// Nos ayudara a alamacenar peticiones del usuario
type Peticiones []registro

var peticion = Peticiones{{

	ID:     1,
	Nombre: "Jorge",
	Genero: "Masculino",
}}

// Funcion mostrar datos en general
func MostrarTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peticion)
}

func MostrarIndividual(w http.ResponseWriter, r *http.Request) {
	mostrar := mux.Vars(r)
	BusquedaId, err := strconv.Atoi(mostrar["id"])
	//Evaluamos errores
	if err != nil {
		return
	}
	// Busqueda de regsitros
	for _, t := range peticion {
		if t.ID == BusquedaId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
		}

	}
}

// Crear registro de persona
func CrearPeticion(w http.ResponseWriter, r *http.Request) {
	var nuevaPeticion registro
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Los ingresados no son correctos")
	}
	json.Unmarshal(reqBody, &nuevaPeticion)
	nuevaPeticion.ID = len(peticion) + 1
	peticion = append(peticion, nuevaPeticion)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(nuevaPeticion)
}

// Funcion eliminar
func Eliminar(w http.ResponseWriter, r *http.Request) {
	eliminacion := mux.Vars(r)

	BusquedaID, err := strconv.Atoi(eliminacion["id"])

	if err != nil {
		return
	}

	for i, t := range peticion {

		if t.ID == BusquedaID {

			peticion = append(peticion[:i], peticion[i+1:]...)

			fmt.Fprintf(w, "El registro %v se ha eliminado correctamente", BusquedaID)
		}
	}
}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	//recibir las respuestas generales del servidor
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	var Datos registro
	if err != nil {
		fmt.Fprintf(w, "No se ha encontrado ninguna conincidencia")
	}
	recibir, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Los datos que ha ingrasado no son correctos o validos")
	}
	json.Unmarshal(recibir, &Datos)
	for i, b := range peticion {

		if b.ID == requestID {
			// Eliminamos el registro
			peticion = append(peticion[:i], peticion[i+1:]...)
			Datos.ID = b.ID
			peticion = append(peticion, Datos)

			fmt.Fprintf(w, "El registro %v ha sido cambiado con exito", requestID)
		}
	}
}

// Funcion index o pagina principal
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenidos a mis registros")
}

func main() {
	ruta := mux.NewRouter().StrictSlash(true)

	// Ruta index
	ruta.HandleFunc("/", Index)
	ruta.HandleFunc("/todo", MostrarTodo).Methods("GET")
	ruta.HandleFunc("/todo", CrearPeticion).Methods("POST")
	ruta.HandleFunc("/todo/{id}", MostrarIndividual).Methods("GET")
	ruta.HandleFunc("/todo/{id}", Eliminar).Methods("DELETE")
	ruta.HandleFunc("/todo/{id}", Actualizar).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", ruta))
}
