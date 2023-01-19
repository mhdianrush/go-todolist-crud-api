package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CategoryController interface {
	// disini buat function-function tiap API nya
	// karena ada 5 API, maka buat 5 function. sebenernya ini sudah standar, karena harus mengikuti handler nya http
	// nah disini parameter nya mengikuti si http handler
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// karena menggunakan library http router, maka kita akan menggunakan parameter ketiga yakni params
	// kalau http biasa kan cukup sampai parameter kedua aja
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}