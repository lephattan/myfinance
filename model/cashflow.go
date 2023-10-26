package model

import (
	"myfinance/database"
)

type Cashflow struct {
	Days DaysCashflow
}

type DayCashflow struct {
	Date    int64                    `db:"date"`
	Inflow  database.Nullable[int64] `db:"inflow"`
	Outflow database.Nullable[int64] `db:"outflow"`
}

type DaysCashflow []*DayCashflow

func (d *DayCashflow) Net() int64 {
	var inflow, outflow int64
	if d.Inflow.Valid {
		inflow = d.Inflow.Actual
	} else {
		inflow = 0
	}

	if d.Outflow.Valid {
		outflow = d.Outflow.Actual
	} else {
		outflow = 0
	}
	return inflow - outflow
}
