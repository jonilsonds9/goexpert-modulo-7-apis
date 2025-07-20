package main

import (
	"github.com/jonilsonds9/goexpert-modulo-7-apis/configs"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/entity"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/infra/database"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}
