# tgPriceScraper
Бот для телеграм, который парсит сайт https://books.toscrape.com и скидывает все книги и цены с него

Чтобы запустить бота нужно:

  Клонировать репозиторий
  
  Создать в нём файл .env с токеном вашего бота с BotFather и данными для подключения к базе данных:
  
    BOT_TOKEN=8322751926:AAEWYYTgmJiUBTyqQpVx9a3UriMvLXBEw1w
    DB_USER=user
    DB_PASSWORD=password
    DB_NAME=database
    DB_HOST=db
    DB_PORT=5432
    
  Убедиться что скачан Docker и написать в терминале "docker-compose up"
