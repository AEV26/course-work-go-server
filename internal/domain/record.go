package domain

import (
	"time"
)

type RUB float64

type Record struct {
	Date         time.Time `json:"date"`
	Rent         RUB       `json:"rent"`
	Heat         RUB       `json:"heat"`
	Exploitation RUB       `json:"exploitation"`
	MOP          RUB       `json:"mop"`
	Renovation   RUB       `json:"renovation"`
	TBO          RUB       `json:"tbo"`
	Electricity  RUB       `json:"electricity"`
	EarthRent    RUB       `json:"earth_rent"`
	Other        RUB       `json:"other"`
	Security     RUB       `json:"security"`
}

type UpdateRecordInput struct {
	Date         *time.Time `json:"date"`
	Rent         *RUB       `json:"rent"`
	Heat         *RUB       `json:"heat"`
	Exploitation *RUB       `json:"exploitation"`
	MOP          *RUB       `json:"mop"`
	Renovation   *RUB       `json:"renovation"`
	TBO          *RUB       `json:"tbo"`
	Electricity  *RUB       `json:"electricity"`
	EarthRent    *RUB       `json:"earth_rent"`
	Other        *RUB       `json:"other"`
	Security     *RUB       `json:"security"`
}

func (r *Record) Expenses() RUB {
	return r.Heat +
		r.Exploitation +
		r.MOP +
		r.Renovation +
		r.TBO +
		r.Electricity +
		r.EarthRent +
		r.Other +
		r.Security
}

func (r *Record) Income() RUB {
	return r.Rent
}

func (r *Record) Profit() RUB {
	return r.Income() - r.Expenses()
}

func (r *Record) Update(inp UpdateRecordInput) Record {
	newRecord := *r

	if inp.Date != nil {
		newRecord.Date = *inp.Date
	}

	if inp.Rent != nil {
		newRecord.Rent = *inp.Rent
	}

	if inp.Heat != nil {
		newRecord.Heat = *inp.Heat
	}

	if inp.Exploitation != nil {
		newRecord.Exploitation = *inp.Exploitation
	}

	if inp.MOP != nil {
		newRecord.MOP = *inp.MOP
	}

	if inp.Renovation != nil {
		newRecord.Renovation = *inp.Renovation
	}

	if inp.TBO != nil {
		newRecord.TBO = *inp.TBO
	}

	if inp.Electricity != nil {
		newRecord.Electricity = *inp.Electricity
	}

	if inp.EarthRent != nil {
		newRecord.EarthRent = *inp.EarthRent
	}

	if inp.Other != nil {
		newRecord.Other = *inp.Other
	}

	if inp.Security != nil {
		newRecord.Security = *inp.Security
	}
	return newRecord
}
