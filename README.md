## Тестовое задание компании Vortex

### Описание работы приложения
Основная задача приложения - синхронизация работы алгоритмов с инфромацией в базе данных. Процессом синхронизации занимается хэндлер, который раз в 5 минут делает запрос к БД, и обновляет информацию об алгоритмах (работает параллельно основной горутине). Для изменния данных в БД есть API:
1. AddClient() - добавление нового пользователя, выполняется через POST-запрос по адресу: http://{HOST}/users. Шаблон json для запроса:
~~~shell
{
    "clientName": "Boba",
    "version": 2,
    "image": "unknown",
    "cpu": "first",
    "memory":"Low",
    "priority": 1.25,
    "needRestart": fasle
}
~~~
2. UpdateClient() - обновление данных о пользователе, выполняется через PUT-запрос по адресу: http://{HOST}/users. Шаблон json для запроса:
~~~shell
{
    "id": 1
    "clientName": "Boba",
    "version": 2,
    "image": "unknown",
    "cpu": "first",
    "memory":"Low",
    "priority": 1.25,
    "needRestart": fasle
}
~~~

3. DeleteClient() - удаление данных о пользователе, выполняется через DELETE-запрос по адресу: http://{HOST}/users. Шаблон json для запроса:

~~~shell
{
    "id": 1
    "clientName": "Boba",
    "version": 2,
    "image": "unknown",
    "cpu": "first",
    "memory":"Low",
    "priority": 1.25,
    "needRestart": fasle
}
~~~

4. UpdateAlgorithmStatus - измение статусов алгоритмов, выполняется через PUT-запрос по адресу: http://{HOST}/algorithmStatus. Шаблон json для запроса:
~~~shell
{
    "id": 1,
    "clientID": 1,
    "VWAP": true,
    "TWAP": false,
    "HFT": false
}
~~~

### Устновка и подготовка к запуску
1. Необходимо склонировать репозиторий командой:
~~~shell
git clone https://github.com/BobaUbisoft17/vortex_test
~~~
2. Необходимо создать файл .env и внести в него следующие значения:
~~~shell
LOGPATH="путь для сохранения логов"
HOST="хост, на котором будет работать приложение"
PORT="порт, который будет прослушивать приложение"
user="имя пользователя базы данных"
password="пароль пользователя"
DBName="название базы данных"
DATABASEURL="путь для подключения к БД"
~~~

### Запуск
Для запуска приложения необходимо ввести команду
~~~shell
docker-compose up --build
~~~