package service

import (
	"context"
	"database/sql"
	"project-restful-api/exception"
	"project-restful-api/helper"
	"project-restful-api/model/domain"
	"project-restful-api/model/web"
	"project-restful-api/repository"

	"github.com/go-playground/validator/v10"
)

type CategoryServiceImpl struct {
	// nah disini kita butuh repository karena nanti manipulasi data nya menggunakan repository
	CategoryRepository repository.CategoryRepository
	/**nah karena CategoryRepository pada package (folder) repository di atas sudah di-set sebagai
	interface, maka kita tidak perlu tambahkan pointer di depan nya
	*/
	// nah kita juga butuh hubungkan dengan database
	DB *sql.DB
	// ini masih dalam bentuk database relational
	/** disini masih pakai database relational karena ketika user meng-input data ke aplikasi, maka data tersebut ya masih dalam bentuk
	relational ==> yakni hanya berisi data per-kategori saja yakni id dan name. oleh karena itu nanti untuk bisa dimasukkan ke database
	dan mengganti data di database, maka kita harus ubah database relational ini menjadi database transactional
	*/
	// nah karena di dalam package sql terdapat struct DB, makanya kita pakai pointer
	Validate *validator.Validate
	// nah ini menggunakan pointer karena di dalam package validator itu terdapat struct yang kita gunakan yakni Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository, DB *sql.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		// dibuat pointer supaya pass by reference. kalau ini diubah, maka yang terhubung dengan nya akan ikut berubah
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

// selanjutnya yakni mengimplementasikan kontrak dari CategoryService nya
func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	// disini struct nya kita validasi dulu yakni yang struct we.CategoryCreateRequest dengan parameter pengganti nya adalah request
	// return value dari function Struct adalah error

	// karena menggunakan database transactional Mysql, maka kita butuh request dalam bentuk transactional (nanti pakai Begin)
	tx, err := service.DB.Begin()
	// menggunakan service.DB.Begin karena di dalam struct CategoryServiceImpl di atas, atribut database nya dalam bentuk relational
	// jadi kita mesti ubah database relational nya menjadi transactional database menggunakan function Begin
	// nah setelah menggunakan function Begin, maka nanti hasil nya dalam bentuk object, yakni return value nya adalah tx dan error
	// tx nya akan menjadi dalam bentuk object ==> menjadi database transactional
	if err != nil {
		panic(err)
	}
	// nah yang tx ini kan akan di-wrap dalam bentuk database transactional
	// jika terjadi error seperti di atas, maka kita akan atasi error nya dengan menjalankan defer function
	defer helper.CommitOrRollback(tx)

	category := domain.Category{
		Name: request.Name,
		// nah request.Name ini diambil dari yang function Create yang atas itu dari struct yang ada pada web.CategiryCreateRequest
		// tujuan nya yakni ketika Name nya sudah ada pada struct CategoryCreateRequest, maka value nya kita akan gunakan di domain.Category
	}
	/** nah karena return value di atas adalah web.CategoryResponse, maka kita harus konversi dulu dari variabel category ini ke
	CategoryResponse
	*/
	category = service.CategoryRepository.Create(ctx, tx, category)
	// syntax ini semacam penghubung antara category service impl dengan category repository dimana di dalam kurung tersebut ada category
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	// disini struct nya kita validasi dulu yakni yang struct we.CategoryCreateRequest dengan parameter pengganti nya adalah request
	// return value dari function Struct adalah error

	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, request.Id)
	// untuk melakukan update, kita harus cek dulu dan pastikan bahwa id nya ada, yakni dengan FindById
	// return value dari FindById pada CategoryRepository ada 2 yakni domain.Category dan error
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	// nah jika id nya tidak ada (tidak ada yang sesuai), maka lakukan panic(err)
	category.Name = request.Name
	// jika id nya ada, maka value dari category.Name nya akan kita kembalikan menjadi value dari request.Name juga

	category = service.CategoryRepository.Update(ctx, tx, category)
	/** jika semua proses sudah terverifikasi, maka selanjutnya ke tahap update. return value dari Update function pada CategoryRepository
	ada 1. perlu diperhatikan juga bahwa nilai parameter ketiga disini diambil dari variabel category di atas yang sudah terverifikasi
	*/
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	// untuk menghapus data tertentu, maka id nya harus di-cek dulu apakah ada atau tidak. jika ada, maka akan diperoleh data Name nya
	// data Name nya tersebut sebenernya ditampung dalam variabel category di atas.
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	// nah kalau ternyata id nya tidak ada, maka kita lakukan panic(err)

	// nah di Delete function ini lebih sederhana
	service.CategoryRepository.Delete(ctx, tx, category)
	// disini kita tidak tampung dalam variabel karena memang return value untuk delete tidak ada
	// karena tidak ada return value nya, maka return nya juga ada di bawah ini.
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	// disini yang categoryId sama dengan yang di struct CategoryRepository
	// return value dari FindById ada 2 yakni domain.Category dan error
	// jika id nya ada, maka value Name nya akan ditampung dalam variabel category
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	// jika id nya tidak ditemukan, maka akan dilakukan panic(err)
	return helper.ToCategoryResponse(category)
	// return nya akan mengembalikan nilai dari value Name yang ada di variabel category
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	categories := service.CategoryRepository.FindAll(ctx, tx)
	// return value dari FindAll hanya ada 1 yakni domain.Category. disini juga tidak menggunakan parameter ketiga
	// karena tujuan kita untuk mendapatkan seluruh data, bukan data berdasarkan id tertentu saja
	// setelah seluruh data hasil query didapat, maka akan ditampung ke variabel categories

	return helper.ToCategoryResponses(categories)
}
