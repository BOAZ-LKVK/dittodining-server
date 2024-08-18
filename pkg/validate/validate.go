package validate

import "github.com/go-playground/validator/v10"

// TODO: request validation middleware를 사용해서 로직 추상화하도록
var Validator = validator.New()
