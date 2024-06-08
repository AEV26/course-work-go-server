package mongorep_test

import (
	"fmt"
	"math/rand"
	"rental-server/internal/domain"
	mongorep "rental-server/internal/repository/mongo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testURI = "mongodb://localhost:27017/test"
var testDatabase = "test"

var dummyUserId int64 = 1
var dummyObject = domain.NewRentObject("", "", 0)
var dummyRecord = domain.Record{}

func TestCreateMongoRepository(t *testing.T) {
	rep, err := mongorep.NewMongoDBRepository(testURI, testDatabase)
	assert.NoError(t, err)
	if rep == nil {
		t.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	object := dummyObject

	t.Run("Should be able to add object", func(t *testing.T) {
		err := rep.Add(dummyUserId, object)

		if !assert.NoError(t, err) {
			t.Fatal(err)
		}

		got, err := rep.GetByName(dummyUserId, object.Name)
		if !assert.NoError(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, object, got)
	})

	t.Run("Should return an error if object already exists", func(t *testing.T) {
		err := rep.Add(dummyUserId, object)
		assert.Error(t, err)
	})
	rep.Clear()
}

func TestDelete(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	object := dummyObject
	rep.Add(dummyUserId, object)

	t.Run("shoud delete object", func(t *testing.T) {
		err := rep.Delete(dummyUserId, object.Name)
		assert.NoError(t, err)
	})
	t.Run("should return object not found if object does not exists", func(t *testing.T) {
		err := rep.Delete(dummyUserId, object.Name)
		assert.Error(t, err)
	})

	rep.Clear()
}

func TestUpdate(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	object := dummyObject
	rep.Add(dummyUserId, object)

	newName := "test"
	input := domain.UpdateRentObjectInput{Name: &newName}
	updated := object.Update(input)

	t.Run("should update object", func(t *testing.T) {
		err := rep.Update(dummyUserId, object.Name, input)
		assert.NoError(t, err)

		got, err := rep.GetByName(dummyUserId, updated.Name)
		assert.NoError(t, err)

		assert.Equal(t, updated, got)
	})

	t.Run("should return object not found if object does not exists", func(t *testing.T) {
		err := rep.Update(dummyUserId, "failed", input)
		assert.Error(t, err)
	})
	rep.Clear()
}

func TestGetAll(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	var objects []domain.RentObject
	for i := 0; i < 10; i++ {
		object := domain.NewRentObject(fmt.Sprintf("%d", i), "", 0)
		objects = append(objects, object)
		rep.Add(dummyUserId, object)
	}

	got, err := rep.GetAll(dummyUserId)
	assert.NoError(t, err)

	assert.Equal(t, objects, got)

	rep.Clear()
}

func TestAddRecord(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, "testing")
	rep.Add(dummyUserId, dummyObject)
	t.Run("should add record to object", func(t *testing.T) {
		index, err := rep.AddRecord(dummyUserId, dummyObject.Name, dummyRecord)
		assert.NoError(t, err)
		got, err := rep.GetRecordByIndex(dummyUserId, dummyObject.Name, index)
		assert.NoError(t, err)

		dummyRecord.Heat = 1
		index, err = rep.AddRecord(dummyUserId, dummyObject.Name, dummyRecord)
		assert.NoError(t, err)

		got, err = rep.GetRecordByIndex(dummyUserId, dummyObject.Name, index)
		assert.NoError(t, err)
		assert.Equal(t, dummyRecord, got)
	})
	t.Run("Shoudl return an error if object does not exists", func(t *testing.T) {
		_, err := rep.AddRecord(dummyUserId, "WTH", dummyRecord)
		assert.Error(t, err)
	})
	rep.Clear()
}

func TestDeleteRecord(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	rep.Add(dummyUserId, dummyObject)
	index, _ := rep.AddRecord(dummyUserId, dummyObject.Name, dummyRecord)

	t.Run("should delete record", func(t *testing.T) {
		err := rep.DeleteRecord(dummyUserId, dummyObject.Name, index)
		assert.NoError(t, err)

		_, err = rep.GetRecordByIndex(dummyUserId, dummyObject.Name, index)
		assert.Error(t, err)
	})
	rep.Clear()
}

func TestUpdateRecord(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	rep.Add(dummyUserId, dummyObject)
	index, _ := rep.AddRecord(dummyUserId, dummyObject.Name, dummyRecord)
	newRent := domain.RUB(1000)
	input := domain.UpdateRecordInput{Rent: &newRent}
	newRecord := dummyRecord.Update(input)

	t.Run("should update record", func(t *testing.T) {
		err := rep.UpdateRecord(dummyUserId, dummyObject.Name, index, input)
		assert.NoError(t, err)

		got, err := rep.GetRecordByIndex(dummyUserId, dummyObject.Name, index)
		assert.Equal(t, newRecord, got)
	})
	rep.Clear()
}

func TestGetAllRecords(t *testing.T) {
	rep, _ := mongorep.NewMongoDBRepository(testURI, testDatabase)
	rep.Add(dummyUserId, dummyObject)
	var records []domain.Record
	for i := 0; i < 10; i++ {
		randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
		record := domain.Record{Date: time.Unix(randomTime, 0)}
		records = append(records, record)
		rep.AddRecord(dummyUserId, dummyObject.Name, record)
	}

	t.Run("should return all records", func(t *testing.T) {
		got, err := rep.GetAllRecords(dummyUserId, dummyObject.Name)
		assert.NoError(t, err)
		assert.Len(t, got, 10)

		for i := 1; i < len(got); i++ {
			if got[i-1].Date.Compare(got[i].Date) == 1 {
				t.Fail()
			}
		}
	})
	rep.Clear()
}
