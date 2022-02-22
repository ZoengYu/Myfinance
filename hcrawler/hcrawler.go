package hcrawler

import (
	"github.com/gocolly/colly/v2"
)

const (
	ustUrl = "https://coinmarketcap.com/currencies/terrausd/"
	terraLunaUrl = "https://coinmarketcap.com/currencies/terra-luna/"
)

type Supply struct {
	UstPrice					string
	UstMarketSupply   string
	UstTotalSupply    string
	LunnaMarketCap		string
	LunnaMarketSupply string
	LunnaTotalSupply  string
}

func GetMarketInfo() Supply {
	ustRes := WebScraper(ustUrl,".statsContainer .statsValue")
	ustPrice := WebScraper(ustUrl,".priceSection .priceValue")
	terraRes := WebScraper(terraLunaUrl,".statsContainer .statsValue")
	terraRes2 := WebScraper(terraLunaUrl,".statsContainer .maxSupplyValue")

	return Supply{
		UstPrice					: ustPrice[0],
		UstMarketSupply   : ustRes[4],
		UstTotalSupply    : ustRes[4],
		LunnaMarketCap    : terraRes[0],
		LunnaMarketSupply : terraRes[4],
		LunnaTotalSupply  : terraRes2[1],
	}
}

func WebScraper(url,queryKey string) []string{
	var marketinfo []string
	collector := colly.NewCollector()

	collector.OnHTML(queryKey, func(element *colly.HTMLElement){
		marketinfo = append(marketinfo,element.Text)
	})

	collector.OnRequest(func(request *colly.Request) {
		// fmt.Println("Visiting", request.URL.String())
	})
	
	collector.Visit(url)

	return marketinfo
}