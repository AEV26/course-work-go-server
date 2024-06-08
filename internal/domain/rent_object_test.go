package domain_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"rental-server/internal/domain"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRentObject(t *testing.T) {
	t.Run("add record to rent object", func(t *testing.T) {
		rentObject := domain.RentObject{}
		record := domain.Record{}
		rentObject.AddRecord(record)

		records := rentObject.GetAllRecords()

		assert.Contains(t, records, record, "Records slice '%v' should contain record '%v'", records, record)
	})

	t.Run("should return sorted slice of records", func(t *testing.T) {
		rentObject := domain.RentObject{}

		for i := 0; i < 10; i++ {
			date := time.Date(10-i, 0, 0, 0, 0, 0, 0, time.Local)
			rentObject.AddRecord(domain.Record{
				Date: date,
			})
		}

		records := rentObject.GetAllRecords()
		isSorted := !sort.SliceIsSorted(records, func(i, j int) bool {
			return records[i].Date.Compare(records[j].Date) <= 0

		})

		assert.Falsef(t, isSorted, "Slice %v should be sorted by date in asc", records)
	})

	t.Run("should calculate income from all records", func(t *testing.T) {
		rentObject := domain.RentObject{}

		for i := 0; i < 10; i++ {
			rentObject.AddRecord(domain.Record{
				Rent: 10000,
			})
		}

		got := rentObject.Income()
		want := domain.RUB(100000)

		assert.Equal(t, want, got, "got %v want %v", got, want)
	})

	t.Run("should calculate expenses from all records", func(t *testing.T) {
		rentObject := domain.RentObject{}

		for i := 0; i < 10; i++ {
			rentObject.AddRecord(domain.Record{
				Heat:         1000,
				Exploitation: 1000,
				MOP:          1000,
				Renovation:   1000,
				TBO:          1000,
				Electricity:  1000,
				EarthRent:    1000,
				Other:        1000,
				Security:     1000,
			})
		}

		got := rentObject.Expenses()
		want := domain.RUB(90000)

		assert.Equal(t, want, got, "got %v want %v", got, want)
	})

	t.Run("should calculate profit from all records", func(t *testing.T) {
		rentObject := domain.RentObject{}

		for i := 0; i < 10; i++ {
			rentObject.AddRecord(domain.Record{
				Rent:         10000,
				Heat:         1000,
				Exploitation: 1000,
				MOP:          1000,
				Renovation:   1000,
				TBO:          1000,
				Electricity:  1000,
				EarthRent:    1000,
				Other:        1000,
				Security:     1000,
			})
		}
		got := rentObject.Profit()
		want := domain.RUB(10000)

		assert.Equal(t, want, got)
	})

	t.Run("should be able to update from UpdateRentObjectInput", func(t *testing.T) {
		rentObject := domain.NewRentObject("Rodionova", "HSE", 10000)
		newName := "Bolshaya pecherskaya"
		newDescription := "HSE Nizhny Novgorod"
		updateRentObject := domain.UpdateRentObjectInput{
			Name:        &newName,
			Description: &newDescription,
		}

		newObject := rentObject.Update(updateRentObject)

		assert.Equal(t, newName, newObject.Name)
		assert.Equal(t, newDescription, newObject.Description)
		assert.Equal(t, rentObject.Area, newObject.Area)
	})

	t.Run("should be marshalled with specific field names", func(t *testing.T) {
		rentObject := domain.NewRentObject("Rodionova", "HSE", 10000)

		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(rentObject)
		assert.NoError(t, err)

		want := fmt.Sprintf(`{"name": %q, "description": %q, "area": %f, "records": %v}`, rentObject.Name, rentObject.Description, rentObject.Area, rentObject.Records)
		assert.JSONEq(t, want, buf.String())
	})
}
