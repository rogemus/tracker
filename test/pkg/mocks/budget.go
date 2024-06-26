package mocks_test

import (
	"math/rand"
	"time"
	"tracker/pkg/model"
)

func GenerateBudget(id ...int) model.Budget {
	mock_time := time.Date(2020, 23, 40, 56, 70, 0, 0, time.UTC)
	mock_id := rand.Intn(9999)

	if id != nil {
		mock_id = id[0]
	}

	return model.Budget{
		ID:          mock_id,
		Uuid:        "mock uuid",
		Created:     mock_time.UTC(),
		Description: "mock description",
		Title:       "mock title",
		UserID:      1,
	}
}
