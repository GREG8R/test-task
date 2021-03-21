package pg_storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
)

type History struct {
	Amount decimal.Decimal
	Hour   time.Time
}

type PgStorage interface {
	SaveMoney(ctx context.Context, amount decimal.Decimal, date time.Time) error
	GetHistory(ctx context.Context, startDate, endDate time.Time) ([]History, error)
}

type Storage struct {
	Conn *pgxpool.Pool
}

func (s Storage) SaveMoney(ctx context.Context, amount decimal.Decimal, date time.Time) error {
	_, err := s.Conn.Exec(ctx, "INSERT INTO transactions (amount, date, created_at) VALUES ($1, $2, $3);", amount, date, time.Now().UTC())

	hour := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 0, 0, 0, date.Location())
	sql := `UPDATE history
			SET amount = amount + $1,
				updated_at = $2
			WHERE hour = $3;`

	result, err := s.Conn.Exec(ctx, sql, amount, time.Now().UTC(), hour)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		_, err = s.Conn.Exec(ctx, "INSERT INTO history (amount, hour, created_at, updated_at) VALUES ($1, $2, $3, $4);",
			amount, hour, time.Now().UTC(), time.Now().UTC())
	}

	return err
}

func (s Storage) GetHistory(ctx context.Context, startDate, endDate time.Time) ([]History, error) {
	var historyRows []History
	rows, err := s.Conn.Query(ctx, "SELECT amount, hour FROM history WHERE hour >= $1 AND hour <= $2;", startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h History
		rows.Scan(&h.Amount, &h.Hour)
		historyRows = append(historyRows, h)
	}

	return historyRows, nil
}
