package controller

import (
	"net/http"
	"project-restful-api/helper"
	"project-restful-api/model/web"
	"project-restful-api/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type CategoryControllerImpl struct {
	// disini akan diimplementasikan sesuai dengan kontrak CategoryController
	// yang dibutuhkan di dalam CategoryControllerImpl ini adalah service sebenarnya
	CategoryService service.CategoryService
	// service.CategoryService disini kita ambil dari folder service lalu ke interface CategoryService
	// nah disini tidak perlu pakai pointer, karena CategoryService ini adalah sebuah interface. NB: slice juga ga perlu pointer

}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	// disini return value nya adalah CategoryController dengan tipe data interface
	// tetapi yang kita return kan tetap yang memiliki tipe data struct, yakni CategoryControllerImpl
	return &CategoryControllerImpl{
		// disini menggunakan pointer karena ketika nilai yang disini diubah, maka nilai yang di struct nya juga agar ikut berubah
		CategoryService: categoryService,
	}
}

func (controller *CategoryControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryCreateRequest := web.CategoryCreateRequest{}
	helper.ReadFromRequestBody(request, &categoryCreateRequest)
	// &categoryCreateRequest menggunakan pointer karena isi dari variabel categoryCreateRequest adalah struct

	// selanjutnya kita akan hubungkan dengan struct CategoryService yang ada pada file category_service.go
	categoryResponse := controller.CategoryService.Create(request.Context(), categoryCreateRequest)
	// disini controller menggerakkan CategoryService
	// variabel ini dari return value function Create yang cuma 1 pada CategoryService
	webResponse := web.WebResponse{
		Code: 200,
		// untuk kode sukses
		Status: "OK",
		Data:   categoryResponse,
		// Data: categoryReponse ==> karena value name nya ditampung variabel categoryCreateResponse dan sudah digunakan di categoryResponse
	}
	helper.WriteToResponseBody(writer, webResponse)
	// disini tidak menggunakan pointer karena Data pada struct WebResponse ini hanya akan ditampilkan ke client saja
	// jadi ga ada data yang akan diubah sehingga membutuhkan pass by reference
}
func (controller *CategoryControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryUpdateRequest := web.CategoryUpdateRequest{}
	helper.ReadFromRequestBody(request, &categoryUpdateRequest)
	// berhasil melakukan decode
	/** tetapi key Id pada struct ini kan hanya mengikut pada auto increment yang dilakukan oleh db mysql. sehingga untuk bisa menemukan
	spesifik Id nya, maka kita mesti menggunakan parameters yang sudah dibuat di OpenAPI sebelumnya
	*/
	// kita akan konversi terlebih dahulu yang Id menjadi int kembali
	categoryId := params.ByName("categoryId")
	IdToInt, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)
	categoryUpdateRequest.Id = IdToInt
	// selesai mengubah value Id ke int dalam format JSON nya

	// nah selanjutnya kita akan kasih data JSON nya ke struct CategoryService agar bisa diteruskan ke Categoryrepository
	categoryResponse := controller.CategoryService.Update(request.Context(), categoryUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}
	// selanjutnya kita akan lakukan encoding
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *CategoryControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// pada Delete, tidak ada body, sehingga tidak perlu melakukan decode seperti function di atas
	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)

	controller.CategoryService.Delete(request.Context(), id)
	// ini tidak ditampung dalam variabel karena memang function Delete di-construct untuk tidak mengembalikan nilai (return value)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		// disini tidak menggunakan atribut Data karena tidak ada data yang akan dikembalikan
	}
	// selanjutnya lakukan encoding
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *CategoryControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)

	categoryResponse := controller.CategoryService.FindById(request.Context(), id)
	// ini ditampung dalam variabel karena memang keyword FindById ini memiliki return value yakni hanya satu
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
		// karena akan ada data yang dikembalikan, maka kita harus sertakan atribut Data
	}
	// selanjutnya lakukan encoding
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *CategoryControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// nah karena FindAll ini tidak mencarri data berdasarkan id, maka tidak ada request khusus yang dilakukan user
	// artinya bahwa tidak ada data id khusus yang dikirimkan oleh user sehingga tidak ada yang perlu di-decode
	// nah disini selanjutnya kita bisa langsung kirimkan saja datanya ke struct CategoryService
	categoryResponses := controller.CategoryService.FindAll(request.Context())
	webResponses := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponses,
	}
	// selanjutnya lakukan encoding
	helper.WriteToResponseBody(writer, webResponses)
}
