package fakers

import (
	"github.com/Aktollkynn/GoProject.git/app/models"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func ProductFaker(db *gorm.DB) *models.Product {
	name := faker.Name()
	return &models.Product{
		ID:               uuid.New().String(),
		UserID:           "",
		Sku:              "",
		Name:             name,
		Slug:             slug.Make(name),
		Price:            decimal.Decimal{},
		Stock:            0,
		Weight:           decimal.NewFromFloat(rand.Float64()),
		ShortDescription: faker.Paragraph(),
		Description:      faker.Paragraph(),
		Status:           1,
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
		DeletedAt:        gorm.DeletedAt{},
	}
}
