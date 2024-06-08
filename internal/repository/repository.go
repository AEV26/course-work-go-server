package repository

import (
	"errors"
	"rental-server/internal/domain"
)

var ObjectNotFoundError = errors.New("Object not found")
var ObjectAlreadyExists = errors.New("Object exists")

type RentObjectRepository interface {
	Add(userID int64, object domain.RentObject) error
	Delete(userID int64, objectName string) error
	Update(userID int64, objectName string, object domain.UpdateRentObjectInput) error
	GetByName(userID int64, objectName string) (domain.RentObject, error)
	GetAll(userID int64) ([]domain.RentObject, error)

	AddRecord(userID int64, objectName string, record domain.Record) (int, error)
	DeleteRecord(userID int64, objectName string, recordIndex int) error
	UpdateRecord(userID int64, objectName string, recordIndex int, record domain.UpdateRecordInput) error
	GetRecordByIndex(userID int64, objectName string, recordIndex int) (domain.Record, error)
	GetAllRecords(userID int64, objectName string) ([]domain.Record, error)
}
