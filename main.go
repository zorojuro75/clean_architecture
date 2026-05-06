package main

import (
    "errors"
    "fmt"

    "productmanager/entity"
    "productmanager/repository"
    "productmanager/usecase"
)

func main() {
    // Wire: repo → usecase (Dependency Injection)
    repo    := repository.NewInMemoryRepo()
    uc      := usecase.NewProductUsecase(repo)

    // Add products
    items := []struct{ name string; price float64; stock int }{
        {"Laptop",      999.99, 5},
        {"Mouse",       29.99,  20},
        {"",            10.00,  5},  // invalid
    }
    fmt.Println("▶ Adding products...")
    for _, item := range items {
        if err := uc.AddProduct(item.name, item.price, item.stock); err != nil {
            fmt.Printf("  ✗ %v\n", err)
        } else {
            fmt.Printf("  ✓ Added: %s\n", item.name)
        }
    }

    // List all
    fmt.Println("\n▶ All products:")
    for _, p := range uc.ListProducts() {
        fmt.Printf("  [%d] %-10s $%.2f\n", p.ID, p.Name, p.Price)
    }

    // Get one
    fmt.Println("\n▶ Looking up ID=1...")
    p, err := uc.GetProduct(1)
    if err != nil {
        if errors.Is(err, entity.ErrNotFound) {
            fmt.Println("  [404] not found")
        }
    } else {
        fmt.Printf("  Found: %s — $%.2f\n", p.Name, p.Price)
    }
}