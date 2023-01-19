package helper

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	// disini menggunakan tx *sql.Tx karena kita kan butuh deklarasi variabel tx.
	// nah file tx ketika sudah dikenai function Begin, maka akan menjadi object ==> menjadi transactional db
	err := recover()
	// ini untuk menangkap panic(err)
	if err != nil {
		// rollback
		errorRollback := tx.Rollback()
		PanicIfError(errorRollback)
		// nah yang function PanicIfError ini kita panggil dari file error.go
		panic(err)
		// nah panic(err) nya harus setelah Rollback, jadi transaksi digagalkan dulu, lalu dimunculkan panic
	} else {
		// jika tidak ada error, maka lakukan commit
		errorCommit := tx.Commit()
		PanicIfError(errorCommit)
		// ini juga function PanicIfError kita panggil dari file error.go
	}
}