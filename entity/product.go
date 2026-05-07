package entity

import "errors"

var (
    ErrNotFound     = errors.New("product not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrDuplicate    = errors.New("product already exists")
    ErrOutOfStock   = errors.New("out of stock")
)

type Product struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Stock int     `json:"stock"`
}

func (p Product) Validate() error {
    if p.Name == ""  { return ErrInvalidInput }
    if p.Price <= 0  { return ErrInvalidInput }
    if p.Stock < 0   { return ErrInvalidInput }
    return nil
}

type ProductRepository interface {
    Save(p Product) (*Product, error)
    FindByID(id uint) (*Product, error)
    FindAll() []Product
    Update(p Product) (*Product, error)
    Delete(id uint) error
}

type ProductUsecase interface {
    CreateProduct(name string, price float64, stock int) (*Product, error)
    GetProduct(id uint) (*Product, error)
    ListProducts() []Product
    UpdateProduct(id uint, name string, price float64, stock int) (*Product, error)
    DeleteProduct(id uint) error
}