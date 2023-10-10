package service

import (
	"context"
	"fmt"
	"myfinance/database"
	"myfinance/model"
	"strings"
)

type PortfolioService interface {
	List(ctx context.Context, dest interface{}) error
	Create(ctx context.Context, t model.Portfolio) (int64, error)
	Get(ctx context.Context, id uint64, dest interface{}) (err error)
	// Update Portfolio to database, return number of affected rows and error
	Update(ctx context.Context, t model.Portfolio) (int, error)
	Delete(ctx context.Context, id uint64) (int, error)
}

func NewPortfolioService(db database.DB) PortfolioService {
	holding_svc := NewHoldingService(db)
	service := &portfolio{db: db, holding: holding_svc}
	return service
}

type portfolio struct {
	db      database.DB
	holding HoldingService
}

func (s *portfolio) List(ctx context.Context, dest interface{}) error {
	opt := database.ListOptions{
		Table: "portfolios",
	}
	q, args := opt.BuildQuery()
	err := s.db.Select(ctx, dest, q, args...)
	return err
}

func (s *portfolio) Create(ctx context.Context, t model.Portfolio) (int64, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	query := fmt.Sprintf("Insert Into %s (name, description) Values (?,?);", t.TableName())
	res, err := s.db.Exec(ctx, query, strings.TrimSpace(t.Name), t.Description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *portfolio) Get(ctx context.Context, id uint64, dest interface{}) (err error) {
	q := fmt.Sprintf("Select * From %s Where `id` = ?;", "portfolios")
	err = s.db.Get(ctx, dest, q, id)
	return
}

// Update Portfolio to database, return number of affected rows and error
func (s *portfolio) Update(ctx context.Context, t model.Portfolio) (int, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	q := fmt.Sprintf(`Update %s
		Set
			name = ?,
			description = ?
		Where %s = ?;
		`, t.TableName(), t.PrimaryKey())
	res, err := s.db.Exec(ctx, q, strings.TrimSpace(t.Name), strings.TrimSpace(t.Description), t.ID)
	if err != nil {
		return 0, err
	}
	n := database.GetAffectedRows(res)
	return n, err
}

func (s *portfolio) Delete(ctx context.Context, id uint64) (n int, err error) {
	q := fmt.Sprintf("Delete From %s Where %s=?", "portfolios", "id")
	res, err := s.db.Exec(ctx, q, id)
	if err != nil {
		return 0, err
	}
	n = database.GetAffectedRows(res)
	return
}
