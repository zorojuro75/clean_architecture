package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound     = errors.New("product not found")
	ErrOutOfStock   = errors.New("product out of stock")
	ErrInvalidInput = errors.New("invalid input")
	ErrDuplicate    = errors.New("product already exists")
)

type Product struct {
	ID        uint
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
}

func (p Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required: %w", ErrInvalidInput)
	}
	if p.Price <= 0 {
		return fmt.Errorf("price is required: %w", ErrInvalidInput)
	}
	if p.Stock <= 0 {
		return fmt.Errorf("stock is required: %w", ErrInvalidInput)
	}
	return nil
}

func (p Product) isAvailable() bool {
	return p.Stock > 0
}

type Order struct {
	ID        uint
	Product   Product
	Quantity  int
	Total     float64
	CreatedAt time.Time
}

type ProductRepository interface {
	Save(p Product) error
	FindByID(id uint) (*Product, error)
	FindAll() []Product
	Update(p Product) error
}

type InMemoryRepo struct {
	products map[uint]Product
	nextID   uint
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make(map[uint]Product),
		nextID:   1,
	}
}

func (r *InMemoryRepo) Save(p Product) error {
	for _, existing := range r.products {
		if existing.Name == p.Name {
			return ErrDuplicate
		}
	}

	p.ID = r.nextID
	r.nextID++
	p.CreatedAt = time.Now()
	r.products[p.ID] = p

	return nil
}

func (r *InMemoryRepo) FindByID(ID uint) (*Product, error) {
	p, ok := r.products[ID]
	if !ok {
		return nil, ErrNotFound
	}
	return &p, nil
}

func (r *InMemoryRepo) FindAll() []Product {
	list := make([]Product, 0, len(r.products))
	for _, p := range r.products {
		list = append(list, p)
	}
	return list
}

func (r *InMemoryRepo) Update(p Product) error {
	_, ok := r.products[p.ID]
	if !ok {
		return ErrNotFound
	}
	r.products[p.ID] = p
	return nil
}

type ProductUseCase struct {
	repo    ProductRepository
	oderSeq uint
}

func NewProductUseCase(repo ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) AddProduct(name string, price float64, stock int) (*Product, error) {
	p := Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}
	err := p.Validate()
	if err != nil {
		return nil, fmt.Errorf("AddProduct: %w", err)
	}
	err = uc.repo.Save(p)
	if err != nil {
		return nil, fmt.Errorf("AddProduct: %w", err)
	}
	return &p, nil
}

func (uc *ProductUseCase) GetProduct(id uint) (*Product, error) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetProduct: %w", err)
	}
	return p, nil
}

func (uc *ProductUseCase) ListProducts() []Product {
	return uc.repo.FindAll()
}

func (uc *ProductUseCase) PlaceOrder(productID uint, quantity int) (*Order, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("PlaceOrder: quantity=%d, %w", quantity, ErrInvalidInput)
	}

	p, err := uc.repo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("PlaceOrder: %w", err)
	}

	if !p.isAvailable() {
		return nil, fmt.Errorf("PlaceOrder: %w", ErrOutOfStock)
	}
	if p.Stock < quantity {
		return nil, fmt.Errorf("PlaceOrder: stock has %d items, %w", p.Stock, ErrOutOfStock)
	}
	p.Stock -= quantity
	errr := uc.repo.Update(*p)
	if errr != nil {
		return nil, fmt.Errorf("PlaceOrder update stock: %w", errr)
	}
	uc.oderSeq++
	order := &Order{
		ID:        uc.oderSeq,
		Product:   *p,
		Quantity:  quantity,
		Total:     float64(quantity) * p.Price,
		CreatedAt: time.Now(),
	}
	return order, nil
}

func printProduct(p Product) {
	status := "✗ out of stock"
	if p.isAvailable() {
		status = fmt.Sprintf("✓ %d in stock", p.Stock)
	}
	fmt.Printf("  [%d] %-22s $%7.2f  %s\n", p.ID, p.Name, p.Price, status)
}

func printOrder(o Order) {
	fmt.Printf("  Order #%d — %s x%d — Total: $%.2f\n",
		o.ID, o.Product.Name, o.Quantity, o.Total)
}

func printErr(context string, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Printf("  [404] %s: %v\n", context, err)
	case errors.Is(err, ErrOutOfStock):
		fmt.Printf("  [409] %s: %v\n", context, err)
	case errors.Is(err, ErrInvalidInput):
		fmt.Printf("  [400] %s: %v\n", context, err)
	case errors.Is(err, ErrDuplicate):
		fmt.Printf("  [409] %s: %v\n", context, err)
	default:
		fmt.Printf("  [500] %s: %v\n", context, err)
	}
}

func main() {
	// 1. Create repo → inject into usecase (Dependency Injection)
	repo := NewInMemoryRepo()
	usecase := NewProductUseCase(repo)

	fmt.Println("════════════════════════════════════")
	fmt.Println(" 📦 Product Manager — Week 1 Project")
	fmt.Println("════════════════════════════════════")

	// 2. Add products
	fmt.Println("\n▶ Adding products...")
	products := []struct {
		name  string
		price float64
		stock int
	}{
		{"Wireless Headphones", 49.99, 10},
		{"Laptop Pro", 999.99, 5},
		{"Mechanical Keyboard", 79.99, 0},
		{"USB-C Hub", 29.99, 20},
		{"", 10.00, 5},            // invalid — no name
		{"Laptop Pro", 999.99, 3}, // duplicate
	}
	for _, item := range products {
		_, err := usecase.AddProduct(item.name, item.price, item.stock)
		if err != nil {
			printErr("AddProduct", err)
		} else {
			fmt.Printf("  ✓ Added: %s\n", item.name)
		}
	}

	// 3. List all products
	fmt.Println("\n▶ All products:")
	for _, p := range usecase.ListProducts() {
		printProduct(p)
	}

	// 4. Find one product
	fmt.Println("\n▶ Finding product ID=2...")
	p, err := usecase.GetProduct(2)
	if err != nil {
		printErr("GetProduct", err)
	} else {
		printProduct(*p)
	}

	// 5. Place orders — success and failure cases
	fmt.Println("\n▶ Placing orders...")
	orders := []struct {
		id   uint
		qty  int
		desc string
	}{
		{1, 3, "3x Headphones"},
		{3, 1, "1x Keyboard (out of stock)"},
		{99, 1, "product not found"},
		{1, -1, "invalid quantity"},
		{2, 2, "2x Laptop"},
	}
	for _, o := range orders {
		order, err := usecase.PlaceOrder(o.id, o.qty)
		if err != nil {
			printErr(o.desc, err)
		} else {
			printOrder(*order)
		}
	}

	// 6. List again — stock should be reduced
	fmt.Println("\n▶ Products after orders (check stock changes):")
	for _, p := range usecase.ListProducts() {
		printProduct(p)
	}

	fmt.Println("\n════════════════════════════════════")
	fmt.Println(" ✓ Week 1 complete. You built this.")
	fmt.Println("════════════════════════════════════")
}
