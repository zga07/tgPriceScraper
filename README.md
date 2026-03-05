# beerNotifierTgBot
Бот для телеграм, который уведомляет всех своих пользователей об открытии бутылочки хмельного

Чтобы запустить бота нужно:

  Клонировать репозиторий
  
  Создать в нём файл .env с токеном вашего бота с BotFather и данными для подключения к базе данных:
  
    BOT_TOKEN=token
    DB_USER=user
    DB_PASSWORD=password
    DB_NAME=database
    DB_HOST=db
    DB_PORT=5432
    
  Убедиться что скачан Docker и написать в терминале "docker-compose up"
