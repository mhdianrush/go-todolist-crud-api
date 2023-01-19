package web

type CategoryCreateRequest struct {
	// disini kita gunakan struct, karena kontrak interface nya kan sudah dibuat di file category_service
	// berikut request untuk melakukan Create
	Name string `validate:"required,max=200,min=1" json:"name"`
	// kita ga pakai id karena nanti id nya auto generate oleh database nya (karena auto increment)
}