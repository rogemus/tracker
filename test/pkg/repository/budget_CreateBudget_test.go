package repository_test

import (
	"testing"
	"tracker/pkg/model"
	"tracker/pkg/repository"
	assert "tracker/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_BudgetRepo_CreateBudget(t *testing.T) {
	testCases := []struct {
		name              string
		expectedErr       error
		expectedSqlErr    error
		budgetId          int
		budgetTitle       string
		budgetDesctiption string
		budgetUserID      int
	}{
		{
			name:              "create budget",
			expectedErr:       nil,
			expectedSqlErr:    nil,
			budgetId:          2,
			budgetTitle:       "Test Title",
			budgetDesctiption: "Test Desc",
			budgetUserID:      1,
		},
		{
			name:              "create budget without desc",
			expectedErr:       nil,
			expectedSqlErr:    nil,
			budgetId:          5,
			budgetTitle:       "Test Title",
			budgetDesctiption: "",
			budgetUserID:      1,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()

			mock.
				ExpectExec("INSERT INTO budget (Title, Description, UserID) VALUES (?, ?, ?)").
				WithArgs(
					test.budgetTitle,
					test.budgetDesctiption,
					test.budgetUserID,
				).
				WillReturnResult(sqlmock.NewResult(int64(test.budgetId), 1)).
				WillReturnError(test.expectedSqlErr)
			defer db.Close()

			repo := repository.NewBudgetRepository(db)

			newBudget := model.Budget{
				Title:       test.budgetTitle,
				Description: test.budgetDesctiption,
				ID:          test.budgetId,
				UserID:      test.budgetUserID,
			}
			newBudgetId, createErr := repo.CreateBudget(newBudget)
			sqlErr := mock.ExpectationsWereMet()

			assert.AssertError(t, createErr, test.expectedErr)
			assert.AssertError(t, sqlErr, test.expectedSqlErr)
			assert.AssertInt(t, int(newBudgetId), test.budgetId)
		})
	}
}
