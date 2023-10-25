package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"myfinance/database"
	"myfinance/helper"
	"myfinance/model"
	"strings"
	"time"
)

type PortfolioService interface {
	List(ctx context.Context, opt database.ListOptions, dest interface{}) error
	Create(ctx context.Context, t model.Portfolio) (int64, error)
	Get(ctx context.Context, id uint64, dest interface{}) (err error)
	// Update Portfolio to database, return number of affected rows and error
	Update(ctx context.Context, t model.Portfolio) (int, error)
	Delete(ctx context.Context, id uint64) (int, error)
	HoldingSymbols(ctx context.Context, id uint64) (holding_symbols []*HoldingSymbol, err error)
	UpdateHolding(ctx context.Context, portfolio_id uint64) (err error)
	UpdateSymbolHolding(ctx context.Context, portfolio_id uint64, symbol string) error
	ClearSymbolHolding(ctx context.Context, portfolio_id uint64, symbol string) error
	HoldingValue(ctx context.Context, portfolio_id uint64, dest interface{}) error
	HoldingCost(ctx context.Context, portfolio_id uint64, dest interface{}) error
	Count(ctx context.Context, opt database.ListOptions) (count uint64, err error)
	HoldingSummarry(ctx context.Context, portfolio_id uint64) (*model.HoldingSummarry, error)
}

func NewPortfolioService(db database.DB) PortfolioService {
	holding_svc := NewHoldingService(db)
	ticker_svc := NewTickerService(db)
	service := &portfolio{db: db, holding: holding_svc, ticker: ticker_svc}
	return service
}

type portfolio struct {
	db      database.DB
	holding HoldingService
	ticker  TickerService
}

func (s *portfolio) List(ctx context.Context, opt database.ListOptions, dest interface{}) error {
	q, args := opt.BuildQuery()
	rows, err := s.db.Select(ctx, q, args...)
	if err != nil {
		return err
	}
	return helper.ModelListScan(dest, rows)
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

type HoldingSymbol struct {
	TickerSymbol string
}

// Return slice of ticker symbols that portfolio_id is holding
func (s *portfolio) HoldingSymbols(ctx context.Context, id uint64) (holding_symbols []*HoldingSymbol, err error) {
	q := fmt.Sprintf("Select Distinct ticker_symbol from %s Where portfolio_id = ?", TransactionTablename)
	rows, err := s.db.Select(ctx, q, id)
	if err != nil {
		return
	}
	defer rows.Close()
	var dest []*HoldingSymbol
	err = helper.ModelListScan(&dest, rows)
	return dest, err
}

// Clear current holding and re-calculate holding for given portfolio_id
func (s *portfolio) UpdateHolding(ctx context.Context, portfolio_id uint64) (err error) {
	log.Printf("Update holding for portfolio_id: %d", portfolio_id)
	err = s.ClearHolding(ctx, portfolio_id)
	if err != nil {
		return
	}
	holding_symbols, err := s.HoldingSymbols(ctx, portfolio_id)
	if err != nil {
		log.Println("Error getting holding symbols")
		return err
	}
	log.Printf("Holding Symbols: %v", holding_symbols)

	for _, holding_symbol := range holding_symbols {
		err = s.UpdateSymbolHolding(ctx, portfolio_id, holding_symbol.TickerSymbol)
		if err != nil {
			return err
		}
	}
	return
}

// Clear all holding record of given portfolio_id
func (s *portfolio) ClearHolding(ctx context.Context, portfolio_id uint64) error {
	log.Printf("Clear holding for portfolio_id: %d", portfolio_id)
	q := fmt.Sprintf("Delete From %s Where portfolio_id=?", model.HoldingTablename)
	_, err := s.db.Exec(ctx, q, portfolio_id)
	return err
}

func (s *portfolio) UpdateSymbolHolding(ctx context.Context, portfolio_id uint64, symbol string) (err error) {
	err = s.ClearSymbolHolding(ctx, portfolio_id, symbol)
	if err != nil {
		return err
	}
	var ticker model.Ticker
	err = s.ticker.Get(ctx, symbol, &ticker)
	if err != nil {
		log.Printf("Error geting ticker for symbol %s", symbol)
		return
	}
	ticker_holding := model.Holding{
		Symbol:      symbol,
		PortfolioID: portfolio_id,
		UpdatedAt: database.Nullable[int64]{
			Valid:  true,
			Actual: time.Now().Unix(),
		},
	}
	var transactions model.Transactions
	transaction_query := fmt.Sprintf(
		"Select * From %s Where %s = :symbol And %s = :portfolio_id;",
		model.TransactionsTablename,
		"ticker_symbol",
		"portfolio_id",
	)
	rows, select_err := s.db.Select(
		ctx,
		transaction_query,
		sql.Named("symbol", symbol),
		sql.Named("portfolio_id", portfolio_id),
	)
	if select_err != nil {
		return select_err
	}

	err = helper.ModelListScan(&transactions, rows)
	if err != nil {
		log.Printf("Error scaning transaction list")
		return err
	}

	for _, transaction := range transactions {
		ticker_holding.HandleTransaction(transaction)
	}
	ticker_holding.AveragePrice = ticker_holding.GetAveragePrice()
	if ticker.CurrentPrice.Valid {
		ticker_holding.CurrentValue.Actual = ticker_holding.TotalShares * ticker.CurrentPrice.Actual
		ticker_holding.CurrentValue.Valid = true
	}
	log.Printf("Ticker holding: %v", ticker_holding)
	err = s.holding.Create(ctx, ticker_holding)
	if err != nil {
		log.Print("Error creating ticker holding")
		return err
	}
	return nil
}

// Clear hold record of given symbol in portfolio
func (s *portfolio) ClearSymbolHolding(ctx context.Context, portfolio_id uint64, symbol string) error {
	log.Printf("Clear holding of %s for portfolio_id %d", symbol, portfolio_id)
	q := fmt.Sprintf("Delete From %s Where portfolio_id=? And symbol=?;", model.HoldingTablename)
	_, err := s.db.Exec(ctx, q, portfolio_id, symbol)
	return err
}

func (s *portfolio) HoldingValue(ctx context.Context, portfolio_id uint64, dest interface{}) error {
	query := "Select Sum(current_value) From holdings Where portfolio_id = ? And current_value Is Not Null;"
	rows, err := s.db.Select(ctx, query, portfolio_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return rows.Scan(dest)
}

func (s *portfolio) HoldingCost(ctx context.Context, portfolio_id uint64, dest interface{}) error {
	query := "Select Sum(total_cost) From holdings Where portfolio_id = ?;"
	rows, err := s.db.Select(ctx, query, portfolio_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return rows.Scan(dest)
}

func (s *portfolio) Count(ctx context.Context, opt database.ListOptions) (count uint64, err error) {
	q, args := opt.BuildCountQuery()
	rows, err := s.db.Select(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, sql.ErrNoRows
	}
	err = rows.Scan(&count)
	return
}

func (p *portfolio) HoldingSummarry(ctx context.Context, portfolio_id uint64) (*model.HoldingSummarry, error) {
	var holding_cost int64
	if err := p.HoldingCost(ctx, portfolio_id, &holding_cost); err != nil {
		return nil, err
	}

	var holding_value int64
	if err := p.HoldingValue(ctx, portfolio_id, &holding_value); err != nil {
		return nil, err
	}
	return &model.HoldingSummarry{
		TotalCost:  holding_cost,
		TotalValue: holding_value,
	}, nil

}
