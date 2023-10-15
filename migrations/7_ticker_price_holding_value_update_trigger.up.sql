CREATE Trigger update_holding_value_on_ticker_price UPDATE OF current_price On tickers
WHEN 
	new.current_price IS NOT NULL AND new.current_price != old.current_price
BEGIN 
	UPDATE holdings SET current_value = total_shares * new.current_price
	Where symbol = new.symbol;
END;
