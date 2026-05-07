package repository

import "productmanager/entity"

type InMemoryRepo struct {
    products map[uint]entity.Product
    nextID   uint
}

func NewInMemoryRepo() *InMemoryRepo {
    r := &InMemoryRepo{products: make(map[uint]entity.Product), nextID: 1}
    // seed data
    r.products[1] = entity.Product{ID: 1, Name: "Laptop",   Price: 999.99, Stock: 5}
    r.products[2] = entity.Product{ID: 2, Name: "Mouse",    Price: 29.99,  Stock: 20}
    r.products[3] = entity.Product{ID: 3, Name: "Keyboard", Price: 79.99,  Stock: 15}
    r.nextID = 4
    return r
}

func (r *InMemoryRepo) Save(p entity.Product) (*entity.Product, error) {
    for _, existing := range r.products {
        if existing.Name == p.Name { return nil, entity.ErrDuplicate }
    }
    p.ID = r.nextID; r.products[p.ID] = p; r.nextID++
    return &p, nil
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

func (r *InMemoryRepo) Update(p entity.Product) (*entity.Product, error) {
    if _, ok := r.products[p.ID]; !ok { return nil, entity.ErrNotFound }
    r.products[p.ID] = p
    return &p, nil
}

func (r *InMemoryRepo) Delete(id uint) error {
    if _, ok := r.products[id]; !ok { return entity.ErrNotFound }
    delete(r.products, id)
    return nil
}