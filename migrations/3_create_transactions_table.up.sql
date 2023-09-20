-- transactions definition

CREATE TABLE transactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date INTEGER Not Null,
	ticker_symbol TEXT(20) Not Null REFERENCES `tickers.symbol` ON DELETE Cascade,
	transaction_type TEXT(20) Not Null,
	volume INTEGER Not Null,
	price INTEGER Not Null,
	commission INTEGER Default 0,
	note TEXT,
	porfolio_id INTEGER Not Null REFERENCES `portfolios.id` ON DELETE Cascade,
	ref_id INTEGER REFERENCES `transactions.id` ON DELETE Set NULL 	
);
