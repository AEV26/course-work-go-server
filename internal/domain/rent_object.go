package domain

import (
	"fmt"
	"sort"
)

type RentObject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Area        float64  `json:"area"`
	Records     []Record `json:"records"`
}

var RecordNotFoundError = fmt.Errorf("Record not found")

type UpdateRentObjectInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Area        *float64 `json:"area"`
}

func NewRentObject(name string, description string, area float64) RentObject {
	return RentObject{
		Name:        name,
		Description: description,
		Area:        area,
		Records:     make([]Record, 0),
	}
}

func NewUpdateRentObjectInput(name string, description string, area float64) UpdateRentObjectInput {
	return UpdateRentObjectInput{
		Name:        &name,
		Description: &description,
		Area:        &area,
	}
}

func (r *RentObject) AddRecord(record Record) int {
	r.Records = append(r.Records, record)
	sort.Slice(r.Records, func(i, j int) bool {
		return r.Records[i].Date.Compare(r.Records[j].Date) < 0
	})
	return len(r.Records) - 1
}

func (r *RentObject) DeleteRecord(recordIndex int) error {
	if err := r.checkRecordIndex(recordIndex); err != nil {
		return err
	}

	size := len(r.Records)
	r.Records[recordIndex] = r.Records[size-1]
	r.Records = r.Records[:size-1]
	return nil

}

func (r *RentObject) UpdateRecord(recordIndex int, input UpdateRecordInput) error {
	if err := r.checkRecordIndex(recordIndex); err != nil {
		return err
	}

	r.Records[recordIndex] = r.Records[recordIndex].Update(input)
	return nil
}

func (r *RentObject) GetRecordByIndex(recordIndex int) (Record, error) {
	if err := r.checkRecordIndex(recordIndex); err != nil {
		return Record{}, err
	}

	return r.Records[recordIndex], nil
}

func (r *RentObject) checkRecordIndex(recordIndex int) error {
	size := len(r.Records)
	if size <= recordIndex || recordIndex < 0 {
		return RecordNotFoundError
	}
	return nil
}

func (r *RentObject) GetAllRecords() []Record {
	sort.Slice(r.Records, func(i, j int) bool {
		return r.Records[i].Date.Compare(r.Records[j].Date) <= 0
	})
	return r.Records
}

func (r *RentObject) Income() RUB {
	accumulator := func(income RUB, record Record) RUB {
		return income + record.Income()
	}

	return Reduce(r.Records, accumulator, RUB(0))
}

func (r *RentObject) Expenses() RUB {
	accumulator := func(expenses RUB, record Record) RUB {
		return expenses + record.Expenses()
	}
	return Reduce(r.Records, accumulator, RUB(0))
}

func (r *RentObject) Profit() RUB {
	accumulator := func(profit RUB, record Record) RUB {
		return profit + record.Profit()
	}
	return Reduce(r.Records, accumulator, RUB(0))
}

func (r *RentObject) Update(inp UpdateRentObjectInput) RentObject {
	newRentObject := *r

	if inp.Name != nil {
		newRentObject.Name = *inp.Name
	}

	if inp.Description != nil {
		newRentObject.Description = *inp.Description
	}

	if inp.Area != nil {
		newRentObject.Area = *inp.Area
	}

	return newRentObject

}

func Reduce[A, B any](collection []A, accumulator func(B, A) B, initalValue B) B {
	var result B
	for _, x := range collection {
		result = accumulator(result, x)
	}
	return result
}
