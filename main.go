package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/telebot.v3"
)

func main() {
	pref := telebot.Settings{
		Token:  "token",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		return c.Send(c.Text())
	})

	bot.Handle("/help", func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		response := fmt.Sprintf("привет, %s, я повторяю сообщения", c.Sender().FirstName)
		return c.Send(response)
	})

	bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, "стикер")
		return c.Send(c.Message().Sticker)
	})

	bot.Handle("/books", func(c telebot.Context) error {
		bookString := GetBooks(1)

		return c.Send(bookString)
	})

	bot.Start()
}

func GetBooks(page int) string {
	url := fmt.Sprintf("https://books.toscrape.com/catalogue/page-%d.html", page)

	response, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Ошибка подключения к сайту: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Sprintf("Ошибка сайта: %s", err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return fmt.Sprintf("Ошибка при чтении данных: %s", err)
	}
	var result string
	result = fmt.Sprintf("Страница №%d:\n\n", page)
	doc.Find("article.product_pod").Each(func(i int, s *goquery.Selection) {
		fullTitle, _ := s.Find("h3 a").Attr("title")
		price := s.Find(".price_color").Text()

		result += fmt.Sprintf("%d. %s\nСтоимость: %s\n\n", i+1, fullTitle, price)
	})

	return result
}
