package entity

import (
	"time"
)

type Fund struct {
	// Title of Fund
	//
	// required: true
	// example: `SPY: SPDR S&P 500 ETF Trust``
	Title string `json:"title" bson:"title"`

	// Inception Date
	//
	// required: true
	// examle: 1993-01-22
	//		   (will be parsed in UTC)
	InceptionDate time.Time `json:"inception_date" bson:"inception_date"`

	// Total Assets Under Management (in thousands)
	//
	// required: true
	// example: 410867000
	TotalAssets int64 `json:"total_assets" bson:"total_assets"`

	// Shares Outstanding
	//
	// required: true
	// example: 914182000
	SharesOutstanding int64 `json:"shares_outstanding" bson:"shares_outstanding"`

	// Expense Ratio
	//
	// required: true
	// example:  0.0945
	ExpenseRatio float64 `json:"expense_ratio" bson:"expense_ratio"`

	// Tracks This Index
	//
	// required: true
	// example:  S&P 500 Index
	TracksThisIndex string `json:"tracks_this_index" bson:"tracks_this_index"`

	// ETF Database Category
	//
	// required: true
	// example:  Large Cap Growth Equities
	ETFDatabaseCategory string `json:"etf_database_category" bson:"etf_database_category"`

	// Issuer
	//
	// required: true
	// example:  State Street
	Issuer string `json:"issuer" bson:"issuer"`

	// Structure
	//
	// required: true
	// example:  UIT
	Structure string `json:"structure" bson:"structure"`

	// Home page link of Fund
	//
	// required: true
	// example:  https://etfdb.com/etf/SPY
	ETFHomePage string `json:"fund_home_page_link" bson:"fund_home_page_link"`
}
