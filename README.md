## Формулировка

Простейшее веб приложение, предоставляющее пользователю набор операций над сущностью Person.
Приложение реализует API:

* `GET /persons/{personId}` – информация о человеке;
* `GET /persons` – информация по всем людям;
* `POST /persons` – создание новой записи о человеке;
* `PATCH /persons/{personId}` – обновление существующей записи о человеке;
* `DELETE /persons/{personId}` – удаление записи о человеке.

[Описание API](person-service.yaml) в формате OpenAPI.

### Требования

* Запросы / ответы должны быть в формате JSON.
* Если запись по id не найдена, то возвращать HTTP статус 404 Not Found.
* При создании новой записи о человека (метод POST /person) возвращать HTTP статус 201 Created с пустым телом и
  Header `Location: /api/v1/persons/{personId}`, где `personId` – id созданной записи.
* Приложение должно использовать БД для хранения записей.