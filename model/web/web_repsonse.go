package web

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	// Code, Status, dan Data disini adalah representasi dari data JSON yang sebelumnya
	// Data menggunakan interface{} karena kan di dalam Data pada JSON tersebut terdapat id dan name
}