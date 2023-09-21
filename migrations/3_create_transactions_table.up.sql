-- transactions definition

CREATE TABLE transactions (
	id 					INTEGER PRIMARY KEY AUTOINCREMENT,
	date	INTEGER Not Null,
	ticker_symbol 		TEXT 	Not Null REFERENCES tickers ON DELETE Cascade,
	portfolio_id 		INTEGER Not Null REFERENCES portfolios ON DELETE Cascade,
	transaction_type 	TEXT 	Not Null,
	volume 				INTEGER Not Null,
	price 				INTEGER Not Null,
	commission 			INTEGER Default 0,
	note 				TEXT,
	ref_id 				INTEGER REFERENCES transactions ON DELETE Set NULL
	);
