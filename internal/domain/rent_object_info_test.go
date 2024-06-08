package domain_test

import (
	"rental-server/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	object := domain.RentObject{
		Area: 100,
		Records: []domain.Record{
			{Rent: 1000, EarthRent: 500},
			{Rent: 1000, EarthRent: 500},
			{Rent: 1000, EarthRent: 500},
			{Rent: 1000, EarthRent: 500},
		},
	}

	got := domain.NewRentObjectInfo(object)

	want := domain.RentObjectInfo{
		Name:        object.Name,
		Description: object.Description,
		Area:        100,
		RecordsInfo: []domain.RecordInfo{
			{Record: domain.Record{Rent: 1000, EarthRent: 500}, Income: 1000, Expenses: 500, Profit: 500, IncomeByArea: 10, ExpensesByArea: 5, ProfitByArea: 5},
			{Record: domain.Record{Rent: 1000, EarthRent: 500}, Income: 1000, Expenses: 500, Profit: 500, IncomeByArea: 10, ExpensesByArea: 5, ProfitByArea: 5},
			{Record: domain.Record{Rent: 1000, EarthRent: 500}, Income: 1000, Expenses: 500, Profit: 500, IncomeByArea: 10, ExpensesByArea: 5, ProfitByArea: 5},
			{Record: domain.Record{Rent: 1000, EarthRent: 500}, Income: 1000, Expenses: 500, Profit: 500, IncomeByArea: 10, ExpensesByArea: 5, ProfitByArea: 5},
		},
	}

	assert.Equal(t, want, got)
}
