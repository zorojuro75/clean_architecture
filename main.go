package main 
import "fmt"

type Product struct {
	ID			uint
	Name		string
	Price		float64
	Stock		int
	Discount	float64
}

func (p *Product) applyDiscount() (float64, error){
	if p.Discount < 0 || p.Discount > 1 {
		return 0, fmt.Errorf("Invalid discount percentage")
	}
	discounted := p.Price * p.Discount
	return discounted, nil
}

func (p *Product)orderTotal(quantity int) float64 {
	p.updateStock(quantity)
	return p.Price * float64(quantity)
}

func (p *Product) updateStock (quantity int){
	p.Stock = p.Stock - quantity
}

func (p *Product) printReceipt(quantity int){
	dicountedPrice, err := p.applyDiscount()
	if err !=nil{
		fmt.Println("Error: ", err)
		return
	}
	total := p.orderTotal(quantity)
	
	fmt.Println("====== Receipt ======")
    fmt.Printf("Product  : %s\n", p.Name)
    fmt.Printf("Price    : $%.2f\n", p.Price)
    fmt.Printf("Discount : %.0f%%\n", p.Discount*100)
    fmt.Printf("New Price: $%.2f\n", dicountedPrice)
    fmt.Printf("Quantity : %d\n", p.Stock)
    fmt.Printf("Total    : $%.2f\n", total)
    fmt.Println("=====================")
}
func main(){
	p := Product{
		ID:			1,
		Name:		"Wireless Headphones",
		Price:		49.99,
		Stock: 		100,
		Discount:	0.15,
	}
	p.printReceipt(5)
}