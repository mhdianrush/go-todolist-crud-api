package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-restful-api/model/domain"
)

type CategoryRepositoryImpl struct {
	// kita akan ikuti kontrak dari category repository
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
	// nah disini karena tidak ada atribut apa pun, maka pakai () dan {} saja
}

func (repository *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	SQL := "insert into category(name) values(?)"
	result, err := tx.ExecContext(ctx, SQL, category.Name)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	category.Id = int(id)
	return category
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	SQL := "update category set name = ? where id = ?"
	_, err := tx.ExecContext(ctx, SQL, category.Name, category.Id)
	// disini menggunakan args nya category.Name dan category.Id karena kalau mau di-update kan harus tau dan jelas di baris dengan id mana
	if err != nil {
		panic(err)
	}
	/** disini kita tidak deklarasikan result nya karena kan kalau kita mengubah data, pasti hanya data tertentu saja dan id nya pun
	sudah jelas. sehingga kalau kita sudah selesai lakukan update, maka kan dari result nya tidak akan dicari LastInsertId nya
	karena ya yang diubah juga kita sudah tau id ke berapa. jadi ga perlu lagi deklarasikan result nya.
	*/
	// selanjutnya kita langsung return kan saja category nya.
	return category
	// category kita return kan karena category disini masih menggantung setelah value id dan name nya diubah.
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category domain.Category) {
	SQL := "delete from category where id = ?"
	_, err := tx.ExecContext(ctx, SQL, category.Id)
	// disini result nya di-ignore aja karena kan ini hanya delete, jadi setelah delete, tidak ada result yang dihasilkan
	if err != nil {
		panic(err)
	}
	// disini tidak pakai return karena memang return value nya tidak ada di ujung function nya
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error) {
	SQL := "select id, name from category where id = ?"
	rows, err := tx.QueryContext(ctx, SQL, categoryId)
	// disini kita gunakan QueryContext karena memang kan kita mau mengambil data dari database
	// jadi proses mengambil data dari database itu dinamakan dengan query
	// balikan nilai dari QueryContext ada 2, yakni rows dan err
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// ini wajib ditutup nantinya di akhir, karena nanti kalau ga ditutup, maka akan terjadi invalid connection
	// invalid connection terjadi karena kita mengirim beberapa perintah tapi salah satu perintah nya belum ditutup
	getData := domain.Category{}
	// variabel category ini kita buat kurung kurawal kosong karena untuk jika tidak ada data nya pada else di bawah, variabel ini dipanggil
	// selanjutnya buat conditional statement jika rows nya di-next, apakah masih ada datanya:
	if rows.Next() {
		// jika ada
		err := rows.Scan(&getData.Id, &getData.Name)
		// return value nya hanya satu yakni error.
		if err != nil {
			panic(err)
		}
		/** menggunakan pointer karena jika data di database nya ada, maka kita langsung hubungkan aja value nya pass by reference ke
		data di database nya menggunakan pointer, jadi bila sewaktu-waktu data di database nya berubah, maka hasil query kita ini
		akan ikut berubah mengikuti data di database juga. nah kalau ga pakai pointer, maka nantihanya akan pass by value, artinya
		data dari database hany di-copy, kalau ada perubahan data di database, maka kita ga akan dapat update perubahan nya.
		*/
		// simple nya, karena domain.Category ini adalah struct Category, maka pakai pointer aja.
		return getData, nil
		/** nah ini sesuai dengan return value yang di atas. artinya jika data by id nya ada maka kan id dan name nya akan dapat,
		sehingga id dan name tersebut akan dimasukkan ke variabel category di atas ini, maka kan otomatis error nya ga ada, makanya nil.
		*/
	} else {
		// jika tidak ada maka kita kasih tau aja error
		return getData, errors.New("category is not found")
		/** ini artinya jika FinndById nya ga ketemu id pencarian nya, maka pasti akan dipanggil variabel category di atas ini,
		dan akan tetap kosong isinya dan sekaligus kita balikkan pesan error.
		*/
	}
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	SQl := "select id, name from category "
	rows, err := tx.QueryContext(ctx, SQl)
	// return value nya ada 2 yakni rows dan error
	// disini pakai functio QueryContext karena memang kan kita mau ambil data dari database
	// disini kita gaperlu pakai parameter ketiga, karena memang di function di atas kita tidak ada sertakan
	// selain itu alasan nya karena kita tidak butuh parameter peng-query karena ya memang bakalan semua data diambil.
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// ini wajib ditutup nantinya di akhir, karena nanti kalau ga ditutup, maka akan terjadi invalid connection
	// invalid connection terjadi karena kita mengirim beberapa perintah tapi salah satu perintah nya belum ditutup

	var categories []domain.Category
	// menggunakan variabel untuk menampung data nya selama ada maka akan dimasukkan ke variabel categories ini
	// tipe data nya sesuai dengan return value nya karena kan return value nya menggunakan slice untuk mengemlompokkan data nya
	for rows.Next() {
		// menggunakan for loop ==> yakni selama data nya ada, maka kita akan balikkan
		// kalau kita pakai if, maka kemungkinan untuk else yakni data nya tidak ada kan masih ada. nah kalau FindAll maka for aja
		getData := domain.Category{}
		err := rows.Scan(&getData.Id, &getData.Name)
		/** nah disini jika data di database masih ada maka akan di-scan, caranya data yang ada di database akan disesuaikan dengan
		tipe data nya. jika int, maka masuk ke Id, jika string maka masuk ke name. karena kan struct Category{} di atas punya Id dan Name
		*/
		if err != nil {
			panic(err)
		}
		categories = append(categories, getData)
		// ini artinya id dan name yang sudah ditampung di variabel getData akak ditambahkan ke variabel categories yang di dalam kurung.
	}
	return categories
	// return categories karena variabel categories masih menggantung dan return berfungsi untuk mengakhiri blok function
}
