package main

import (
	"testing"
)

func TestGetSymbolParameter(t *testing.T) {
	result := GetSymbolParameter()
	if result != "005930.KS+" {
		t.Fatalf("symbol code is wrong. your result is %s", result)
	}
}

func TestCollectStockLog(t *testing.T) {

}

/*
func TestParseYahooCSV(t *testing.T) {
	result := ParseYahooCSV("http://finance.yahoo.com/d/quotes.csv?s=AAPL+GOOG+MSFT&f=nab")
	for _, stockLog := range result {
		if stockLog.Code != "AAPL" && stockLog.Code != "GOOG" && stockLog.Code != "MSFT" {
			t.Fatalf("code is wrong! your code is %s", stockLog.Code)
		}
	}
}
*/
