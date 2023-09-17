-- tickers definition

CREATE TABLE tickers (
	symbol TEXT(20),
	name TEXT,
	CONSTRAINT tickers_PK PRIMARY KEY (symbol)
);
