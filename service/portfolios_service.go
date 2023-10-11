package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"myfinance/database"
	"myfinance/helper"
	"myfinance/model"
	"strings"
	"time"
)

type PortfolioService interface {
	List(ctx context.Context, dest interface{}) error
	Create(ctx context.Context, t model.Portfolio) (int64, error)
	Get(ctx context.Context, id uint64, dest interface{}) (err error)
	// Update Portfolio to database, return number of affected rows and error
	Update(ctx context.Context, t model.Portfolio) (int, error)
	Delete(ctx context.Context, id uint64) (int, error)
	HoldingSymbols(ctx context.Context, id uint64) (holding_symbols []*HoldingSymbol, err error)
	UpdateHolding(ctx context.Context, portfolio_id uint64) (err error)
}

func NewPortfolioService(db database.DB) PortfolioService {
	holding_svc := NewHoldingService(db)
	ticker_svc := NewTickerService(db)
	transaction_svc := NewTransactionService(db)
	service := &portfolio{db: db, holding: holding_svc, ticker: ticker_svc, transaction: transaction_svc}
	return service
}

type portfolio struct {
	db          database.DB
	holding     HoldingService
	ticker      TickerService
	transaction TransactionService
}

func (s *portfolio) List(ctx context.Context, dest interface{}) error {
	opt := database.ListOptions{
		Table: "portfolios",
	}
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
		var ticker model.Ticker
		err = s.ticker.Get(ctx, holding_symbol.TickerSymbol, &ticker)
		if err != nil {
			log.Printf("Error geting ticker for symbol %s", holding_symbol.TickerSymbol)
			return
		}
		ticker_holding := model.Holding{
			Symbol:      holding_symbol.TickerSymbol,
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
			sql.Named("symbol", holding_symbol.TickerSymbol),
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
			switch transaction.TransactionType {
			case "buy":
				{
					ticker_holding.TotalShares += int64(transaction.Volume)
					total, err := transaction.Total()
					if err != nil {
						return err
					}
					ticker_holding.TotalCost += total
				}
			case "sell":
				{
					ticker_holding.TotalShares -= int64(transaction.Volume)
					total, err := transaction.Total()
					if err != nil {
						return err
					}
					ticker_holding.TotalCost -= total
				}
			default:
				{
					return errors.New(
						fmt.Sprintf("cannot calculate holding from transaction type %s", transaction.TransactionType),
					)
				}

			}
		}
		ticker_holding.AveragePrice = ticker_holding.GetAveragePrice()
		log.Printf("Ticker holding: %v", ticker_holding)
		err = s.holding.Create(ctx, ticker_holding)
		if err != nil {
			log.Print("Error creating ticker holding")
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
