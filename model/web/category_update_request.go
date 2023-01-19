package web

type CategoryUpdateRequest struct {
	Id   int    `validate:"required" json:"id"`
	Name string `validate:"required,max=200,min=1" json:"name"`
	/** nah di dalam validator tersebut harus ditambahkan tanda ``. jika wajib, maka tambahkan instruksi required. di dalam `` tidak boleh
	ada spasi. yang max=200 itu adalah maksimal huruf yang kita buat pada database sebelumnya yakni varchar(200) dan min=1 itu artinya
	huruf nya tidak boleh tidak ada karena perintah nya adalah not null.
	sementara yang id itu kan sudah auto generate dari database nya menggunakan perintah auto increment, cukup required aja
	*/
}
