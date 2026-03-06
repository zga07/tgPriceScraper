package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"tgBot/internal/postgresDB"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gopkg.in/telebot.v3"
)

var userStates = make(map[int64]string)

func main() {
	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка подгрузки окружения: ", err)
	}

	db := postgresDB.InitDB()
	defer db.Close()
	postgresDB.CreateTable(db)

	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	register := func(c telebot.Context) {
		postgresDB.SaveUser(db, c.Sender().ID, c.Sender().Username)
	}

	bot.Handle("/beer", func(c telebot.Context) error {
		register(c)

		postgresDB.UpdateLeaderboard(db, c.Sender().ID)

		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		allUsers := postgresDB.GetAllUsers(db)
		notification := fmt.Sprintf("‼️Минуточку внимания, @%s открыл бутылочку хмельного‼️", c.Sender().Username)
		if c.Sender().Username == "" {
			notification = fmt.Sprintf("‼️Минуточку внимания, %s открыл бутылочку хмельного‼️", c.Sender().FirstName)
		}

		for _, id := range allUsers {
			if id == c.Sender().ID {
				continue
			}
			bot.Send(&telebot.User{ID: id}, notification)
			time.Sleep(50 * time.Millisecond)
		}
		return c.Send("Всем пришло уведомление о твоём намерении выпить пива")
	})

	bot.Handle("/top", func(c telebot.Context) error {
		register(c)
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		return c.Send(postgresDB.DisplayLeaderboard(db))
	})

	bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		register(c)
		log.Printf("[%s] написал: %s", c.Sender().Username, "стикер")
		return c.Send(c.Message().Sticker)
	})

	bot.Handle("/broadcast", func(c telebot.Context) error {
		if c.Sender().ID != adminID {
			return c.Send("У вас нет прав")
		}

		userStates[c.Sender().ID] = "awaiting_broadcast"
		return c.Send("Режим рассылки активирован. Пришлите текст, который нужно отправить всем пользователям:")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		register(c)

		if userStates[c.Sender().ID] == "awaiting_broadcast" {
			allUsers := postgresDB.GetAllUsers(db)
			message := c.Text()

			count := 0
			for _, id := range allUsers {
				if id == c.Sender().ID {
					continue
				}
				_, err := bot.Send(&telebot.User{ID: id}, message)
				if err == nil {
					count++
				}
				time.Sleep(50 * time.Millisecond)
			}
			delete(userStates, c.Sender().ID)
			return c.Send(fmt.Sprintf("Рассылка завершена! Получили: %d человек", count))
		}

		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		return c.Send(c.Message().Text)
	})

	bot.Start()
}
