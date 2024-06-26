package repository_test

import (
	"testing"
	"tracker/pkg/model"
	"tracker/pkg/repository"
	assert "tracker/pkg/utils"
	mocks "tracker/test/pkg/mocks"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_TransactionRepo_UpdateTransaction(t *testing.T) {
	transaction := mocks.GenerateTransaction(1, 1)

	testCases := []struct {
		name           string
		expectedErr    error
		expectedSqlErr error
		transaction    model.Transaction
		transactionID  int
	}{
		{
			name:           "Update transaction",
			expectedErr:    nil,
			expectedSqlErr: nil,
			transaction:    transaction,
			transactionID:  1,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()

			mock.
				ExpectExec("UPDATE transactions SET amount=$1, description=$2 WHERE id=$3").
				WithArgs(
					test.transaction.Amount,
					test.transaction.Description,
					test.transaction.ID,
				).
				WillReturnResult(sqlmock.NewResult(int64(test.transaction.ID), 1)).
				WillReturnError(test.expectedSqlErr)

			repo := repository.NewTransactionRepository(db)

			newTransactionId, createErr := repo.UpdateTransaction(test.transaction)
			err := mock.ExpectationsWereMet()

			assert.AssertInt(t, int(newTransactionId), test.transaction.ID)
			assert.AssertError(t, err, test.expectedSqlErr)
			assert.AssertError(t, createErr, test.expectedErr)
		})
	}
}
