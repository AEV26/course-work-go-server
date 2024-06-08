package memory

import (
	"rental-server/internal/domain"
	"rental-server/internal/repository"
	"sort"
)

type MemoryStore map[int64]map[string]domain.RentObject

type MemoryObjectRepository struct {
	store MemoryStore
}

func NewMemoryObjectRepository(store MemoryStore) *MemoryObjectRepository {
	if store == nil {
		store = MemoryStore{}
	}

	return &MemoryObjectRepository{
		store: store,
	}
}

func (m *MemoryObjectRepository) Add(userID int64, object domain.RentObject) error {
	if m.store[userID] == nil {
		objectStore := make(map[string]domain.RentObject)
		m.store[userID] = objectStore
	}
	m.store[userID][object.Name] = object
	return nil
}

func (m *MemoryObjectRepository) Delete(userID int64, objectName string) error {
	if _, ok := m.store[userID][objectName]; !ok {
		return repository.ObjectNotFoundError
	}
	delete(m.store[userID], objectName)
	return nil
}

func (m *MemoryObjectRepository) Update(userID int64, objectName string, input domain.UpdateRentObjectInput) error {
	object, err := m.GetByName(userID, objectName)
	if err != nil {
		return err
	}

	newObject := object.Update(input)
	m.store[userID][objectName] = newObject
	return nil
}

func (m *MemoryObjectRepository) GetByName(userID int64, objectName string) (domain.RentObject, error) {
	var obj domain.RentObject
	obj, ok := m.store[userID][objectName]

	if !ok {
		return obj, repository.ObjectNotFoundError
	}

	return obj, nil
}

func (m *MemoryObjectRepository) GetAll(userID int64) ([]domain.RentObject, error) {
	var objects []domain.RentObject
	for _, object := range m.store[userID] {
		objects = append(objects, object)
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Name < objects[j].Name
	})
	return objects, nil
}

func (m *MemoryObjectRepository) AddRecord(userID int64, objectName string, record domain.Record) (int, error) {
	object, err := m.GetByName(userID, objectName)
	if err != nil {
		return 0, err
	}
	index := object.AddRecord(record)
	m.store[userID][objectName] = object
	return index, nil
}

func (m *MemoryObjectRepository) DeleteRecord(userID int64, objectName string, recordIndex int) error {
	object, err := m.GetByName(userID, objectName)
	if err != nil {
		return err
	}

	err = object.DeleteRecord(recordIndex)
	if err != nil {
		return err
	}
	m.store[userID][objectName] = object
	return nil
}

func (m *MemoryObjectRepository) UpdateRecord(userID int64, objectName string, recordIndex int, input domain.UpdateRecordInput) error {
	object, err := m.GetByName(userID, objectName)
	if err != nil {
		return err
	}
	err = object.UpdateRecord(recordIndex, input)
	if err != nil {
		return err
	}
	m.store[userID][objectName] = object
	return nil
}

func (m *MemoryObjectRepository) GetRecordByIndex(userID int64, objectName string, recordIndex int) (domain.Record, error) {
	var zero domain.Record
	object, err := m.GetByName(userID, objectName)
	if err != nil {
		return zero, err
	}

	return object.GetRecordByIndex(recordIndex)

}

func (m *MemoryObjectRepository) GetAllRecords(userID int64, objectName string) ([]domain.Record, error) {
	object, err := m.GetByName(userID, objectName)
	if err != nil {
		return nil, err
	}
	return object.GetAllRecords(), nil
}
