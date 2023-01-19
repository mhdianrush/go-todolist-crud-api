package exception

import (
	"net/http"
	"project-restful-api/helper"
	"project-restful-api/model/web"

	"github.com/go-playground/validator/v10"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if notFoundError(writer, request, err) {
		return
		// nah ini menggunakan if karena kalau yang function ini dieksekusi, maka function yang setelahnya jangan di-eskekusi
		// nanti return value nya kita gunakan boolean aja di function nya
		// jika ini return nya true, maka function setelah nya jangan di-eksekusi lagi
		// return berfungsi untuk menghentikan eksekusi blok fungsi
	}

	if validationErrors(writer, request, err) {
		return
		// ini hanya untuk request yang aneh-aneh saja seperti ketidak cocokan data dan sebagainya
		// nah ini menggunakan if karena kalau yang function ini dieksekusi, maka function yang setelahnya jangan di-eskekusi
		// return berfungsi untuk menghentikan eksekusi blok fungsi
	}

	// jika terjadi error, kita bisa beri tau jika ada internal server error
	internalServerError(writer, request, err)
}

func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	// ingat bahwa deklarasi 2 variabel seperti ini mirip seperti variabel, err
	// artinya bahwa variabel ok itu sebenernya menggantikan posisi variabel err
	// ValidationsErrors ini ada di dalam package validator
	if ok {
		// jika ok (err), maka akan dilakukan pengembalian ungkapan validasi bahwa terdapat kesalahan ketika input oleh client
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	// notFoundError diawali huruf kecil agar tidak redeclared
	// selanjutnya kita akan lakukan pengecekan
	exception, ok := err.(NotFoundError)
	/** menggunakan err.(NotFoundError) karena err ini diambil dari parameter ketiga di atas (isi pesan error),
	diikuti dengan struct NotFoundError
	*/
	// ingat bahwa deklarasi 2 variabel seperti ini mirip seperti variabel, err
	// artinya bahwa variabel ok itu sebenernya menggantikan posisi variabel err
	if ok {
		// artinya bisa dikonversi ==> artinya jika data dengan id yang dicari adalah err (tidak ada)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
			/** value dari Data disini yakni Error karena di atas itu kita menggunakan err.(NotFoundError) sehingga Error itu sebenernya
			adalah atribut dari struct NotFoundError. atribut Error ini kan akan diisi oleh err.Error() yang merupakan parameter dari
			panic(NewNotFoundError(err.Error())) di file category_service_impl.go tadi
			*/
		}
		helper.WriteToResponseBody(writer, webResponse)
		return true
		// return true ini karena return value nya bool. tetapi ini karena jika ok (artinya jika err ==> tidak ada data dengan id dimaksud)
	} else {
		// artinya tidak bisa dikonversi
		return false
		// return false artinya jika tidak ok (jika tidak err ==> artinya data dengan id yang dicari itu ada)
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	// disini cuma set aja untuk header nya
	writer.WriteHeader(http.StatusInternalServerError)
	// nah disini untuk write status code nya di header
	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
		// nah Data disini hanya berisi err saja karena di parameter ketiga yakni err tersebut tempat menyimpan pesan error nya
		// jadi cukup dibuat Data: err saja untuk mengembalikan data pesan error nya
	}
	helper.WriteToResponseBody(writer, webResponse)
	// nah ini nanti yang akan disampaikan ke client pesan error nya
	// nah jadi harus nya kalau ada exception, maka dapatnya adalah INTERNAL SERVER ERROR
}
