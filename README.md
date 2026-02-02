# GraphQL Comments System

Система управления постами и комментариями на Go с использованием GraphQL.  
Реализована функциональность, аналогичная комментариям на популярных платформах (Reddit, Хабр).  

- Язык: Go  
- GraphQL: gqlgen  
- Хранилища: PostgreSQL и Redis (кеш)  
- Тесты: Unit-тесты на бизнес-логику  
- Docker + Docker Compose для быстрого запуска  

---

## Функционал

### Посты
- Просмотр списка постов  
- Просмотр поста с комментариями  
- Возможность запретить оставление комментариев на пост  

### Комментарии
- Иерархическая вложенность (дерево комментариев)  
- Ограничение длины текста комментария до 2000 символов    
- Подписка на новые комментарии через GraphQL Subscriptions  

---

## Быстрый старт

1. Клонируем репозиторий:

```bash
git clone https://github.com/limon4ik-black/graphql-comments-system.git
cd graphql-comments-system
```

2. Запускаем Docker:

```bash
make docker-up
```

3. Миграции:

```bash
make migrate
```

4. Запускаем сервер:

```bash
make run
```
5. Открываем GraphQL Playground:
```bash
http://localhost:8080/
```
# GraphQL примеры

Создать пост:

```bash
mutation {
  createPost(title: "Hello", content: "World", author: "Zhenya") {
    id
    title
    content
    author
    commentsAllowed
  }
}
```
Получить все посты:

```bash
query {
  posts {
    id
    title
    author
    commentsAllowed
  }
}
```
Добавить комментарий к посту:

```bash
mutation {
  addComment(postID: "POST_ID", text: "Hello comment") {
    id
    postID
    author
    text
  }
}
```

Подписка на новые комментарии

```bash
subscription {
  commentAdded(postID: "POST_ID") {
    id
    postID
    author
    text
  }
}
```
Чтобы получить комментарий по подписке нужно сделать запрос выше в одной вкладке браузера и не трогать ее, а уже в другой вкладке http://localhost:8080/ отправить комментарий к посту с подпиской.)

## Запуск unit-тестов:
```bash
make test
```

# Структура проекта:
```bash
.
├── cmd/server            # Точка входа сервера
├── internal/service      # Бизнес-логика
├── internal/repository   # Репозитории
├── internal/domain       # Модели и структуры
├── internal/config       # Конфигурация
├── graph                 # GraphQL схема и резолверы
├── migrations            # SQL миграции для PostgreSQL
├── docker-compose.yml    # Docker Compose для зависимостей
└── Makefile              # Утилиты для запуска
```

