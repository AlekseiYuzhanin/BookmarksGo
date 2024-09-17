# Bookmark Manager

## Установка

### Предварительные требования

Убедитесь, что у вас установлены следующие компоненты:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Шаги по установке

1. Скопируйте файл конфигурации:
   bash
   mv .env.example .env

2. Соберите и запустите приложение с помощью Docker Compose:
   bash
   docker-compose up --build
## Использование

После успешного запуска приложения вы можете использовать следующие маршруты:

### Добавление закладки

**Метод:** POST  
**URL:** http://localhost:8080/bookmarks/add

**Тело запроса:**
json
{
"content": "Название закладки",
}
### Получение закладок

**Метод:** GET  
**URL:** http://localhost:8080/bookmarks

Этот запрос вернет список всех сохраненных закладок.
