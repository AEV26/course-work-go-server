package requests_test

import (
	"rental-server/internal/domain"
	"rental-server/internal/server/requests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerRequest(t *testing.T) {
	t.Run("All fields are represented", func(t *testing.T) {
		userId := int64(1)
		object := domain.RentObject{}
		req := requests.AddObjectRequest{
			UserID: &userId,
			Object: &object,
		}

		err := requests.CheckRequest(req)
		assert.NoError(t, err)
	})
	t.Run("Field userId is missing", func(t *testing.T) {
		object := domain.RentObject{}
		req := requests.AddObjectRequest{
			Object: &object,
		}

		err := requests.CheckRequest(req)
		if assert.Error(t, err) {
			assert.Equal(t, err, requests.MissingFieldsError{Missing: []string{"user_id"}})
		}
	})
}
