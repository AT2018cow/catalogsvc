package main

import (
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	stdopentracing "github.com/opentracing/opentracing-go"
)

// Product struct
type Product struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name             string        `json:"name"`
	ShortDescription string        `json:"shortDescription"`
	Description      string        `json:"description"`
	ImageURL1        string        `json:"imageUrl1"`
	ImageURL2        string        `json:"imageUrl2"`
	ImageURL3        string        `json:"imageUrl3"`
	Price            float32       `json:"price"`
	Tags             []string      `json:"tags"`
}

// GetProducts accepts context as input and returns JSON with all the products
func GetProducts(c *gin.Context) {
	var products []Product

	tracer := stdopentracing.GlobalTracer()
	println("printing header from golang ", c.Request.Header)

	for k, v := range c.Request.Header {
		fmt.Println("Header field %q, Value %q\n", k, v)
	}

	productSpanCtx, _ := tracer.Extract(stdopentracing.HTTPHeaders, stdopentracing.HTTPHeadersCarrier(c.Request.Header))

	print("product span context")
	fmt.Println(productSpanCtx)

	productSpan := tracer.StartSpan("db_get_products", stdopentracing.FollowsFrom(productSpanCtx))

	// productSpan, _ := stdopentracing.StartSpanFromContext(c, "db_get_products_new")

	error := collection.Find(nil).All(&products)

	if error != nil {
		message := "Products " + error.Error()
		ext.Error.Set(productSpan, true) // Tag the span as errored
		productSpan.LogEventWithPayload("GET service error", message)
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": message})
		return
	}

	defer productSpan.Finish()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": products})

}

// GetProduct accepts a context as input along with a specific product ID and returns details about that product
// If a product is not found, it returns 404 NOT FOUND
func GetProduct(c *gin.Context) {
	var product Product

	span, _ := stdopentracing.StartSpanFromContext(c, "get_product")
	defer span.Finish()

	productID := c.Param("id")

	span.LogFields(
		log.String("event", "string-format"),
		log.String("ProductID", productID),
	)

	// Check if the Product ID is formatted correctly. If not return an Error - Bad Request
	if bson.IsObjectIdHex(productID) {
		error := collection.FindId(bson.ObjectIdHex(productID)).One(&product)

		if error != nil {
			message := "Product " + error.Error()
			ext.Error.Set(span, true) // Tag the span as errored
			span.LogEventWithPayload("GET product error", message)
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": message})
			return
		}

	} else {
		message := "Incorrect Format for ProductID"
		ext.Error.Set(span, true) // Tag the span as errored
		span.LogEventWithPayload("Incorrect Format for ProductID", message)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": product})

}

// CreateProduct adds a new product item to the database
func CreateProduct(c *gin.Context) {
	var product Product

	error := c.ShouldBindJSON(&product)

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Field Name(s)/ Value(s)"})
		return
	}

	product.ID = bson.NewObjectId()

	error = collection.Insert(&product)

	if error != nil {
		message := "Product " + error.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Product created successfully!", "resourceId": product})

}
