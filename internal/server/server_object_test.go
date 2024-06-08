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

var dummyObject = domain.NewRentObject("Name", "Description", 1000)
var dummyUserID int64 = 1

func TestAddObject(t *testing.T) {
	t.Run("Should add object", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		s := server.NewRentObjectServer(rep)

		request := newAddObjectRequest(dummyUserID, dummyObject)
		responce := httptest.NewRecorder()

		s.ServeHTTP(responce, request)
		assertStatus(t, responce.Code, http.StatusCreated)

		objects, _ := rep.GetAll(dummyUserID)
		if !assert.Len(t, objects, 1) {
			t.Fatal()
		}
	})

	t.Run("Should return UnprocessableEntity on wrong request body", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		s := server.NewRentObjectServer(rep)

		request := httptest.NewRequest(http.MethodPost, "/addObject", nil)
		responce := httptest.NewRecorder()
		s.ServeHTTP(responce, request)

		assertStatus(t, responce.Code, http.StatusUnprocessableEntity)
	})
}

func TestDeleteObject(t *testing.T) {
	t.Run("Should be able to delete object", func(t *testing.T) {
		store := memory.MemoryStore{
			dummyUserID: {
				dummyObject.Name: dummyObject,
			},
		}
		rep := memory.NewMemoryObjectRepository(store)

		s := server.NewRentObjectServer(rep)

		request := newDeleteObjectRequest(dummyUserID, dummyObject.Name)
		responce := httptest.NewRecorder()

		s.ServeHTTP(responce, request)
		assertStatus(t, responce.Code, http.StatusOK)

		objects, _ := rep.GetAll(dummyUserID)
		assert.Len(t, objects, 0)
	})
	t.Run("Should return UnprocessableEntity on wrong request", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		s := server.NewRentObjectServer(rep)

		request := httptest.NewRequest(http.MethodPost, "/deleteObject", nil)
		responce := httptest.NewRecorder()
		s.ServeHTTP(responce, request)

		assertStatus(t, responce.Code, http.StatusUnprocessableEntity)
	})
	t.Run("Should return NotFound if object doesnt exists", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		s := server.NewRentObjectServer(rep)

		request := newDeleteObjectRequest(dummyUserID, dummyObject.Name)
		responce := httptest.NewRecorder()
		s.ServeHTTP(responce, request)

		assertStatus(t, responce.Code, http.StatusNotFound)
	})
}

func TestUpdateObject(t *testing.T) {
	t.Run("Should update object", func(t *testing.T) {
		store := memory.MemoryStore{
			dummyUserID: {
				dummyObject.Name: dummyObject,
			},
		}
		rep := memory.NewMemoryObjectRepository(store)

		s := server.NewRentObjectServer(rep)

		updateInput := domain.NewUpdateRentObjectInput("NewName", "NewDescription", 10000)

		request := newUpdateObjectRequest(dummyUserID, dummyObject.Name, updateInput)
		responce := httptest.NewRecorder()

		s.ServeHTTP(responce, request)
		assertStatus(t, responce.Code, http.StatusOK)

		newObject := dummyObject.Update(updateInput)

		got, _ := rep.GetByName(dummyUserID, dummyObject.Name)

		assert.Equal(t, newObject, got)
	})
	t.Run("Should return UnprocessableEntity on wrong request", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		s := server.NewRentObjectServer(rep)

		request := httptest.NewRequest(http.MethodPost, "/updateObject", nil)
		responce := httptest.NewRecorder()
		s.ServeHTTP(responce, request)

		assertStatus(t, responce.Code, http.StatusUnprocessableEntity)
	})
	t.Run("Should return NotFound if object doesnt exists", func(t *testing.T) {
		rep := memory.NewMemoryObjectRepository(nil)
		s := server.NewRentObjectServer(rep)

		request := newUpdateObjectRequest(dummyUserID, dummyObject.Name, domain.UpdateRentObjectInput{})
		responce := httptest.NewRecorder()
		s.ServeHTTP(responce, request)

		assertStatus(t, responce.Code, http.StatusNotFound)
	})
}

func TestGetObjectByID(t *testing.T) {
	object := dummyObject
	rep := memory.NewMemoryObjectRepository(nil)
	_ = rep.Add(dummyUserID, object)

	s := server.NewRentObjectServer(rep)

	request := newGetObjectRequest(dummyUserID, object.Name)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)
	assertStatus(t, responce.Code, http.StatusOK)

	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(object)

	assert.JSONEq(t, buf.String(), responce.Body.String())
}

func TestGetAll(t *testing.T) {
	rep := memory.NewMemoryObjectRepository(nil)
	var objects []domain.RentObject
	for i := 0; i < 10; i++ {
		object := dummyObject
		object.Name = fmt.Sprintf("Name%d", i)
		_ = rep.Add(dummyUserID, object)
		objects = append(objects, object)
	}

	s := server.NewRentObjectServer(rep)

	request := newGetAllRequest(dummyUserID)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)
	assertStatus(t, responce.Code, http.StatusOK)

	var gotObjects []domain.RentObject
	json.NewDecoder(responce.Body).Decode(&gotObjects)

	assert.Equal(t, objects, gotObjects)
}

func TestGetObjectInfo(t *testing.T) {
	object := domain.RentObject{
		Area: 100,
		Records: []domain.Record{
			{Rent: 100, EarthRent: 50},
			{Rent: 100, EarthRent: 50},
			{Rent: 100, EarthRent: 50},
			{Rent: 100, EarthRent: 50},
		},
	}

	store := memory.MemoryStore{
		dummyUserID: {
			dummyObject.Name: object,
		},
	}
	rep := memory.NewMemoryObjectRepository(store)

	s := server.NewRentObjectServer(rep)

	request := newGetObjectInfoRequest(dummyUserID, dummyObject.Name)
	responce := httptest.NewRecorder()

	s.ServeHTTP(responce, request)

	assertStatus(t, responce.Code, http.StatusOK)

	var got domain.RentObjectInfo
	_ = json.NewDecoder(responce.Body).Decode(&got)

	want := domain.NewRentObjectInfo(object)

	assert.Equal(t, want, got)

}

func newAddObjectRequest(userId int64, object domain.RentObject) *http.Request {
	buf := &bytes.Buffer{}
	data := requests.AddObjectRequest{
		UserID: &userId,
		Object: &object,
	}
	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest(http.MethodPost, "/addObject", buf)
	return req
}

func newDeleteObjectRequest(userId int64, objectName string) *http.Request {
	buf := &bytes.Buffer{}

	data := requests.DeleteObjectRequest{
		UserID:     &userId,
		ObjectName: &objectName,
	}

	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest(http.MethodPost, "/deleteObject", buf)
	return req
}

func newUpdateObjectRequest(userId int64, objectName string, updateObjectInput domain.UpdateRentObjectInput) *http.Request {
	buf := &bytes.Buffer{}

	data := requests.UpdateObjectRequest{
		UserID:      &userId,
		ObjectName:  &objectName,
		UpdateInput: &updateObjectInput,
	}

	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest(http.MethodPost, "/updateObject", buf)
	return req
}

func newGetObjectRequest(userId int64, objectName string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/getObject?%s=%d&%s=%s", server.UserIdQueryParam, userId, server.ObjectNameQueryParam, objectName), nil)
	return req
}

func newGetAllRequest(userId int64) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/getAll?%s=%d", server.UserIdQueryParam, userId), nil)
	return req
}

func newGetObjectInfoRequest(userId int64, objectName string) *http.Request {
	path := fmt.Sprintf("/getObjectInfo?%s=%d&%s=%s", server.UserIdQueryParam, userId, server.ObjectNameQueryParam, objectName)
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
