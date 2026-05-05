package main

import (
	"fmt"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrStockOut = errors.New("Out of stock")
	ErrInvalidInput = errors.New("Invalid input")
)

type Product struct {
	ID		uint
	Name	string
	Price 	float64
	Stock 	int
}

type ProductRepository interface {
	FindByID(id uint) (*Product, error)
}

type FakeRepo struct {}

func (r *FakeRepo) FindByID(id uint) (*Product, error) {
	products := map[uint]*Product{
        1: {ID: 1, Name: "Laptop",    Price: 999.99, Stock: 5},
        2: {ID: 2, Name: "Headphones", Price: 49.99,  Stock: 0},
    }

	p, ok := products[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func placeOrder(r ProductRepository, productID uint, quantity int) (float64, error){
	if quantity<=0{
		return 0, fmt.Errorf("placeOrder: quantity=%d: %w", quantity, ErrInvalidInput)
	}
	p, err:= r.FindByID(productID)
	if err!=nil{
		return 0, fmt.Errorf("placeOrder: productID=%d: %w", productID, err)
	}

	if p.Stock < quantity {
		return 0, fmt.Errorf("placeOrder: %s: %w", p.Name, ErrStockOut)
	}

	total := float64(quantity) * p.Price
	return total, nil
}

func handleOrder(r ProductRepository, productID uint, quantity int) {
	total, err := placeOrder(r, productID, quantity)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
            fmt.Printf("[404] product %d does not exist\n", productID)
        case errors.Is(err, ErrStockOut):
            fmt.Printf("[409] product %d is out of stock\n", productID)
        case errors.Is(err, ErrInvalidInput):
            fmt.Printf("[400] invalid quantity: %d\n", quantity)
        default:
            fmt.Printf("[500] internal error: %v\n", err)
        }
		return
	}
	fmt.Printf("[200] order placed! total: $%.2f\n", total)
}


func main() {
    repo := &FakeRepo{}

    fmt.Println("--- Order attempts ---")
    handleOrder(repo, 1, 2)  // success
    handleOrder(repo, 2, 1)  // out of stock
    handleOrder(repo, 9, 1)  // not found
    handleOrder(repo, 1, -1) // invalid input
}