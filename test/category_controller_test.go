package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"project-restful-api/app"
	"project-restful-api/controller"
	"project-restful-api/helper"
	"project-restful-api/middleware"
	"project-restful-api/model/domain"
	"project-restful-api/repository"
	"project-restful-api/service"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	// *sql.DB disini menggunakan pointer karena DB itu adalah struct di dalam package sql
	db, err := sql.Open("mysql", "root:Ulang.ko.PutusAsa.daa.mang.17@tcp(localhost:3306)/belajar_golang_restful_api_test")
	helper.PanicIfError(err)
	// "mysql" disini adalah driver nya
	// 3306 ini adalah default nya

	// jika mau set connection pooling nya, silahkan:
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	// ini 60 menit
	db.SetConnMaxIdleTime(10 * time.Minute)
	// jadi kalau misal idle nya 10 menit, maka lebih baik kita close aja
	return db
	// menggunakan return db karena variabel db disini masih menggantung sehingga dengan return, maka kita menghentikan eksekusi nya
}

// ini akan menjadi lebih ke integration test (langsung nge-hit endpoint nya)
func setupRouter(db *sql.DB) http.Handler {

	validate := validator.New()
	// nah validate disini untuk membuat validasi
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoryController)
	// router nya sudah dipindahkan ke file router.go di folder app. disini kita panggil aja router nya dengan ditampung dalam var router
	return middleware.NewAuthMiddleware(router)
	// disini sekalian aja kita masukkan middleware nya
	// disini return middleware karena return value nya di atas adalah http.Handler
	// middleware kan dibuat menggunakan aturan dari Handler
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
	// fitur TRUNCATE ini adalah fitur pada database mysql untuk menghapus seleuruh data yang ada pada tabel
	/** menggunakan keyword Exec karena isi di dalam kurung nya membutuhkan string sebagai parameter pertamanya, sehingga bisa diisi
	statement dengan fitur TRUNCATE yang ada pada db mysql
	*/
}

// selanjutnya yakni skenario test
func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	// selanjutnya buat request body nya
	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	// yang di dalam kurung tersebut adalah name yang akan di-test. kita buat dalam bentuk JSON ==> karena OpenAPI kita menggunakan JSON

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// selanjutnya buat request Header nya
	request.Header.Add("Content-Type", "application/json")
	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 200, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
	// ini karena di dalam key data itu kan terdapat 2 nilai yakni id dan name, makanya dibuat begini biar spesifik
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	// selanjutnya buat request body nya
	requestBody := strings.NewReader(`{"name": ""}`)
	// yang di dalam kurung tersebut adalah name yang akan di-test. kita buat dalam bentuk JSON ==> karena OpenAPI kita menggunakan JSON
	// karena ini test untuk yang failed, maka kita buat aja misal jika menggunakan string kosong

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// selanjutnya buat request Header nya
	request.Header.Add("Content-Type", "application/json")
	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 400, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya
	// status code nya 400 karena ini kan function untuk yang gagal (bad request)

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
	// tidak menggunakan yang responseBody["data"] karena kan tidak ada nama yang dimasukkan (hanya string kosong)
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	tx, _ := db.Begin()
	// ini membuat menjadi object, yakni dari database relational ke database transactional agar bisa transfer data
	categoryRepository := repository.NewCategoryRepository()
	// pakai repository karena kan repository ini di-construct untuk bertransaksi dengan database yakni terdapat tx sebagai db transactional
	category := categoryRepository.Create(context.Background(), tx, domain.Category{
		// nah disini kita Create dulu, karena kan data di database sudah di-truncate, jadi perlu kita tambahkan terlebih dahulu
		// lalu di-update
		Name: "Gadget",
		// nah untuk yang id nya nanti auto increment dari database nya langsung
	})
	tx.Commit()
	// kalau sudah selesai maka harus di-commit

	router := setupRouter(db)

	// selanjutnya buat request body nya
	requestBody := strings.NewReader(`{"name": "Laptop"}`)
	// yang di dalam kurung tersebut adalah name yang akan di-test. kita buat dalam bentuk JSON ==> karena OpenAPI kita menggunakan JSON

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// selanjutnya buat request Header nya
	request.Header.Add("Content-Type", "application/json")
	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 200, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Laptop", responseBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	tx, _ := db.Begin()
	// ini membuat menjadi object, yakni dari database relational ke database transactional agar bisa transfer data
	categoryRepository := repository.NewCategoryRepository()
	// pakai repository karena kan repository ini di-construct untuk bertransaksi dengan database yakni terdapat tx sebagai db transactional
	category := categoryRepository.Create(context.Background(), tx, domain.Category{
		// nah disini kita Create dulu, karena kan data di database sudah di-truncate, jadi perlu kita tambahkan terlebih dahulu
		// lalu di-update
		Name: "Gadget",
		// nah untuk yang id nya nanti auto increment dari database nya langsung
	})
	tx.Commit()
	// kalau sudah selesai maka harus di-commit

	router := setupRouter(db)

	// selanjutnya buat request body nya
	requestBody := strings.NewReader(`{"name": ""}`)
	// yang di dalam kurung tersebut adalah name yang akan di-test. kita buat dalam bentuk JSON ==> karena OpenAPI kita menggunakan JSON

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// selanjutnya buat request Header nya
	request.Header.Add("Content-Type", "application/json")
	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 400, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 400 di tengah harapan nya bad request. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	tx, _ := db.Begin()
	// ini membuat menjadi object, yakni dari database relational ke database transactional agar bisa transfer data
	categoryRepository := repository.NewCategoryRepository()
	// pakai repository karena kan repository ini di-construct untuk bertransaksi dengan database yakni terdapat tx sebagai db transactional
	category := categoryRepository.Create(context.Background(), tx, domain.Category{
		// nah disini kita Create dulu, karena kan data di database sudah di-truncate, jadi perlu kita tambahkan terlebih dahulu
		// lalu di-update
		Name: "Gadget",
		// nah untuk yang id nya nanti auto increment dari database nya langsung
	})
	tx.Commit()
	// kalau sudah selesai maka harus di-commit

	router := setupRouter(db)

	// disini tidak perlu ada requestBody, karena kita ga mengirimkan apapun, hanya GET sesuatu saja

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// tidak perlu membuat request Header Content-Type nya application/json, karena request body nya tidak ada

	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 200, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	// tidak perlu ada data yang dimasukkan ke database, karena kan ini yang Failed nya
	// jadi seolah-olah akan nge-GET data yang memang tidak ada

	router := setupRouter(db)

	// disini tidak perlu ada requestBody, karena kita ga mengirimkan apapun, hanya GET sesuatu saja

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// tidak perlu membuat request Header Content-Type nya application/json, karena request body nya tidak ada

	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 404, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	tx, _ := db.Begin()
	// ini membuat menjadi object, yakni dari database relational ke database transactional agar bisa transfer data
	categoryRepository := repository.NewCategoryRepository()
	// pakai repository karena kan repository ini di-construct untuk bertransaksi dengan database yakni terdapat tx sebagai db transactional
	category := categoryRepository.Create(context.Background(), tx, domain.Category{
		// nah disini kita Create dulu, karena kan data di database sudah di-truncate, jadi perlu kita tambahkan terlebih dahulu
		// lalu di-update
		Name: "Gadget",
		// nah untuk yang id nya nanti auto increment dari database nya langsung
	})
	tx.Commit()
	// kalau sudah selesai maka harus di-commit

	// ini perlu ada data yang dimasukkan ke database karena ingin dihapus datanya sebagai percobaan

	router := setupRouter(db)

	// disini tidak ada request body, karena tidak ada data yang dikirimkan, hanya menghapus data saja

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// selanjutnya buat request Header nya
	request.Header.Add("Content-Type", "application/json")
	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 200, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	// tidak perlu ada data yang dimasukkan ke database

	router := setupRouter(db)

	// disini tidak ada request body, karena tidak ada data yang dikirimkan, hanya menghapus data saja

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// selanjutnya buat request Header nya
	request.Header.Add("Content-Type", "application/json")
	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 404, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestListCategoriesSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	tx, _ := db.Begin()
	// ini membuat menjadi object, yakni dari database relational ke database transactional agar bisa transfer data
	categoryRepository := repository.NewCategoryRepository()
	// pakai repository karena kan repository ini di-construct untuk bertransaksi dengan database yakni terdapat tx sebagai db transactional
	category1 := categoryRepository.Create(context.Background(), tx, domain.Category{
		// nah disini kita Create dulu, karena kan data di database sudah di-truncate, jadi perlu kita tambahkan terlebih dahulu
		// lalu di-update
		Name: "Gadget",
		// nah untuk yang id nya nanti auto increment dari database nya langsung
	})
	category2 := categoryRepository.Create(context.Background(), tx, domain.Category{
		// nah disini kita Create dulu, karena kan data di database sudah di-truncate, jadi perlu kita tambahkan terlebih dahulu
		// lalu di-update
		Name: "Computer",
		// nah untuk yang id nya nanti auto increment dari database nya langsung
	})
	tx.Commit()
	// kalau sudah selesai maka harus di-commit

	router := setupRouter(db)

	// disini tidak perlu ada requestBody, karena kita ga mengirimkan apapun, hanya GET sesuatu saja

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// tidak perlu membuat request Header Content-Type nya application/json, karena request body nya tidak ada

	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "RAHASIA")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 200, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	// nah kita tampung dalam slice saja
	var categories = responseBody["data"].([]interface{})
	// disini data nya dikelompok kan dulu dalam slice dengan tipe data interface{}
	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})
	// disini dipisah satu per satu data yang tadi ada di dalam slice

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])

	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	// nah setelah function truncateCategory() di jalankan, maka kita harus create data ke database dulu

	router := setupRouter(db)

	// disini tidak perlu ada requestBody, karena kita ga mengirimkan apapun, hanya GET sesuatu saja

	// selanjutnya buat http request nya
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	// untuk path id nya, kita pakai strconv aja dari int ke string menggunakan Itoa. nanti kan yang category punya Id tuh ==> category.Id
	// ini seperti yang di pembelajaran golang web
	// parameter pertama yakni http method nya. parameter kedua yakni URL nya. parameter ketiga yakni body nya, karena ada body, masukin.
	// kalau tidak ada body, maka bisa diisi nil saja di parameter ketiga

	// tidak perlu membuat request Header Content-Type nya application/json, karena request body nya tidak ada

	// request Header nya juga mesti ada yang X-API-Key
	request.Header.Add("X-API-Key", "SALAH")

	// selanjutnya buat recorder nya
	recorder := httptest.NewRecorder()

	// selanjutnya tinggal panggil
	router.ServeHTTP(recorder, request)

	// untuk mendapatkan response nya
	response := recorder.Result()
	// jika response nya mau dibaca body nya bisa. ==> response.Body()

	// disini akan di-cek
	assert.Equal(t, 401, response.StatusCode)
	// t di awal itu emang udah dari sana nya gitu. 200 di tengah harapan nya OK. response.StatusCode ini untuk mengambil status code nya

	// nah misal jika ingin membaca data yang ada di response body
	body, _ := io.ReadAll(response.Body)
	// variabel body memiliki kembalian []byte yang artinya ini juga bisa dianggap JSON
	var responseBody map[string]interface{}
	// variabel responseBody ini akan dijadikan variabel untuk menampung nilai konversi dari JSON ke tipe data di golang
	/** kita menggunakan map[string]interface{} disini yakni map[string] nya nanti sebagai key, lalu interface{} bisa menampung berbagai
	macam tipe data nantinya
	*/
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}
	// function Unmarshal digunakan untuk melakukan konversi data JSON ke tipe data golang
	// nanti yang di dalam kurung itu yang body, nilainya akan di-unmarshal dari JSON ke golang lalu di masukkan ke paramater &responseBody
	/** menggunakan pointer &responseBody karena supaya pass by reference. jika yang di dalam kurung sudah menjadi tipe data golang, maka
	yang variabel di atas nya itu yakni responseBody kan nanti di-print, maka biar jadi ke tipe data golang juga
	*/
	fmt.Println(responseBody)

	// nah jika mau mengecek secara lebih detail lagi, maka bisa lakukan seperti berikut:
	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}
