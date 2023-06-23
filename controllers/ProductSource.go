package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateProductSourceInput struct {
	Slug          string `json:"slug" binding:"required"`
	ProductSource string `json:"product_source" binding:"required"`
}

func CreateProductSource(c *gin.Context) {
	var input CreateProductSourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productSource := models.ProductSource{Slug: input.Slug, ProductSource: input.ProductSource}
	result := models.DB.First(&productSource, "slug = ?", input.Slug)
	if result.RowsAffected == 0 {
		models.DB.Create(&productSource)
		c.JSON(http.StatusOK, gin.H{"data": productSource})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": "Sudah ada"})
	}
}

func GetAllProductSources(c *gin.Context) {
	var ps []models.ProductSource
	models.DB.Find(&ps)

	c.JSON(http.StatusOK, gin.H{"data": ps})
}
