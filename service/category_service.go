package service

import (
	"context"
	"project-restful-api/model/web"
)

type CategoryService interface {
	// ini adalah kontrak nya dalam tipe data interface
	// service ini akan mengikuti jumlah function API nya
	// dalam 1 API maka hanya akan memanggil 1 function di sebuah service
	Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse
	Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse
	Delete(ctx context.Context, categoryId int)
	FindById(ctx context.Context, categoryId int) web.CategoryResponse
	FindAll(ctx context.Context) []web.CategoryResponse
	// sekilas ini mirip seperti yang di category_repository.go ==> hanya saja nanti disini terdapat logic
	// parameter pertama diawali dengan context juga
}
