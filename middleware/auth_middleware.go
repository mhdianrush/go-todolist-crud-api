package middleware

import (
	"net/http"
	"project-restful-api/helper"
	"project-restful-api/model/web"
)

// kita akan membuat struct dengan kontrak handler, karena middleware harus dalam bentuk handler
type AuthMiddleware struct {
	// AuthMiddleware nya nanti ditempatkan di yang paling atas, sehingga otomatis dia perlu meneruskan request nya ke handler berikutnya
	Handler http.Handler
}

// selanjutnya kita buat function untuk membuat atribut Handler di atas
func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
	// disini menggunakan pointer karena AuthMiddleware adalah sebuah struct
}

// selanjutnya buat kontrak handler sesuai yang ada di server
func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// menggunakan pointer *AuthMiddleware karena AuthMiddleware adalah sebuah struct
	if "RAHASIA" == request.Header.Get("X-API-Key") {
		// misal untuk API Key nya kita pakai kata "RAHASIA"
		// ini berarti ok, berarti tidak ada terjadi masalah. maka kita tinggal teruskan ke header berikutnya saja
		middleware.Handler.ServeHTTP(writer, request)
		// jadi disini tinggal dilanjutkan saja
	} else {
		// ini berarti error ==> terjadi masalah
		// jadi kalau dia tidak punya X-API-Key di header nya, maka kita akan bilang Unauthorized
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			// disini tidak menggunakan atribut data karena memang tidak ada nilai yang akan ditampung seperti sebelumnya
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}
