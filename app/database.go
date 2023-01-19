package app

import (
	"database/sql"
	"project-restful-api/helper"
	"time"
)

func NewDB() *sql.DB {
	// *sql.DB disini menggunakan pointer karena DB itu adalah struct di dalam package sql
	db, err := sql.Open("mysql", "root:Ulang.ko.PutusAsa.daa.mang.17@tcp(localhost:3306)/belajar_golang_restful_api")
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
