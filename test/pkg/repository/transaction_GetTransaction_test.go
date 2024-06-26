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
)

func Test_TransactionRepo_GetTransaction(t *testing.T) {
	transaction := mocks.GenerateTransaction(1, 1)
	var empty_transaction model.Transaction

	testCases := []struct {
		name           string
		expected       model.Transaction
		transactionID  int
		expectedErr    error
		expectedSqlErr error
	}{
		{
			name:           "returns row",
			expected:       transaction,
			transactionID:  1,
			expectedErr:    nil,
			expectedSqlErr: nil,
		},
		{
			name:           "returns err if not found",
			expected:       empty_transaction,
			transactionID:  9999,
			expectedErr:    errs.Transaction404Err,
			expectedSqlErr: sql.ErrNoRows,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

			columns := []string{
				"id",
				"uuid",
				"description",
				"amount",
				"created",
				"budgetid",
				"userid",
			}
			expectedRows := sqlmock.NewRows(columns)

			expectedRows.AddRow(
				transaction.ID,
				transaction.Uuid,
				transaction.Description,
				transaction.Amount,
				transaction.Created,
				transaction.BudgetID,
				transaction.UserID,
			)

			mock.
				ExpectQuery("SELECT id, uuid, description, amount, created, budgetid, userid FROM transactions WHERE id = $1").
				WithArgs(test.transactionID).
				WillReturnRows(expectedRows).
				WillReturnError(test.expectedSqlErr)

			defer db.Close()

			repo := repository.NewTransactionRepository(db)
			result, getErr := repo.GetTransaction(test.transactionID)
			sqlErr := mock.ExpectationsWereMet()

			assert.AssertInt(t, result.ID, test.expected.ID)
			assert.AssertStruct(t, result, test.expected)
			assert.AssertError(t, getErr, test.expectedErr)
			assert.AssertError(t, sqlErr, nil)
		})
	}
}
