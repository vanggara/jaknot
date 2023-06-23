package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type Product struct {
	Name       string `json:"name" binding:"required"`
	Weight     int32  `json:"weight" binding:"required"`
	Detail     string `json:"detail" binding:"required"`
	UrlVideo   string `json:"url_video"`
	Brand      string `json:"brand" binding:"required"`
	BoxItem    string `json:"box_item" binding:"required"`
	SKU        string `json:"sku" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Image1     string `json:"image_1" binding:"required"`
	Image2     string `json:"image_2"`
	Image3     string `json:"image_3"`
	Image4     string `json:"image_4"`
	Image5     string `json:"image_5"`
	Prices     int64  `json:"prices" binding:"required"`
	Stock      int16  `json:"stock" binding:"required"`
	UrlProduct string `json:"url_product" gorm:"not null"`
}

type RequestJaknot struct {
	OperationName string `json:"operationName"`
	Variables     struct {
		Slug string `json:"slug"`
	} `json:"variables"`
	Query string `json:"query"`
}

func CreateUpdateProduct(c *gin.Context) {
	product_source := GetAllProductSource()
	var result = ""
	for i := 0; i < len(product_source); i++ {
		slugParam := product_source[i].Slug
		response, err := GetExternalAPI(slugParam)
		if err != nil {
			panic(err)
		}

		var obj map[string]interface{}
		err2 := json.Unmarshal([]byte(response), &obj)
		if err2 != nil {
			panic(err2)
		}
		data := obj["data"].(map[string]interface{})
		if data["productBySlug"] != nil {
			productBySlug := data["productBySlug"].(map[string]interface{})
			name := productBySlug["name"].(string)
			weight := productBySlug["weight"].(map[string]interface{})["value"].(float64) * 1000
			detail := productBySlug["detail"].(map[string]interface{})["content"].(string)
			brand := productBySlug["detail"].(map[string]interface{})["brand"].(map[string]interface{})["name"].(string)
			boxItem := productBySlug["detail"].(map[string]interface{})["boxItem"].(string)
			skus := productBySlug["skus"].([]interface{})[0]
			sku := skus.(map[string]interface{})["id"].(string)
			slug := skus.(map[string]interface{})["slug"].(string)
			images := skus.(map[string]interface{})["images"]
			images1 := images.([]interface{})[0].(map[string]interface{})["large"].(string)
			images2 := images.([]interface{})[1].(map[string]interface{})["large"].(string)
			images3 := images.([]interface{})[2].(map[string]interface{})["large"].(string)
			images4 := images.([]interface{})[3].(map[string]interface{})["large"].(string)
			images5 := images.([]interface{})[4].(map[string]interface{})["large"].(string)
			price := skus.(map[string]interface{})["prices"].([]interface{})[0].(map[string]interface{})["top"].(float64) + 13000
			stock := skus.(map[string]interface{})["stocks"].([]interface{})[9].(map[string]interface{})["stockRemaining"]
			location := skus.(map[string]interface{})["stocks"].([]interface{})[9].(map[string]interface{})["name"].(string)
			urlVideo := os.DevNull
			if stock == nil {
				stock = 0.0
			}
			isAvailable := stock.(float64) > 0

			// fmt.Println(sku)
			product := models.Product{Name: name, Weight: int32(weight), Detail: detail, UrlVideo: urlVideo,
				Brand: brand, BoxItem: boxItem, SKU: sku, Slug: slug, Image1: images1, Image2: images2,
				Image3: images3, Image4: images4, Image5: images5, Location: location, Prices: price,
				Stock: int64(stock.(float64)), IsAvailable: isAvailable}
			result := models.DB.First(&product, "sku = ?", &sku)
			fmt.Println(result)
			if result.RowsAffected > 0 {
				product.Prices = price
				product.Stock = int64(stock.(float64))
				product.IsAvailable = isAvailable
				models.DB.Save(&product)
			} else {
				models.DB.Create(&product)
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}

func GetExternalAPI(slug string) ([]byte, error) {
	variables := map[string]string{"slug": slug}
	json_variables, _ := json.Marshal(variables)

	values := map[string]string{"operationName": "ProductDetail",
		"variables": string(json_variables),
		"query":     "query ProductDetail($slug: String!) {productBySlug(slug: $slug) {id name categoryIds weight {unit value __typename} detail { hasWarranty warranty warrantyInfo content specification { name value __typename } info guideUrl videoUrls displayVideoUrls brand { name brandPageUrl __typename } boxItem __typename}skus { id slug isWatched isActive isComingSoon color { border code colorCode colorCode2 name __typename } images { large __typename } prices { top bottom memberType name promoEndAt promoName promoEndTimeRemaining __typename } stocks { branchId isReminded lastStockAt name stockAvailableAt stockRemaining isStockAvailableSoon isMainWarehouse __typename } volume { height length unit width __typename } __typename}variantValues { label value __typename}variants { id labels options { label value options { label value skus { color { border code colorCode colorCode2 name __typename } id image stocks { branchId isStockAvailable name __typename } price { top bottom memberType name promoEndAt promoName promoEndTimeRemaining __typename } slug __typename } __typename } skus { color { border code colorCode colorCode2 name __typename } id image price { top bottom memberType name promoEndAt promoName promoEndTimeRemaining __typename } stocks { branchId isStockAvailable name __typename } slug __typename } __typename } __typename}__typename}__typename}"}
	json_data, err := json.Marshal(values)

	response, error := http.Post("https://graphql.jakartanotebook.com/graphql/v3", "application/json", bytes.NewBuffer(json_data))

	if error != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, errRes := ioutil.ReadAll(response.Body)
	if errRes != nil {
		log.Fatal(errRes)
	}
	return responseData, errRes
}

func CountAllProductSource() int64 {
	var ps []models.ProductSource
	results := models.DB.Find(&ps)
	count := results.RowsAffected
	return count
}

func GetAllProductSource() []*models.ProductSource {
	ps := []*models.ProductSource{}
	models.DB.Find(&ps)
	return ps
}

func ExportDataUpdate(c *gin.Context) {
	f := excelize.NewFile()
	// Create a new sheet.
	// Set value of a cell.
	product_source := GetAllProduct()
	f.SetCellStr("Sheet1", "B1", "Product Name")
	f.SetCellStr("Sheet1", "C1", "")
	f.SetCellStr("Sheet1", "D1", "")
	f.SetCellStr("Sheet1", "E1", "Stock")
	f.SetCellStr("Sheet1", "F1", "")
	f.SetCellStr("Sheet1", "G1", "Prices")
	f.SetCellStr("Sheet1", "H1", "Stock")
	f.SetCellStr("Sheet1", "I1", "")
	f.SetCellStr("Sheet1", "J1", "SKU")
	f.SetCellStr("Sheet1", "K1", "Active")
	f.SetCellStr("Sheet1", "L1", "Weight")
	for i := 2; i <= len(product_source); i++ {
		f.SetCellValue("Sheet1", "B"+fmt.Sprint(i), product_source[i-1].Name)
		f.SetCellValue("Sheet1", "C"+fmt.Sprint(i), nil)
		f.SetCellValue("Sheet1", "D"+fmt.Sprint(i), 1)
		f.SetCellValue("Sheet1", "E"+fmt.Sprint(i), product_source[i-1].Stock)
		f.SetCellValue("Sheet1", "F"+fmt.Sprint(i), nil)
		f.SetCellValue("Sheet1", "G"+fmt.Sprint(i), product_source[i-1].Prices)
		f.SetCellValue("Sheet1", "H"+fmt.Sprint(i), product_source[i-1].Stock)
		f.SetCellValue("Sheet1", "I"+fmt.Sprint(i), nil)
		f.SetCellValue("Sheet1", "J"+fmt.Sprint(i), product_source[i-1].SKU)
		f.SetCellValue("Sheet1", "K"+fmt.Sprint(i), "Aktif")
		f.SetCellValue("Sheet1", "L"+fmt.Sprint(i), product_source[i-1].Weight)
	}
	// Set active sheet of the workbook.
	// Save xlsx file by the given path.
	if err := f.SaveAs("Jaknot_" + time.Now().Format("01-02-2006") + ".xlsx"); err != nil {
		println(err.Error())
	}
}

func GetAllProduct() []*models.Product {
	ps := []*models.Product{}
	models.DB.Find(&ps)
	return ps
}

func ExportDataInsert(c *gin.Context) {
	f := excelize.NewFile()
	// Create a new sheet.
	// Set value of a cell.
	product_source := GetAllProduct()
	f.SetCellStr("Sheet1", "B1", "Product Name")
	f.SetCellStr("Sheet1", "C1", "Description")
	f.SetCellStr("Sheet1", "D1", "Category Code")
	f.SetCellStr("Sheet1", "E1", "Weight")
	f.SetCellStr("Sheet1", "F1", "Minimum Order")
	f.SetCellStr("Sheet1", "G1", "Etalase Number")
	f.SetCellStr("Sheet1", "H1", "Preorder Estimate")
	f.SetCellStr("Sheet1", "I1", "Condition")
	f.SetCellStr("Sheet1", "J1", "Product Foto 1")
	f.SetCellStr("Sheet1", "K1", "Product Foto 2")
	f.SetCellStr("Sheet1", "L1", "Product Foto 3")
	f.SetCellStr("Sheet1", "M1", "Product Foto 4")
	f.SetCellStr("Sheet1", "N1", "Product Foto 5")
	f.SetCellStr("Sheet1", "R1", "SKU")
	f.SetCellStr("Sheet1", "S1", "Status")
	f.SetCellStr("Sheet1", "T1", "Total Stock")
	f.SetCellStr("Sheet1", "U1", "Price")
	f.SetCellStr("Sheet1", "V1", "Courier")
	f.SetCellStr("Sheet1", "W1", "Insurance")
	for i := 2; i <= len(product_source); i++ {
		f.SetCellValue("Sheet1", "B"+fmt.Sprint(i), product_source[i-1].Name)
		f.SetCellValue("Sheet1", "C"+fmt.Sprint(i), nil)
		f.SetCellValue("Sheet1", "D"+fmt.Sprint(i), 1)
		f.SetCellValue("Sheet1", "E"+fmt.Sprint(i), product_source[i-1].Stock)
		f.SetCellValue("Sheet1", "F"+fmt.Sprint(i), nil)
		f.SetCellValue("Sheet1", "G"+fmt.Sprint(i), product_source[i-1].Prices)
		f.SetCellValue("Sheet1", "H"+fmt.Sprint(i), product_source[i-1].Stock)
		f.SetCellValue("Sheet1", "I"+fmt.Sprint(i), nil)
		f.SetCellValue("Sheet1", "J"+fmt.Sprint(i), product_source[i-1].SKU)
		f.SetCellValue("Sheet1", "K"+fmt.Sprint(i), "Aktif")
		f.SetCellValue("Sheet1", "L"+fmt.Sprint(i), product_source[i-1].Weight)
	}
	// Set active sheet of the workbook.
	// Save xlsx file by the given path.
	if err := f.SaveAs("Jaknot_" + time.Now().Format("01-02-2006") + ".xlsx"); err != nil {
		println(err.Error())
	}
}
