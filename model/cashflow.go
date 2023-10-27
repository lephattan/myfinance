package model

import (
	"myfinance/database"
	"myfinance/helper"
)

type Cashflow struct {
	Days DaysCashflow
}

func (c *Cashflow) ChartLabels() []string {
	var labels []string
	for _, v := range c.Days {
		labels = append(labels, helper.UnixTimeFmt(v.Date, "2006-02-01"))
	}
	return labels
}

type CashflowChartDataset struct {
	Label string                     `json:"label"`
	Data  []database.Nullable[int64] `json:"data"`
}

type CashflowChartDatasets []*CashflowChartDataset

func (c *Cashflow) Datasets() (datasets *CashflowChartDatasets) {
	var inflow, outflow CashflowChartDataset
	for i, v := range c.Days {
		inflow.Data = append(inflow.Data, v.Inflow)
		outflow.Data = append(outflow.Data, v.Outflow)
		// Outflow is preresented as negative
		outflow.Data[i].Actual = -outflow.Data[i].Actual
	}
	inflow.Label = "inflow"
	outflow.Label = "outflow"
	datasets = &CashflowChartDatasets{&inflow, &outflow}
	return
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
