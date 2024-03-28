package models

import (
	"database/sql"
	"fmt"
	"time"
	"tracker/pkg/utils"
)

type Budget struct {
	ID          int
	Uuid        string
	Created     time.Time
	Description string
	Title       string
}

func (db *Database) GetBudget(id int) (Budget, error) {
	var b Budget
	query := "SELECT ID, Uuid, Title, Created, Description FROM budget WHERE ID = ?"
	utils.LogInfo(fmt.Sprintf("Looking for Budget(%b)...", id))
	row := db.QueryRow(query, id)

	if err := row.Scan(&b.ID, &b.Uuid, &b.Title, &b.Created, &b.Description); err != nil {
		if err == sql.ErrNoRows {
			return b, fmt.Errorf("GetBudget %d: no such budget", id)
		}

		return b, fmt.Errorf("GetBudget %d: %v", id, err)
	}

	utils.LogInfo(fmt.Sprintf("Found: Budget(%d, %s)", b.ID, b.Title))
	return b, nil
}

func (db *Database) GetBudgets() ([]Budget, error) {
	var budgets []Budget
	query := "SELECT ID, Uuid, Title, Created, Description FROM budget"
	utils.LogInfo("Looking for []Budget...")
	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("GetBudgets: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var b Budget

		if err := rows.Scan(&b.ID, &b.Uuid, &b.Title, &b.Created, &b.Description); err != nil {
			return nil, fmt.Errorf("GetBudgets: %v", err)
		}

		budgets = append(budgets, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBudgets: %v", err)
	}

	utils.LogInfo("Found: []Budget")
	return budgets, nil
}
