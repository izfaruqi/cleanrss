package controllers

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func ValidatorInit() {
	validate = validator.New()
}