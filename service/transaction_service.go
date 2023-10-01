package service

import (
	"context"
	"database/sql"
	"fmt"
	"myfinace/database"
	"myfinace/model"
)

type TransactionService interface {
	List(ctx context.Context, opt database.ListOptions, dest interface{}) error
	Create(ctx context.Context, t model.Transaction) (int64, error)
	Get(ctx context.Context, t model.Transaction, dest interface{}) (err error)
	Update(ctx context.Context, t model.Transaction) (int, error)
	Delete(ctx context.Context, id uint64) (int, error)
}

func NewTransactionService(db database.DB) (service TransactionService) {
	service = &transaction{db: db}
	return
}

type transaction struct {
	db  database.DB
	rec database.Record
}

func (s *transaction) List(ctx context.Context, opt database.ListOptions, dest interface{}) error {
	q, args := opt.BuildQuery()
	err := s.db.Select(ctx, dest, q, args...)
	return err
}

func (s *transaction) Create(ctx context.Context, t model.Transaction) (int64, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	query, args, err := t.GenerateInsertStatement()
	if err != nil {
		return 0, err
	}
	res, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *transaction) Get(ctx context.Context, t model.Transaction, dest interface{}) (err error) {
	q := fmt.Sprintf("Select * From %s Where `%s` = ?;", t.TableName(), t.PrimaryKey())
	err = s.db.Get(ctx, dest, q, t.ID)
	return
}

func (s *transaction) Update(ctx context.Context, t model.Transaction) (int, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	q := fmt.Sprintf(`Update %s
		Set
			date = :date,
			ticker_symbol = :ticker_symbol,
			transaction_type = :transaction_type,
			volume = :volume,
			commission = :commission,
			note = :note,
			portfolio_id = :portfolio_id,
			ref_id = :ref_id
		Where %s = :id;
		`, t.TableName(), t.PrimaryKey())
	res, err := s.db.Exec(ctx, q,
		sql.Named("id", t.ID),
		sql.Named("date", t.Date),
		sql.Named("ticker_symbol", t.TickerSymbol),
		sql.Named("transaction_type", t.TransactionType),
		sql.Named("volume", t.Volume),
		sql.Named("commission", t.Commission),
		sql.Named("note", t.Note),
		sql.Named("portfolio_id", t.PortfolioID),
		sql.Named("ref_id", t.RefID),
	)
	if err != nil {
		return 0, err
	}
	n := database.GetAffectedRows(res)
	return n, err
}

func (s *transaction) Delete(ctx context.Context, id uint64) (n int, err error) {
	q := fmt.Sprintf("Delete From %s Where %s=?", "transactions", "id")
	res, err := s.db.Exec(ctx, q, id)
	if err != nil {
		return 0, err
	}
	n = database.GetAffectedRows(res)
	return
}
