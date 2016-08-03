package newscrawler

import (
	"github.com/PuerkitoBio/goquery"
)

type News struct {
	URL   string `redis:"url"`
	Title string `redis:"title"`
	Date  string `redis:"date"`
	Kind  string `redis:"kind"`
}

func fetchNews() (*[]News, error) {
	var (
		news     News
		newslist []News = make([]News, 0)
		URL             = "http://www.studiareinformatica.uniroma1.it/avvisi"
		baseURL         = "http://www.studiareinformatica.uniroma1.it/"
	)

	// Fetch the web page
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return nil, err
	}

	// Find the list of news
	content := doc.Find(".item-list > ul")

	content.Children().Each(func(i int, s *goquery.Selection) {
		// Scrape the title and URL
		link := s.Find(".views-field-title > .field-content > a")
		news.Title = link.Text()
		news.URL, _ = link.Attr("href")

		news.URL = baseURL + news.URL

		// Scrape the date
		news.Date, _ = s.Find("span[property='dc:date']").Attr("content")

		// Scrape the tag
		news.Kind, _ = s.Find(".views-field-field-archivio > .field-content > a").Attr("href")

		newslist = append(newslist, news)
	})

	return &newslist, nil
}
