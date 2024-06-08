package requests

import (
	"fmt"
	"reflect"
	"rental-server/internal/domain"
)

type MissingFieldsError struct {
	Missing []string
}

func (e MissingFieldsError) Error() string {
	return fmt.Sprintf("Missing fields in request: %#q", e.Missing)
}

type ServerRequest interface{}

func CheckRequest(r ServerRequest) error {
	v := reflect.ValueOf(r)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	var missing []string

	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).IsNil() {
			field := t.Field(i)
			missing = append(missing, field.Tag.Get("json"))
		}
	}

	if len(missing) != 0 {
		return MissingFieldsError{Missing: missing}
	}

	return nil
}

type AddObjectRequest struct {
	UserID *int64             `json:"user_id"`
	Object *domain.RentObject `json:"object"`
}

type DeleteObjectRequest struct {
	UserID     *int64  `json:"user_id"`
	ObjectName *string `json:"object_name"`
}

type UpdateObjectRequest struct {
	UserID      *int64                        `json:"user_id"`
	ObjectName  *string                       `json:"object_name"`
	UpdateInput *domain.UpdateRentObjectInput `json:"update_input"`
}

type AddRecordRequest struct {
	UserID     *int64         `json:"user_id"`
	ObjectName *string        `json:"object_name"`
	Record     *domain.Record `json:"record"`
}

type DeleteRecordRequest struct {
	UserID      *int64  `json:"user_id"`
	ObjectName  *string `json:"object_name"`
	RecordIndex *int    `json:"record_index"`
}

type UpdateRecordRequest struct {
	UserID      *int64                    `json:"user_id"`
	ObjectName  *string                   `json:"object_name"`
	RecordIndex *int                      `json:"record_index"`
	UpdateInput *domain.UpdateRecordInput `json:"update_input"`
}
