-- tickers definition

CREATE TABLE tickers (
	symbol TEXT Not NULL,
	name TEXT,
	CONSTRAINT tickers_PK PRIMARY KEY (symbol)
);
