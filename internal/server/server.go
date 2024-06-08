package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"rental-server/internal/domain"
	"rental-server/internal/repository"
	"rental-server/internal/server/requests"
	"strconv"
)

var UserIdQueryParam = "userId"
var ObjectNameQueryParam = "objectName"
var RecordIndexQueryParam = "recordIndex"

type appHandler func(w http.ResponseWriter, r *http.Request) *appError

type appError struct {
	Error error
	Msg   string
	Code  int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Msg, err.Code)
	}
}

type RentObjectServer struct {
	rep repository.RentObjectRepository
	http.Handler
}

func NewRentObjectServer(rep repository.RentObjectRepository) *RentObjectServer {
	server := &RentObjectServer{
		rep: rep,
	}

	router := http.NewServeMux()
	router.Handle("/addObject", appHandler(server.addObject))
	router.Handle("/deleteObject", appHandler(server.deleteObject))
	router.Handle("/updateObject", appHandler(server.updateObject))
	router.Handle("/getObject", appHandler(server.getObject))
	router.Handle("/getObjectInfo", appHandler(server.getObjectInfo))
	router.Handle("/getAll", appHandler(server.getAll))
	router.Handle("/addRecord", appHandler(server.addRecord))
	router.Handle("/deleteRecord", appHandler(server.deleteRecord))
	router.Handle("/updateRecord", appHandler(server.updateRecord))
	router.Handle("/getRecord", appHandler(server.getRecord))
	router.Handle("/getRecords", appHandler(server.getRecords))

	server.Handler = router

	return server
}

func (s *RentObjectServer) addObject(w http.ResponseWriter, r *http.Request) *appError {
	var addObjectRequest requests.AddObjectRequest

	if err := parseRequest(r.Body, &addObjectRequest); err != nil {
		return &appError{err, "Error while parsing body", http.StatusUnprocessableEntity}
	}

	err := s.rep.Add(*addObjectRequest.UserID, *addObjectRequest.Object)
	if err != nil {
		return processRepositoryError(err)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *RentObjectServer) deleteObject(w http.ResponseWriter, r *http.Request) *appError {
	var deleteObjectRequest requests.DeleteObjectRequest

	if err := parseRequest(r.Body, &deleteObjectRequest); err != nil {
		return &appError{err, "Error while parsing body", http.StatusUnprocessableEntity}
	}

	err := s.rep.Delete(*deleteObjectRequest.UserID, *deleteObjectRequest.ObjectName)
	if err != nil {
		return processRepositoryError(err)
	}
	return nil
}

func (s *RentObjectServer) updateObject(w http.ResponseWriter, r *http.Request) *appError {
	var updateObjectRequest requests.UpdateObjectRequest

	if err := parseRequest(r.Body, &updateObjectRequest); err != nil {
		return &appError{err, "Error while parsing body", http.StatusUnprocessableEntity}
	}
	err := s.rep.Update(*updateObjectRequest.UserID, *updateObjectRequest.ObjectName, *updateObjectRequest.UpdateInput)

	if err != nil {
		return processRepositoryError(err)
	}
	return nil
}

func (s *RentObjectServer) getObject(w http.ResponseWriter, r *http.Request) *appError {
	query := r.URL.Query()
	if !query.Has(ObjectNameQueryParam) || !query.Has(UserIdQueryParam) {
		return &appError{errors.New("getObject: incorrect query parameters name"), "Incorrect query parameters name", http.StatusUnprocessableEntity}
	}

	userID, errUsr := getUserIdParam(query)
	objectName := getObjectNameParam(query)

	if errUsr != nil {
		return &appError{errors.New("getObject: incorrect query parameters value"), "Incorrect query parameters value", http.StatusUnprocessableEntity}
	}

	object, err := s.rep.GetByName(userID, objectName)

	if err != nil {
		return processRepositoryError(err)
	}

	json.NewEncoder(w).Encode(object)
	return nil
}

func (s *RentObjectServer) getAll(w http.ResponseWriter, r *http.Request) *appError {
	query := r.URL.Query()

	if !query.Has(UserIdQueryParam) {
		return &appError{errors.New("getAll: incorrect query parameters name"), "Incorrect query parameters name", http.StatusUnprocessableEntity}
	}

	userID, errUsr := getUserIdParam(query)

	if errUsr != nil {
		return &appError{errors.New("getAll: incorrect query parameters value"), "Incorrect query parameters value", http.StatusUnprocessableEntity}
	}

	objects, err := s.rep.GetAll(userID)
	if err != nil {
		return processRepositoryError(err)
	}

	json.NewEncoder(w).Encode(objects)
	return nil
}

func (s *RentObjectServer) addRecord(w http.ResponseWriter, r *http.Request) *appError {
	var addRecordRequest requests.AddRecordRequest

	if err := parseRequest(r.Body, &addRecordRequest); err != nil {
		return &appError{err, "Error while parsing body", http.StatusUnprocessableEntity}
	}

	_, err := s.rep.AddRecord(*addRecordRequest.UserID, *addRecordRequest.ObjectName, *addRecordRequest.Record)
	if err != nil {
		return processRepositoryError(err)
	}
	return nil
}

func (s *RentObjectServer) deleteRecord(w http.ResponseWriter, r *http.Request) *appError {
	var deleteRecordRequest requests.DeleteRecordRequest

	if err := parseRequest(r.Body, &deleteRecordRequest); err != nil {
		return &appError{err, "Error while parsing body", http.StatusUnprocessableEntity}
	}

	err := s.rep.DeleteRecord(*deleteRecordRequest.UserID, *deleteRecordRequest.ObjectName, *deleteRecordRequest.RecordIndex)

	if err != nil {
		return processRepositoryError(err)
	}
	return nil
}

func (s *RentObjectServer) updateRecord(w http.ResponseWriter, r *http.Request) *appError {
	var updateRecordRequest requests.UpdateRecordRequest

	if err := parseRequest(r.Body, &updateRecordRequest); err != nil {
		return &appError{err, "Error while parsing body", http.StatusUnprocessableEntity}
	}

	err := s.rep.UpdateRecord(*updateRecordRequest.UserID, *updateRecordRequest.ObjectName, *updateRecordRequest.RecordIndex, *updateRecordRequest.UpdateInput)
	if err != nil {
		return processRepositoryError(err)
	}
	return nil
}

func (s *RentObjectServer) getRecord(w http.ResponseWriter, r *http.Request) *appError {
	query := r.URL.Query()
	if !isQueryHasParameters(query, UserIdQueryParam, ObjectNameQueryParam, RecordIndexQueryParam) {
		return &appError{errors.New("getRecord: incorrect query parameters name"), "Incorrect query parameters name", http.StatusUnprocessableEntity}
	}

	userID, errUsr := getUserIdParam(query)
	objectName := getObjectNameParam(query)
	recordIndex, errRec := getRecordIndexParam(query)

	if errRec != nil || errUsr != nil {
		return &appError{errors.New("getRecord: incorrect query parameters value"), "Incorrect query parameters value", http.StatusUnprocessableEntity}
	}
	record, err := s.rep.GetRecordByIndex(userID, objectName, recordIndex)

	if err != nil {
		return processRepositoryError(err)
	}

	json.NewEncoder(w).Encode(record)
	return nil
}

func (s *RentObjectServer) getRecords(w http.ResponseWriter, r *http.Request) *appError {
	query := r.URL.Query()
	if !query.Has(ObjectNameQueryParam) || !query.Has(UserIdQueryParam) {
		return &appError{errors.New("getRecords: incorrect query parameters name"), "Incorrect query parameters name", http.StatusUnprocessableEntity}
	}

	userID, errUsr := getUserIdParam(query)
	objectName := getObjectNameParam(query)

	if errUsr != nil {
		return &appError{errors.New("getRecords: incorrect query parameters value"), "Incorrect query parameters value", http.StatusUnprocessableEntity}
	}

	records, err := s.rep.GetAllRecords(userID, objectName)
	if err != nil {
		return processRepositoryError(err)
	}

	json.NewEncoder(w).Encode(records)
	return nil
}

func (s *RentObjectServer) getObjectInfo(w http.ResponseWriter, r *http.Request) *appError {
	query := r.URL.Query()
	if !query.Has(ObjectNameQueryParam) || !query.Has(UserIdQueryParam) {
		return &appError{errors.New("getObjectInfo: incorrect query parameters name"), "Incorrect query parameters name", http.StatusUnprocessableEntity}
	}

	userID, errUsr := getUserIdParam(query)
	objectName := getObjectNameParam(query)

	if errUsr != nil {
		return &appError{errors.New("getObjectInfo: incorrect query parameters value"), "Incorrect query parameters value", http.StatusUnprocessableEntity}
	}

	object, err := s.rep.GetByName(userID, objectName)

	if err != nil {
		return processRepositoryError(err)
	}

	json.NewEncoder(w).Encode(domain.NewRentObjectInfo(object))
	return nil
}

func getUserIdParam(query url.Values) (int64, error) {
	return strconv.ParseInt(query.Get(UserIdQueryParam), 10, 64)
}

func getObjectNameParam(query url.Values) string {
	return query.Get(ObjectNameQueryParam)
}

func getRecordIndexParam(query url.Values) (int, error) {
	return strconv.Atoi(query.Get(RecordIndexQueryParam))
}

func isQueryHasParameters(query url.Values, parameters ...string) bool {
	for _, p := range parameters {
		if !query.Has(p) {
			return false
		}
	}
	return true
}

func parseRequest(r io.Reader, req any) error {
	if err := json.NewDecoder(r).Decode(req); err != nil {
		return err
	}

	if err := requests.CheckRequest(req); err != nil {
		return err
	}

	return nil
}

func processRepositoryError(err error) *appError {
	switch err {
	case domain.RecordNotFoundError:
		return &appError{err, "Record not found", http.StatusNotFound}
	case repository.ObjectNotFoundError:
		return &appError{err, "Object not found", http.StatusNotFound}
	case repository.ObjectAlreadyExists:
		return &appError{err, "Object already exists", http.StatusConflict}
	default:
		return &appError{err, "Error happend on server", http.StatusInternalServerError}
	}
}
