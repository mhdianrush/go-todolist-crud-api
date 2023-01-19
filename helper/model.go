package helper

import (
	"project-restful-api/model/domain"
	"project-restful-api/model/web"
)

func ToCategoryResponse(category domain.Category) web.CategoryResponse {
	/** nah di dalam kurung itu kita menggunakan seperti itu karena value yang sudah diterima oleh struct
	Category dari struct CategoryCreateRequest harus kita passing juga ke return value nya yakni
	web.CategoryResponse
	*/
	return web.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
		/** nah walaupun data yang diperoleh oleh struct Category hanya atribut Name, namun kita tetap
		gunakan pengambilan Id disini, ya kalau nanti Id nya ga ada ya pasti 0 saja default value nya
		*/
	}
	/** ingat disini kita langsung menggunakan return karena ketika datanya sudah didapat, kita sekaligus
	menghentikan blok function ini
	*/
}

func ToCategoryResponses(categories []domain.Category) []web.CategoryResponse {
	var categoryResponses []web.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, ToCategoryResponse(category))
	}
	return categoryResponses
}
