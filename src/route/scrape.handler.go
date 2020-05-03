package route

import (
	"fmt"
	"log"
	"net"
	"time"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/redisstorage"
	"golang.org/x/net/html"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2/queue"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ALLOWED_MYNIMO = "mynimo.com"
const ALLOWED_SOURCE = ""
var panickRecover = func() {
	if err := recover(); err != nil {
		log.Println("Error: ", err)
	}
}

func createCollector(redis, redispass string) *colly.Collector{
	sc := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains(ALLOWED_MYNIMO),
		colly.MaxDepth(2),
	)
	sc.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5, RandomDelay: 5 * time.Second})
	sc.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	})

	// Extensions
	extensions.RandomUserAgent(sc)
    extensions.Referer(sc)
	
	// Db
	if redis != "" {
		storage := &redisstorage.Storage{
			Address:  redis,
			Password: redispass,
			DB:       0,
			Prefix:   "job01",
		}
		err := sc.SetStorage(storage)
		if err != nil {
			panic(err)
		}
	}
	return sc
}

func myNimoCollector(cl *colly.Collector, c echo.Context, q *queue.Queue) (bool, string) {
	noError := true
	jobLocation := c.QueryParam("location")
	jobType := c.QueryParam("type")

	if (jobLocation != "" && jobType != "") {
		cl.OnXML("//a[@class=\"item\" and @type=\"nextItem\"]", func(x  *colly.XMLElement) {
			newUrl := x.Request.AbsoluteURL(x.Attr("href"))
			if newUrl != "" { 
				q.AddURL(newUrl)
			}	
		})
		cl.OnXML("//*[@id=\"job-browse-card\"]", func(x  *colly.XMLElement) {
			nodes, _ := htmlquery.QueryAll(x.DOM.(*html.Node), "//div[contains(@class,\"job-browse-card-element\")]")
			for _, node := range nodes {
				defer panickRecover()
				link := htmlquery.FindOne(node, "//a[@class=\"jobTitleLink-v2\"]/@href")
				title := htmlquery.FindOne(node, "//a[@class=\"jobTitleLink-v2\"]")
				salary := htmlquery.FindOne(node, "//a[@class=\"jobTitleLink-v2\"]/following-sibling::div")
				company := htmlquery.FindOne(node, "//span[@class=\"company-browse-info\"]")
				jobColor := "orange"
				fmt.Println(htmlquery.InnerText(title))
				fmt.Println(htmlquery.InnerText(salary))
				fmt.Println(htmlquery.InnerText(company))
				fmt.Println(x.Request.AbsoluteURL(htmlquery.InnerText(link)))
				fmt.Println(jobColor)
				fmt.Println("---------------------------")
			}
		})
		cl.OnError(func(r *colly.Response, err error){
			if err != nil {
				noError = false
				log.Println(err)
				
			}
		})
		cl.OnScraped(func(r *colly.Response) {
			log.Println("+Scraping Mynimo Done...+")
		}) 
		q.AddURL(fmt.Sprintf("https://mynimo.com/%s/%s", jobLocation, jobType))
	} else {
		return noError, "Nothing to do."
	}	
	return noError, ""
}

func ScrapeHandler(redisHost, redisPass string) echo.HandlerFunc{
	q, _ := queue.New(
		8, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	return func (c echo.Context) (err error) {
		collector := createCollector(redisHost, redisPass)
		noError, msg := myNimoCollector(collector, c, q)
		go func() {
			q.Run(collector)
		}()
		if noError {
			if msg != "" {return c.JSON(http.StatusOK, msg)
			}else{return c.JSON(http.StatusOK, "Job searching started...")}
		} else {
			return c.JSON(http.StatusInternalServerError, "Failed searching jobs...")
		}
	}
} 
