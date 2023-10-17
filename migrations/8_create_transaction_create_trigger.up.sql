Create Trigger transaction_create_update_holding AFTER INSERT On transactions
BEGIN
	INSERT OR IGNORE INTO holdings (symbol, portfolio_id, total_shares ,total_cost , average_price ,current_value, updated_at) 
		Values (NEW.ticker_symbol, NEW.portfolio_id, 0, 0, 0, null, unixepoch('now'));
	UPDATE holdings 
	SET 
		total_cost = CASE 
			WHEN NEW.transaction_type = 'buy' THEN total_cost + NEW.volume * NEW.price + NEW.commission
			WHEN NEW.transaction_type = 'sell' THEN total_cost - NEW.volume * average_price 
			ELSE total_cost END,
		total_shares = CASE 
			WHEN NEW.transaction_type = 'buy' Then total_shares + NEW.volume
			WHEN NEW.transaction_type = 'sell' Then total_shares - NEW.volume
			ELSE total_shares END,
		average_price = CASE 
			WHEN NEW.transaction_type = 'buy' THEN (total_cost + NEW.volume * NEW.price + NEW.commission) / (total_shares + NEW.volume) 			
			ELSE average_price END,
		current_value = CASE 
			WHEN tickers.current_price IS NOT NULL THEN 
				CASE 
					WHEN NEW.transaction_type = 'buy' Then (total_shares + NEW.volume) * tickers.current_price 
					WHEN NEW.transaction_type = 'sell' Then (total_shares - NEW.volume) * tickers.current_price 
					ELSE total_shares * tickers.current_price 
				END
			ELSE null END,				
		updated_at = unixepoch('now')
	FROM tickers
	Where holdings.symbol = NEW.ticker_symbol And holdings.portfolio_id = NEW.portfolio_id AND tickers.symbol = NEW.ticker_symbol;
END;
