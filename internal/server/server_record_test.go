package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rental-server/internal/domain"
	"rental-server/internal/repository/memory"
	"rental-server/internal/server"
	"rental-server/internal/server/requests"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dummyRecord = domain.Record{}

func TestAddRecord(t *testing.T) {
	t.Run("Happy path. Add record", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		_ = rep.Add(dummyUserID, dummyObject)

		s := server.NewRentObjectServer(rep)

		request := newAddRecordRequest(dummyUserID, dummyObject.Name, dummyRecord)
		responce := httptest.NewRecorder()

		s.ServeHTTP(responce, request)
		assertStatus(t, responce.Code, http.StatusOK)

		object, _ := rep.GetByName(dummyUserID, dummyObject.Name)

		assert.Contains(t, object.Records, dummyRecord)
	})
}

func TestDeleteRecord(t *testing.T) {
	rep := memory.NewMemoryObjectRepository(nil)
	_ = rep.Add(dummyUserID, dummyObject)
	recordPos, _ := rep.AddRecord(dummyUserID, dummyObject.Name, dummyRecord)

	s := server.NewRentObjectServer(rep)

	request := newDeleteRecordRequest(dummyUserID, dummyObject.Name, recordPos)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)
	assertStatus(t, responce.Code, http.StatusOK)

	object, _ := rep.GetByName(dummyUserID, dummyObject.Name)
	assert.NotContains(t, object.Records, dummyRecord)

}

func TestUpdateRecord(t *testing.T) {
	rep := memory.NewMemoryObjectRepository(nil)
	object := dummyObject
	recordPos := object.AddRecord(dummyRecord)
	_ = rep.Add(dummyUserID, object)

	s := server.NewRentObjectServer(rep)

	newRent := domain.RUB(1000)
	updateInput := domain.UpdateRecordInput{Rent: &newRent}

	request := newUpdateRecordRequest(dummyUserID, dummyObject.Name, recordPos, updateInput)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)
	assertStatus(t, responce.Code, http.StatusOK)

	want := dummyRecord.Update(updateInput)
	got, _ := rep.GetRecordByIndex(dummyUserID, dummyObject.Name, recordPos)

	assert.Equal(t, got, want)
}

func TestGetRecordByIndex(t *testing.T) {
	rep := memory.NewMemoryObjectRepository(nil)
	object := dummyObject
	record := domain.Record{Rent: 1000}
	recordPos := object.AddRecord(record)
	_ = rep.Add(dummyUserID, object)

	s := server.NewRentObjectServer(rep)

	request := newGetRecordRequest(dummyUserID, dummyObject.Name, recordPos)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)
	assertStatus(t, responce.Code, http.StatusOK)

	var got domain.Record
	json.NewDecoder(responce.Body).Decode(&got)

	assert.Equal(t, record, got)
}

func TestGetRecords(t *testing.T) {
	rep := memory.NewMemoryObjectRepository(nil)
	object := dummyObject

	object.AddRecord(dummyRecord)
	object.AddRecord(dummyRecord)
	_ = rep.Add(dummyUserID, object)

	s := server.NewRentObjectServer(rep)

	request := newGetRecordsRequest(dummyUserID, dummyObject.Name)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)
	assertStatus(t, responce.Code, http.StatusOK)

	var got []domain.Record
	json.NewDecoder(responce.Body).Decode(&got)

	assert.Len(t, got, 2)
}

func newAddRecordRequest(userID int64, objectName string, record domain.Record) *http.Request {
	buf := &bytes.Buffer{}

	data := requests.AddRecordRequest{
		UserID:     &userID,
		ObjectName: &objectName,
		Record:     &record,
	}
	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest(http.MethodPost, "/addRecord", buf)
	return req
}

func newDeleteRecordRequest(userID int64, objectName string, recordIndex int) *http.Request {
	buf := &bytes.Buffer{}

	data := requests.DeleteRecordRequest{
		UserID:      &userID,
		ObjectName:  &objectName,
		RecordIndex: &recordIndex,
	}
	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest(http.MethodPost, "/deleteRecord", buf)
	return req
}

func newUpdateRecordRequest(userID int64, objectName string, recordIndex int, input domain.UpdateRecordInput) *http.Request {
	buf := &bytes.Buffer{}

	data := requests.UpdateRecordRequest{
		UserID:      &userID,
		ObjectName:  &objectName,
		RecordIndex: &recordIndex,
		UpdateInput: &input,
	}
	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest(http.MethodPost, "/updateRecord", buf)
	return req
}

func newGetRecordRequest(userID int64, objectName string, recordIndex int) *http.Request {
	uri := fmt.Sprintf(
		"/getRecord?%s=%d&%s=%s&%s=%d", server.UserIdQueryParam, userID, server.ObjectNameQueryParam, objectName, server.RecordIndexQueryParam, recordIndex,
	)

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return req
}

func newGetRecordsRequest(userID int64, objectName string) *http.Request {
	uri := fmt.Sprintf("/getRecords?%s=%d&%s=%s", server.UserIdQueryParam, userID, server.ObjectNameQueryParam, objectName)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return req
}
