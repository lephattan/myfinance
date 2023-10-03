-- holdings definition

CREATE TABLE holdings (
	symbol TEXT,
	portfolio_id INTEGER,
	total_shares INTEGER NOT NULL,
	total_cost INTEGER NOT NULL ,
	average_price INTEGER NOT NULL,
	current_value INTEGER,
	updated_at INTEGER,
	CONSTRAINT holdings_FK FOREIGN KEY (symbol) REFERENCES tickers(symbol) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT holdings_FK_1 FOREIGN KEY (portfolio_id) REFERENCES transactions(id) ON DELETE CASCADE ON UPDATE RESTRICT
);
