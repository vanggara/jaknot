package main

import (
	"main/controllers"
	"main/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	models.ConnectDatabase() // new!

	router.GET("/product", controllers.CreateUpdateProduct)         // here!
	router.POST("/product-source", controllers.CreateProductSource) // here!
	router.GET("/product-source", controllers.GetAllProductSources) // here!
	router.GET("/excel-dowload-insert", controllers.ExportDataInsert)
	router.GET("/excel-dowload-update", controllers.ExportDataUpdate)

	router.Run("localhost:8080")
}
