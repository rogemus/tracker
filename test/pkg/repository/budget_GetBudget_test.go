package repository_test

import (
	"database/sql"
	"testing"
	"tracker/pkg/errs"
	"tracker/pkg/model"
	"tracker/pkg/repository"
	assert "tracker/pkg/utils"
	mocks "tracker/test/pkg/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func Test_BudgetRepo_GetBudget(t *testing.T) {
	var empty_budget model.Budget
	budget := mocks.GenerateBudget()

	testCases := []struct {
		name           string
		expected       model.Budget
		budgetID       int
		expectedErr    error
		expectedSqlErr error
	}{
		{
			name:           "returns rows for budgetID(1)",
			expected:       budget,
			budgetID:       1,
			expectedErr:    nil,
			expectedSqlErr: nil,
		},
		{
			name:           "returns empty row for budgetID(9999)",
			expected:       empty_budget,
			budgetID:       9999,
			expectedErr:    errs.Budget404Err,
			expectedSqlErr: sql.ErrNoRows,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

			columns := []string{
				"id",
				"uuid",
				"created",
				"description",
				"title",
				"userid",
			}
			expectedRows := sqlmock.NewRows(columns)

			if !cmp.Equal(test.expected, empty_budget) {
				bud := test.expected

				expectedRows.AddRow(
					bud.ID,
					bud.Uuid,
					bud.Created,
					bud.Description,
					bud.Title,
					bud.UserID,
				)
			}

			if test.expectedSqlErr != nil {
				mock.
					ExpectQuery("SELECT id, uuid, created, description, title, userid FROM budgets WHERE id = $1").
					WithArgs(test.budgetID).
					WillReturnError(test.expectedSqlErr)
			} else {
				mock.
					ExpectQuery("SELECT id, uuid, created, description, title, userid FROM budgets WHERE id = $1").
					WithArgs(test.budgetID).
					WillReturnRows(expectedRows)
			}

			defer db.Close()

			repo := repository.NewBudgetRepository(db)
			result, err := repo.GetBudget(test.budgetID)

			assert.AssertStruct[model.Budget](t, result, test.expected)
			assert.AssertError(t, err, test.expectedErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
