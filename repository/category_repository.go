package repository

import (
	"context"
	"database/sql"
	"project-restful-api/model/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category
	// ini adalah function untuk melakukan create (menambah data)
	// menggunakan parameter ketiga untuk menambahkan data name sesuai yang ada di struct Category
	Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category
	// ini adalah function untuk mengubah/mengupdate data
	// menggunakan parameter ketiga untuk mengubah data name lama menjadi data name baru berdasarkan id
	Delete(ctx context.Context, tx *sql.Tx, category domain.Category)
	// ini adalah function untuk menghapus data
	// menggunakan parameter ketiga untuk menghapus data name berdasarkan id
	// tidak menggunakan return value karena ini kan menghapus data saja, ga ada pengembalian data
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error)
	// ini adalah function untuk menemukan data yang dicari berdasarkan id
	/**
	disini menggunakan parameter ketiga hanya berupa id dengan tipe data int saja karena kita hanya mencari
	data menggunakan id. nanti setelah itu return value nya akan mengambalikan data id dan name
	*/
	// menggunakan return value karena akan ada pengembalian data setelah dicari data nya
	// return value nya kita buat 2 aja, jika ada maka dibalikkan nilainya, kalau ga ada kita kasih respon error aja
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Category
	// ini function untuk memperlihatkan seluruh data
	/** tidak menggunakan parameter ketiga karena tidak seperti function sebelumnya yang melakukan perubahan
	terhadap baris data yang membutuhkan id dan name maupun id saja, disini setelah database transactional,
	maka seluruh data akan diberikan, oleh karena itu di return value nya ditampung dalam bentuk slice.
	*/
}
