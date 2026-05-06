package usecase

import (
    "fmt"
    "productmanager/entity"
)

type ProductUsecase struct {
    repo entity.ProductRepository  // interface from entity package
}

func NewProductUsecase(repo entity.ProductRepository) *ProductUsecase {
    return &ProductUsecase{repo: repo}
}

func (uc *ProductUsecase) AddProduct(name string, price float64, stock int) error {
    p := entity.Product{Name: name, Price: price, Stock: stock}
    if err := p.Validate(); err != nil {
        return fmt.Errorf("AddProduct: %w", err)
    }
    return uc.repo.Save(p)
}

func (uc *ProductUsecase) GetProduct(id uint) (*entity.Product, error) {
    p, err := uc.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetProduct id=%d: %w", id, err)
    }
    return p, nil
}

func (uc *ProductUsecase) ListProducts() []entity.Product {
    return uc.repo.FindAll()
}