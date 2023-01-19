package main

import "testing"

func TestMainTest(t *testing.T) {
	server := InitializeServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}