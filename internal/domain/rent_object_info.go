package domain

type RentObjectInfo struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Area        float64      `json:"area"`
	RecordsInfo []RecordInfo `json:"records_info"`
}

type RecordInfo struct {
	Record
	Income         RUB `json:"income"`
	Expenses       RUB `json:"expenses"`
	Profit         RUB `json:"profit"`
	IncomeByArea   RUB `json:"income_by_area"`
	ExpensesByArea RUB `json:"expenses_by_area"`
	ProfitByArea   RUB `json:"profit_by_area"`
}

func NewRentObjectInfo(object RentObject) RentObjectInfo {
	objectInfo := RentObjectInfo{
		Name:        object.Name,
		Description: object.Description,
		Area:        object.Area,
	}

	for _, record := range object.GetAllRecords() {
		var incomeByArea, expensesByArea, profitByArea RUB

		if object.Area != 0 {
			incomeByArea = RUB(float64(record.Income()) / object.Area)
			expensesByArea = RUB(float64(record.Expenses()) / object.Area)
			profitByArea = RUB(float64(record.Profit()) / object.Area)
		}

		recordInfo := RecordInfo{
			Record:         record,
			Income:         record.Income(),
			Expenses:       record.Expenses(),
			Profit:         record.Profit(),
			IncomeByArea:   incomeByArea,
			ExpensesByArea: expensesByArea,
			ProfitByArea:   profitByArea,
		}
		objectInfo.RecordsInfo = append(objectInfo.RecordsInfo, recordInfo)
	}

	return objectInfo
}
