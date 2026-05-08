package usecase

import (
    "fmt"
    "productmanager/entity"
)

type productUsecase struct {
    repo entity.ProductRepository
}

func NewProductUsecase(repo entity.ProductRepository) entity.ProductUsecase {
    return &productUsecase{repo: repo}
}

func (uc *productUsecase) CreateProduct(name string, price float64, stock int) (*entity.Product, error) {
    p := entity.Product{Name: name, Price: price, Stock: stock}
    err := p.Validate()
    if err != nil {
        return nil, fmt.Errorf("CreateProduct: %w", err)
    }
    return uc.repo.Save(p)
}

func (uc *productUsecase) GetProduct(id uint) (*entity.Product, error) {
    p, err := uc.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetProduct id=%d: %w", id, err)
    }
    return p, nil
}

func (uc *productUsecase) ListProducts(page, limit int) ([]entity.Product, int) {
    if page < 1   { page  = 1 }
    if limit < 1  { limit = 10 }
    if limit > 100 { limit = 100 }

    all := uc.repo.FindAll()
    total := len(all)

    start := (page - 1) * limit
    if start >= total {
        return []entity.Product{}, total
    }

    end := start + limit
    if end > total { 
        end = total 
    }

    return all[start:end], total
}


func (uc *productUsecase) UpdateProduct(id uint, name string, price float64, stock int) (*entity.Product, error) {
    existing, err := uc.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("UpdateProduct: %w", err)
    }
    existing.Name  = name
    existing.Price = price
    existing.Stock = stock
    if err := existing.Validate(); err != nil {
        return nil, fmt.Errorf("UpdateProduct: %w", err)
    }
    return uc.repo.Update(*existing)
}

func (uc *productUsecase) DeleteProduct(id uint) error {
    err := uc.repo.Delete(id)
    if err != nil {
        return fmt.Errorf("DeleteProduct id=%d: %w", id, err)
    }
    return nil
}