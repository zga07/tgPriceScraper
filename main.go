package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"tgBot/internal/scraper"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gopkg.in/telebot.v3"
)

var userStates = make(map[int64]string)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка подгрузки окружения: ", err)
	}

	db := initDB()
	defer db.Close()
	createTable(db)

	token := os.Getenv("BOT_TOKEN")
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	bot.Handle("/books", func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		userStates[c.Sender().ID] = "/books"
		return c.Send("Какую страницу книг мне открыть? Пришли номер (от 1 до 50):")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		username := c.Sender().Username
		firstName := c.Sender().FirstName

		saveUser(db, userID, username)

		if c.Text() == "/beer" {
			log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
			allUsers := getAllUsers(db)
			notification := fmt.Sprintf("‼️Минуточку внимания, @%s открыл бутылочку хмельного‼️", username)
			if username == "" {
				notification = fmt.Sprintf("‼️Минуточку внимания, %s открыл бутылочку хмельного‼️", firstName)
			}

			for _, id := range allUsers {
				if id == userID {
					continue
				}
				bot.Send(&telebot.User{ID: id}, notification)
			}
			return c.Send("Всем пришло уведомление о твоём намерении выпить пива")
		}

		if userStates[userID] == "/books" {
			log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
			pageNum, err := strconv.Atoi(c.Text())
			if err != nil || pageNum < 1 || pageNum > 50 {
				return c.Send("Пожалуйста, введи число от 1 до 50.")
			}

			delete(userStates, userID)
			c.Send(fmt.Sprintf("Ищу книги на странице %d...", pageNum))

			result := scraper.GetBooks(pageNum)
			return c.Send(result)
		}

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

	bot.Start()
}

func initDB() *sql.DB {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("База не ответила на пинг: ", err)
	}

	fmt.Println("Бот подключен к базе данных!")
	return db
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		tg_id BIGINT UNIQUE NOT NULL,
		username TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Не удалось создать/найти таблицу:", err)
	}
}

func saveUser(db *sql.DB, tgID int64, username string) {
	query := `INSERT INTO users (tg_id, username) VALUES ($1, $2) ON CONFLICT (tg_id) DO NOTHING`
	_, err := db.Exec(query, tgID, username)
	if err != nil {
		log.Println("Ошибка сохранения пользователя:", err)
	}
}

func getAllUsers(db *sql.DB) []int64 {
	rows, err := db.Query("SELECT tg_id FROM users")
	if err != nil {
		log.Println("Ошибка получения списка юзеров:", err)
		return nil
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		} else {
			log.Println("Ошибка добавления пользователя: ", id)
		}
	}
	return ids
}
