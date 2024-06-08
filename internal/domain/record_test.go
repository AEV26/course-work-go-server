package domain_test

import (
	"rental-server/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	record := domain.Record{
		Date:         time.Now(),
		Rent:         10000,
		Heat:         100,
		Exploitation: 100,
		MOP:          100,
		Renovation:   200,
		TBO:          300,
		Electricity:  400,
		EarthRent:    500,
		Other:        600,
		Security:     700,
	}

	t.Run("record should calculate expenses", func(t *testing.T) {

		got := record.Expenses()
		want := domain.RUB(300. + 200 + 300 + 400 + 500 + 600 + 700)

		assert.Equal(t, got, want, "got %f want %f", got, want)
	})

	t.Run("record should calculate income", func(t *testing.T) {
		got := record.Income()
		want := domain.RUB(10000)

		assert.Equal(t, got, want, "got %f want %f", got, want)
	})

	t.Run("record should calculate profit", func(t *testing.T) {
		got := record.Profit()
		want := domain.RUB(7000)
		assert.Equal(t, got, want, "got %f want %f", got, want)
	})

	t.Run("Should be able update record from UpdateRecordInput", func(t *testing.T) {
		record := domain.Record{}
		newTime := time.Now()
		newRent := domain.RUB(1000)
		newHeat := domain.RUB(1000)
		newExploitation := domain.RUB(1000)
		newMOP := domain.RUB(1000)
		newRenovation := domain.RUB(1000)
		newTBO := domain.RUB(1000)
		newElectricity := domain.RUB(1000)
		newEarthRent := domain.RUB(1000)
		newOther := domain.RUB(1000)
		newSecurity := domain.RUB(1000)

		inp := domain.UpdateRecordInput{
			Date:         &newTime,
			Rent:         &newRent,
			Heat:         &newHeat,
			Exploitation: &newExploitation,
			MOP:          &newMOP,
			Renovation:   &newRenovation,
			TBO:          &newTBO,
			Electricity:  &newElectricity,
			EarthRent:    &newEarthRent,
			Other:        &newOther,
			Security:     &newSecurity,
		}

		got := record.Update(inp)

		want := domain.Record{
			Date:         newTime,
			Rent:         newRent,
			Heat:         newHeat,
			Exploitation: newExploitation,
			MOP:          newMOP,
			Renovation:   newRenovation,
			TBO:          newTBO,
			Electricity:  newElectricity,
			EarthRent:    newEarthRent,
			Other:        newOther,
			Security:     newSecurity,
		}

		assert.Equal(t, want, got)

	})

}
