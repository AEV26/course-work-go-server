package memory_test

import (
	"rental-server/internal/domain"
	"rental-server/internal/repository"
	"rental-server/internal/repository/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dummyRecord = domain.Record{}
var dummyObject = domain.RentObject{}
var dummyUserID int64 = 1

func TestAdd(t *testing.T) {
	t.Run("Happy path. Should be able to add object", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)

		object := domain.NewRentObject("Name", "Description", 1000)

		_ = rep.Add(dummyUserID, object)

		objects, _ := rep.GetAll(dummyUserID)

		assert.Contains(t, objects, object)
	})
}

func TestDelete(t *testing.T) {

	t.Run("Happy path. should be able to delete object", func(t *testing.T) {
		store := memory.MemoryStore{
			dummyUserID: {
				dummyObject.Name: dummyObject,
			},
		}
		rep := memory.NewMemoryObjectRepository(store)
		objects, _ := rep.GetAll(dummyUserID)
		assert.NotEmpty(t, objects)

		rep.Delete(dummyUserID, dummyObject.Name)

		objects, _ = rep.GetAll(dummyUserID)
		assert.Empty(t, objects)
	})

	t.Run("If object doesnt exist should return error", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		err := rep.Delete(dummyUserID, "")
		assert.ErrorIs(t, err, repository.ObjectNotFoundError)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Happy path. should be able to update object", func(t *testing.T) {
		store := memory.MemoryStore{
			dummyUserID: {
				dummyObject.Name: dummyObject,
			},
		}
		rep := memory.NewMemoryObjectRepository(store)

		newName := "Rodionova Street"
		newDescription := "HSE Campus"
		newArea := 1000.0

		update := domain.NewUpdateRentObjectInput(newName, newDescription, newArea)

		rep.Update(dummyUserID, dummyObject.Name, update)

		got, _ := rep.GetByName(dummyUserID, dummyObject.Name)

		want := domain.RentObject{
			Name:        newName,
			Description: newDescription,
			Area:        newArea,
		}

		assert.Equal(t, want, got)
	})

	t.Run("Should return an error if object doesnt exist", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		updateInput := domain.UpdateRentObjectInput{}
		err := rep.Update(dummyUserID, "", updateInput)

		assert.ErrorIs(t, err, repository.ObjectNotFoundError)
	})
}

func TestMemoryRepositoryRecords(t *testing.T) {

	t.Run("test add ", func(t *testing.T) {
		t.Run("Happy path. add record to object", func(t *testing.T) {
			record := domain.Record{}
			store := memory.MemoryStore{
				dummyUserID: {
					dummyObject.Name: {
						Records: []domain.Record{record},
					},
				},
			}

			rep := memory.NewMemoryObjectRepository(store)

			records, _ := rep.GetAllRecords(dummyUserID, dummyObject.Name)

			assert.Contains(t, records, record)
		})

		t.Run("If object doesnt exist should return an error", func(t *testing.T) {

			rep := memory.NewMemoryObjectRepository(nil)
			record := domain.Record{}
			_, err := rep.AddRecord(dummyUserID, "", record)

			assert.ErrorIs(t, err, repository.ObjectNotFoundError)
		})
	})

	t.Run("test delete", func(t *testing.T) {
		t.Run("Happy path. Delete record from object", func(t *testing.T) {
			recordIndex := 0
			store := memory.MemoryStore{
				dummyUserID: {
					dummyObject.Name: domain.RentObject{Records: []domain.Record{0: dummyRecord}},
				},
			}
			rep := memory.NewMemoryObjectRepository(store)

			rep.DeleteRecord(dummyUserID, dummyObject.Name, recordIndex)

			records, _ := rep.GetAllRecords(dummyUserID, dummyObject.Name)
			assert.Empty(t, records)
		})

		t.Run("If object doesnt exists should return ObjectNotFoundError", func(t *testing.T) {
			rep := memory.NewMemoryObjectRepository(nil)
			recordIndex := 0
			err := rep.DeleteRecord(dummyUserID, dummyObject.Name, recordIndex)

			assert.ErrorIs(t, err, repository.ObjectNotFoundError)
		})
		t.Run("If record doesnt exists should return RecordNotFoundError", func(t *testing.T) {
			rep := memory.NewMemoryObjectRepository(nil)
			_ = rep.Add(dummyUserID, dummyObject)
			recordIndex := 0

			err := rep.DeleteRecord(dummyUserID, dummyObject.Name, recordIndex)

			assert.ErrorIs(t, err, domain.RecordNotFoundError)
		})
	})

	t.Run("Test update record", func(t *testing.T) {
		t.Run("Happy path. update record", func(t *testing.T) {
			recordIndex := 0
			store := memory.MemoryStore{
				dummyUserID: {
					dummyObject.Name: domain.RentObject{Records: []domain.Record{0: dummyRecord}},
				},
			}
			rep := memory.NewMemoryObjectRepository(store)

			newRent := domain.RUB(1000)
			update := domain.UpdateRecordInput{
				Rent: &newRent,
			}

			rep.UpdateRecord(dummyUserID, dummyObject.Name, recordIndex, update)

			got, _ := rep.GetRecordByIndex(dummyUserID, dummyObject.Name, recordIndex)

			want := domain.Record{
				Rent: newRent,
			}

			assert.Equal(t, want, got)
		})
		t.Run("If object doesnt exists should return ObjectNotFoundError", func(t *testing.T) {
			rep := memory.NewMemoryObjectRepository(nil)
			objectID, recordIndex := "", 0
			updateInput := domain.UpdateRecordInput{}

			err := rep.UpdateRecord(dummyUserID, objectID, recordIndex, updateInput)

			assert.ErrorIs(t, err, repository.ObjectNotFoundError)
		})

		t.Run("If record doesnt exists should return RecordNotFoundError", func(t *testing.T) {
			rep := memory.NewMemoryObjectRepository(nil)
			_ = rep.Add(dummyUserID, dummyObject)
			recordIndex := 0
			updateInput := domain.UpdateRecordInput{}

			err := rep.UpdateRecord(dummyUserID, dummyObject.Name, recordIndex, updateInput)

			assert.ErrorIs(t, err, domain.RecordNotFoundError)
		})
	})
}
