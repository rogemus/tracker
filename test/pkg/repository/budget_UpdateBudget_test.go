package repository_test

import (
	"testing"
	"tracker/pkg/model"
	"tracker/pkg/repository"
	assert "tracker/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_BudgetRepo_UpdateBudget(t *testing.T) {
	testCases := []struct {
		name              string
		expectedErr       error
		expectedSqlErr    error
		budgetId          int
		budgetTitle       string
		budgetDesctiption string
	}{
		{
			name:              "update budget",
			expectedErr:       nil,
			expectedSqlErr:    nil,
			budgetId:          2,
			budgetTitle:       "Test Title",
			budgetDesctiption: "Test Desc",
		},
		{
			name:              "update budget without desc",
			expectedErr:       nil,
			expectedSqlErr:    nil,
			budgetId:          5,
			budgetTitle:       "Test Title",
			budgetDesctiption: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()

			mock.
				ExpectExec("UPDATE budget SET Title=?, Description=? WHERE ID = ?").
				WithArgs(test.budgetTitle, test.budgetDesctiption, test.budgetId).
				WillReturnResult(sqlmock.NewResult(int64(test.budgetId), 1)).
				WillReturnError(test.expectedSqlErr)
			defer db.Close()

			repo := repository.NewBudgetRepository(db)
			newBudget := model.Budget{Title: test.budgetTitle, Description: test.budgetDesctiption, ID: test.budgetId}
			updateErr := repo.UpdateBudget(newBudget, test.budgetId)
			sqlErr := mock.ExpectationsWereMet()

			assert.AssertError(t, sqlErr, test.expectedSqlErr)
			assert.AssertError(t, updateErr, test.expectedErr)
		})
	}
}
