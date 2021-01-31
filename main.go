package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ShippingMethod struct {
	gorm.Model
	Name string
}

type ShippingType struct {
	gorm.Model
	Name string
}

type Order struct {
	gorm.Model
	OrderNr        string
	ShippingMethod uint
	ShippingType   uint
}

type Validation struct {
	OrderNr        string `json:"order_nr"`
	ShippingMethod string `json:"shipping_method"`
	ShippingType   string `json:"shipping_type"`
}

func main() {

	dsn := "root:root@tcp(127.0.0.1:3306)/api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error trying connect to database")
	}

	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {

		p := new(Validation)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		c1 := make(chan uint, 1)
		c2 := make(chan uint, 1)

		go createShippingMethod(db, p.ShippingMethod, c1)
		go createShippingType(db, p.ShippingType, c2)
		var order Order
		order.OrderNr = "Order 1"
		order.ShippingType = <-c2
		order.ShippingMethod = <-c1

		db.Create(&order)
		fmt.Println("------------------------------------------------")
		c.JSON(order)
		return nil
	})

	app.Listen(":3000")
}

func createShippingMethod(db *gorm.DB, name string, c chan uint) {
	fmt.Println("Executing Shipping Method")
	var shippingMethod ShippingMethod
	result := db.FirstOrCreate(&shippingMethod, ShippingMethod{Name: name})
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	c <- shippingMethod.ID
	close(c)
}

func createShippingType(db *gorm.DB, name string, c chan uint) {
	fmt.Println("Executing Shipping Type")
	var shippingType ShippingType
	result := db.FirstOrCreate(&shippingType, ShippingType{Name: name})
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	c <- shippingType.ID
	close(c)
}
