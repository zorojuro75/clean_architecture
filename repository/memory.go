package repository

import "productmanager/entity"  // imports the entity package

type InMemoryRepo struct {
    products map[uint]entity.Product
    nextID   uint
}

func NewInMemoryRepo() *InMemoryRepo {
    return &InMemoryRepo{products: make(map[uint]entity.Product), nextID: 1}
}

func (r *InMemoryRepo) Save(p entity.Product) error {
    p.ID = r.nextID
    r.products[p.ID] = p
    r.nextID++
    return nil
}

func (r *InMemoryRepo) FindByID(id uint) (*entity.Product, error) {
    p, ok := r.products[id]
    if !ok { return nil, entity.ErrNotFound }
    return &p, nil
}

func (r *InMemoryRepo) FindAll() []entity.Product {
    list := []entity.Product{}
    for _, p := range r.products { list = append(list, p) }
    return list
}