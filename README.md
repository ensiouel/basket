## Basket - сервис для хранения и передачи файлов с использованием gRPC 

---

## Стек технологий

#### API
- gin
- gRPC
#### Хранилище:
- PostgreSQL
#### Логирование:
- slog

---

#### Проблемы возникшие при разработке и их решение:

1) ***Проблема:*** Как хранить и отдавать пользовательские файлы?  
   ***Решение:*** Был написан [микросервис](https://github.com/ensiouel/static) с использованием gRPC, который принимает, хранит и отдает файлы.

---

## Примеры использования

### Получить информацию о файле

#### Запрос
```http request
GET http://localhost:8080/api/v1/file/dc109eb3-5b07-4d03-91f1-a70a90970f67
Accept: application/json
```

#### Ответ
```json
{
  "response": {
    "id": "dc109eb3-5b07-4d03-91f1-a70a90970f67",
    "title": "My cat",
    "name": "cat.png",
    "description": "Look at my beautiful FAT cat ^-^",
    "size": 514784,
    "download_count": 12,
    "created_at": "2023-06-30T20:03:00.084235+03:00",
    "updated_at": "2023-07-01T21:48:54.804697+03:00"
  }
}
```

---

### Загрузить файл

#### Запрос
```http request
POST http://localhost:8080/api/v1/file/upload
Accept: application/json
Content-Type: multipart/form-data; boundary=boundary

--boundary
Content-Disposition: form-data; name="file"; filename="doggie.gif"

< ./doggie.gif
```

#### Ответ
```json
{
   "response": {
      "id": "9c446873-a656-4c51-ae8f-ff7dc7963324",
      "title": "doggie",
      "name": "doggie.gif",
      "description": "",
      "size": 257878,
      "download_count": 0,
      "created_at": "2023-07-02T10:13:19.2345379+03:00",
      "updated_at": "2023-07-02T10:13:19.2345379+03:00"
   }
}
```

---

### Скачать файл

#### Запрос
```http request
GET http://localhost:8080/api/v1/file/9c446873-a656-4c51-ae8f-ff7dc7963324/download
Accept: application/octet-stream
```

#### Ответ
![doggie.gif](https://user-images.githubusercontent.com/92665549/175049146-cadc63b4-0f29-4a53-9088-cba9701035fd.gif)

---

### Изменить информацию о файле

#### Запрос
```http request
PATCH http://localhost:8080/api/v1/file/9c446873-a656-4c51-ae8f-ff7dc7963324
Accept: application/json
Content-Type: application/json

{
  "title": "DogePet DogePet DogePet"
}
```

#### Ответ
```json
{
   "response": {
      "id": "9c446873-a656-4c51-ae8f-ff7dc7963324",
      "title": "DogePet DogePet DogePet",
      "name": "doggie.gif",
      "description": "",
      "size": 257878,
      "download_count": 1,
      "created_at": "2023-07-02T10:13:19.2345379+03:00",
      "updated_at": "2023-07-02T10:17:41.9376175+03:00"
   }
}
```

---

### Удалить файл

#### Запрос
```http request
DELETE http://localhost:8080/api/v1/file/04c403a3-1856-4759-87b5-af043cca0e64
Accept: application/json
```

#### Ответ
```json
{
   "response": 1
}
```

---

## Развертывание

**Собрать** приложение

```shell
make build
```

**Поднять** приложение

```shell
make up
```

---

## Конфигурации

### Все параметры загружаются из файта **[.env](.env)**

```dotenv
GIN_MODE=release

LOGGER_LEVEL=debug

SERVER_ADDR=:8080

GRPC_ADDR=:50051

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=file

FILE_SERVICE_MAX_FILE_SIZE=1073741824
```

---

#### TODO:
- [ ] Аутентификация
- [ ] Кеширование файлов
- [ ] Swagger документация
- [ ] Метрики (Prometheus, Grafana)
- [ ] Admin client