package entity

import "errors"   // only stdlib — no project imports!

var (
    ErrNotFound     = errors.New("product not found")
    ErrOutOfStock   = errors.New("out of stock")
    ErrInvalidInput = errors.New("invalid input")
)

type Product struct {
    ID    uint
    Name  string
    Price float64
    Stock int
}

func (p Product) Validate() error {
    if p.Name == ""  { return ErrInvalidInput }
    if p.Price <= 0  { return ErrInvalidInput }
    return nil
}

func (p Product) IsAvailable() bool { return p.Stock > 0 }

// Repository interface — owned by the domain layer (Day 5)
type ProductRepository interface {
    Save(p Product) error
    FindByID(id uint) (*Product, error)
    FindAll() []Product
}