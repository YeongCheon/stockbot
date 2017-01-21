package bot

import (
	"testing"
)

var stockCrawler StockCrawler

func init() {
	stockCrawler = StockCrawler{}
}

func Test_CollectStockData(t *testing.T) {
	collectChannel := make(chan string)
	stockCrawler.CollectStockData(collectChannel)
}
