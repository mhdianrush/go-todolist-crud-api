//go:build wireinject
// +build wireinject

package main

import (
	"net/http"
	"project-restful-api/app"
	"project-restful-api/controller"
	"project-restful-api/middleware"
	"project-restful-api/repository"
	"project-restful-api/service"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

func InitializeServer() *http.Server {
	wire.Build(
		app.NewDB,
		validator.New,
		repository.NewCategoryRepository,
		service.NewCategoryService,
		controller.NewCategoryController,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)
	return nil
}
