package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=vanggar4 dbname=jaknot port=5432 timezone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // change the database provider if necessary

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Post{})          // register Post model
	database.AutoMigrate(&Product{})       // register Product model
	database.AutoMigrate(&ProductSource{}) // register Product Source model

	DB = database
}
