Create Trigger update_ticker_price_updated_at Update Of current_price on tickers
BEGIN 
	UPDATE tickers SET price_updated_at = unixepoch('now') Where symbol = old.symbol And new.current_price != old.current_price;
END

