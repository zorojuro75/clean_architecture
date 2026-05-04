package main 
import (
    "errors"
    "fmt"
)

type Product struct {
	ID			uint
	Name		string
	Price		float64
	Stock		int
	Discount	float64
}

func (p Product) Validate() error {
	if p.Name == ""{
		return errors.New("product name cannot be empty")
	}
	if p.Price <= 0 {
        return errors.New("product price must be greater than zero")
    }
    if p.Stock < 0 {
        return errors.New("product stock cannot be negative")
    }
	if p.Discount < 0 || p.Discount >= 1 {
		return errors.New("Discount must be greater than or equal 0 and less than 1")
	}
	return nil
}

func (p *Product) applyDiscount() (float64, error){
	if p.Discount < 0 || p.Discount > 1 {
		return 0, errors.New("Invalid discount percentage")
	}
	discounted := p.Price * p.Discount
	return discounted, nil
}

func (p *Product)orderTotal(discountedPrice float64, quantity int) (float64, error) {
	err := p.updateStock(quantity)
	if err != nil{
		return 0, err
	}
	return discountedPrice * float64(quantity), nil
}

func (p *Product) updateStock (quantity int) error{
	if p.Stock<=0 {
		return errors.New("Product unavailable")
	}
	p.Stock = p.Stock - quantity
	return nil
}

func (p *Product) printReceipt(quantity int) error{
	dicountedAmount, err := p.applyDiscount()
	if err !=nil{
		fmt.Println("Error: ", err)
		return err
	}
	discountedPrice := p.Price - dicountedAmount
	total, err := p.orderTotal(discountedPrice, quantity)
	if err != nil{

	}
	
	fmt.Println("====== Receipt ======")
    fmt.Printf("Product  : %s\n", p.Name)
    fmt.Printf("Price    : $%.2f\n", p.Price)
    fmt.Printf("Discount : %.0f%%\n", p.Discount*100)
    fmt.Printf("New Price: $%.2f\n", discountedPrice)
    fmt.Printf("Quantity : %d\n", quantity)
    fmt.Printf("Total    : $%.2f\n", total)
    fmt.Println("=====================")
	return nil
}

type Address struct {
	Street		string
	City		string
	Zip			string
}

type User struct {
	ID 		uint
	Name	string
	Email	string
	Address Address
}
func main(){
	u := User{
		ID: 1,
		Name: "Banna",
		Email: "smasayedalbanna75@gmail.com",
		Address: Address{
			Street: "66/2/2 Maniknagar, Wasa Road",
			City: "Dhaka",
			Zip: "1203",
		},
	}

	fmt.Println(u.Name)
	fmt.Println(u.Address.City)
	good := Product{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10, Discount: .10}
    bad  := Product{ID: 2, Name: "", Price: -5, Stock: 0}

    if err := good.Validate(); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("✓ Product is valid:", good.Name)
    }

    if err := bad.Validate(); err != nil {
        fmt.Println("✗ Invalid product:", err)
    }

	good.printReceipt(3)
}