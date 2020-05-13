package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

var dbconexion *sql.DB
var err error

type Mensaje struct {
	Imprimir string `json:"imprimir"`
}

type Marcacion struct {
	IDMARCACION    string `json:"idmarcacion"`
	FOTOURL        string `json:"foto_url"`
	FECHAMARCACION string `json:"fecha_marcacion"`
	LATITUD        string `json:"latitud"`
	LONGITUD       string `json:"longitud"`
	CELULAR        string `json:"celular"`
}

type Usuario struct {
	IDUSUARIO       string `json:"idusuario"`
	NOMBRE          string `json:"nombre"`
	APELLIDOMATERNO string `json:"apellido_materno"`
	APELLIDOPATERNO string `json:"apellido_paterno"`
	CELULAR         string `json:"celular"`
	FECHAALTA       string `json:"fecha_alta"`
}

func setupRoutes() {
	dbconexion, err = sql.Open("mysql", "root:Asd123**@tcp(localhost:3306)/AsistenciaDB")

	if err != nil {
		fmt.Println("Error al abrir la base de datos")
		panic(err.Error())
	}

	defer dbconexion.Close()

	router := mux.NewRouter()

	router.HandleFunc("/v1/marcacion", SetMarcacion).Methods("POST")
	router.HandleFunc("/v1/marcacion", GetMarcaciones).Methods("GET")
	router.HandleFunc("/v1/marcacion/{idmarcacion}", GetMarcacion).Methods("GET")
	router.HandleFunc("/v1/marcacion/phone/{phone}", GetMarcacionPhone).Methods("GET")
	router.HandleFunc("/v1/marcacion/lastphone/{phone}", GetMarcacionLastPhone).Methods("GET")

	router.HandleFunc("/v1/usuario", SetUsuario).Methods("POST")
	router.HandleFunc("/v1/usuario", GetUsuarios).Methods("GET")
	router.HandleFunc("/v1/usuario/{idusuario}", GetUsuario).Methods("GET")
	router.HandleFunc("/v1/usuario/phone/{phone}", GetUsuarioPhone).Methods("GET")

	http.ListenAndServe(":8000", router)

}

func main() {
	setupRoutes()
}

func SetMarcacion(w http.ResponseWriter, r *http.Request) {

	resultado, err := dbconexion.Prepare("INSERT INTO Marcacion (foto_url,fecha_marcacion,latitud,longitud,celular) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	fotourl := keyVal["foto_url"]
	fechamarcacion := keyVal["fecha_marcacion"]
	latitud := keyVal["latitud"]
	longitud := keyVal["longitud"]
	celular := keyVal["celular"]

	_, err = resultado.Exec(fotourl, fechamarcacion, latitud, longitud, celular)
	if err != nil {
		panic(err.Error())
	}

	var men Mensaje
	men.Imprimir = "Marcacion Creada con exito"
	json.NewEncoder(w).Encode(men)

}

func GetMarcaciones(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var marcacion []Marcacion

	result, err := dbconexion.Query("SELECT idmarcacion,foto_url,fecha_marcacion,latitud,longitud,celular FROM Marcacion")

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var marcacion_temp Marcacion
		err := result.Scan(&marcacion_temp.IDMARCACION, &marcacion_temp.FOTOURL, &marcacion_temp.FECHAMARCACION, &marcacion_temp.LATITUD, &marcacion_temp.LONGITUD, &marcacion_temp.CELULAR)
		if err != nil {
			panic(err.Error())
		}
		marcacion = append(marcacion, marcacion_temp)
	}

	json.NewEncoder(w).Encode(marcacion)
}

func GetMarcacion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var marcacion []Marcacion

	result, err := dbconexion.Query("SELECT idmarcacion,foto_url,fecha_marcacion,latitud,longitud,celular FROM Marcacion WHERE idmarcacion = ? ", params["idmarcacion"])

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var marcacion_temp Marcacion
		err := result.Scan(&marcacion_temp.IDMARCACION, &marcacion_temp.FOTOURL, &marcacion_temp.FECHAMARCACION, &marcacion_temp.LATITUD, &marcacion_temp.LONGITUD, &marcacion_temp.CELULAR)
		if err != nil {
			panic(err.Error())
		}
		marcacion = append(marcacion, marcacion_temp)
	}

	json.NewEncoder(w).Encode(marcacion)
}

func GetMarcacionPhone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var marcacion []Marcacion

	result, err := dbconexion.Query("SELECT idmarcacion,foto_url,fecha_marcacion,latitud,longitud,celular FROM Marcacion WHERE celular = ? ", params["phone"])

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var marcacion_temp Marcacion
		err := result.Scan(&marcacion_temp.IDMARCACION, &marcacion_temp.FOTOURL, &marcacion_temp.FECHAMARCACION, &marcacion_temp.LATITUD, &marcacion_temp.LONGITUD, &marcacion_temp.CELULAR)
		if err != nil {
			panic(err.Error())
		}
		marcacion = append(marcacion, marcacion_temp)
	}

	json.NewEncoder(w).Encode(marcacion)
}

func GetMarcacionLastPhone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var marcacion []Marcacion

	result, err := dbconexion.Query("SELECT idmarcacion,foto_url,fecha_marcacion,latitud,longitud,celular FROM Marcacion WHERE celular = ? order by idmarcacion desc limit 1 ", params["phone"])

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var marcacion_temp Marcacion
		err := result.Scan(&marcacion_temp.IDMARCACION, &marcacion_temp.FOTOURL, &marcacion_temp.FECHAMARCACION, &marcacion_temp.LATITUD, &marcacion_temp.LONGITUD, &marcacion_temp.CELULAR)
		if err != nil {
			panic(err.Error())
		}
		marcacion = append(marcacion, marcacion_temp)
	}

	json.NewEncoder(w).Encode(marcacion)
}

func SetUsuario(w http.ResponseWriter, r *http.Request) {

	resultado, err := dbconexion.Prepare("INSERT INTO Usuario (nombre,apellido_materno,apellido_paterno,celular,fecha_alta) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	nombre := keyVal["nombre"]
	apellidomaterno := keyVal["apellido_materno"]
	apellidopaterno := keyVal["apellido_paterno"]
	celular := keyVal["celular"]
	fecha_alta := keyVal["fecha_alta"]

	_, err = resultado.Exec(nombre, apellidomaterno, apellidopaterno, celular, fecha_alta)
	if err != nil {
		panic(err.Error())
	}

	var men Mensaje
	men.Imprimir = "User was created!"
	json.NewEncoder(w).Encode(men)

}

func GetUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usuario []Usuario

	result, err := dbconexion.Query("SELECT idusuario,nombre,apellido_materno,apellido_paterno,celular,fecha_alta FROM Usuario")

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var usuario_temp Usuario
		err := result.Scan(&usuario_temp.IDUSUARIO, &usuario_temp.NOMBRE, &usuario_temp.APELLIDOMATERNO, &usuario_temp.APELLIDOPATERNO, &usuario_temp.CELULAR, &usuario_temp.FECHAALTA)
		if err != nil {
			panic(err.Error())
		}
		usuario = append(usuario, usuario_temp)
	}

	json.NewEncoder(w).Encode(usuario)
}

func GetUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var usuario []Usuario

	result, err := dbconexion.Query("SELECT idusuario,nombre,apellido_materno,apellido_paterno,celular,fecha_alta FROM Usuario where idusuario = ?", params["idusuario"])

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var usuario_temp Usuario
		err := result.Scan(&usuario_temp.IDUSUARIO, &usuario_temp.NOMBRE, &usuario_temp.APELLIDOMATERNO, &usuario_temp.APELLIDOPATERNO, &usuario_temp.CELULAR, &usuario_temp.FECHAALTA)
		if err != nil {
			panic(err.Error())
		}
		usuario = append(usuario, usuario_temp)
	}

	json.NewEncoder(w).Encode(usuario)
}

func GetUsuarioPhone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var usuario []Usuario

	result, err := dbconexion.Query("SELECT idusuario,nombre,apellido_materno,apellido_paterno,celular,fecha_alta FROM Usuario where celular = ?", params["phone"])

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var usuario_temp Usuario
		err := result.Scan(&usuario_temp.IDUSUARIO, &usuario_temp.NOMBRE, &usuario_temp.APELLIDOMATERNO, &usuario_temp.APELLIDOPATERNO, &usuario_temp.CELULAR, &usuario_temp.FECHAALTA)
		if err != nil {
			panic(err.Error())
		}
		usuario = append(usuario, usuario_temp)
	}

	json.NewEncoder(w).Encode(usuario)
}
