package newscrawler

import (
	"github.com/PuerkitoBio/goquery"
)

type News struct {
	Url   string `redis:"url"`
	Title string `redis:"title"`
	Date  string `redis:"date"`
	Kind  string `redis:"content"`
}

func fetchNews() (*[]News, error) {
	var (
		news     News
		newslist []News = make([]News, 0)
		url             = "http://www.studiareinformatica.uniroma1.it/avvisi"
	)

	// Fetch the web page
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	// Find the list of news
	content := doc.Find(".item-list > ul")

	content.Children().Each(func(i int, s *goquery.Selection) {
		// Scrape the title and url
		link := s.Find(".views-field-title > .field-content > a")
		news.Title = link.Text()
		news.Url, _ = link.Attr("href")

		news.Url = url + news.Url

		// Scrape the date
		news.Date, _ = s.Find("span[property='dc:date']").Attr("content")

		// Scrape the tag
		news.Kind, _ = s.Find(".views-field-field-archivio > .field-content > a").Attr("href")

		newslist = append(newslist, news)
	})

	return &newslist, nil
}
