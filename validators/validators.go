package validators

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required"`
    Description string  `json:"description" validate:"required"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    BasePrice   float64 `json:"base_price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"required,gte=0"`
}

func ValidateCreateProduct(request CreateProductRequest) error {
    return validate.Struct(request)
}