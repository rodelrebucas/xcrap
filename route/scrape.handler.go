package route

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/gocolly/colly/v2"
)


type Task struct{
	Name string `json:name xml:name`
}

func ScrapeHandler(c echo.Context) (err error) {
	sc := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("mynimo.com"),
	)
	sc.OnXML("/html/body/div[3]/div[1]/div/div[8]/div[2]/div/div[3]", func(x  *colly.XMLElement) {
		// fmt.Println(x.Text, x.Name)
		links := x.ChildTexts("//@href")

		for _, link := range links {
			fmt.Println(link)
		}
	})
	sc.Visit("https://mynimo.com/cebu/it-jobs")

	return c.JSON(http.StatusOK, &Task{
		Name: "Scrape one",
	})
}