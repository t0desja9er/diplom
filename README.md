# Дипломный проект

## Оисание

Проект предазначен для изучения нового языка программирования, технологий разработки.
Проект способен работать с разными биржами криптовалют, для этого необходимо заменить все ссылки API в проекте на API той биржи которая вам необходима.
Изначально проект работает с биржой криптовалют binance.com, так найти API для этой биржи не составляло особого труда.

## Команды которые используются в боте
В ходе разработки было написано несколько команд.

###/start, /help 
Команда /start вызывается только при начале работы бота. И также как и команда /help выводит сообщения с перечнем всех команд которые можно использовать.
Реализовано жто для того что бы в любой момент можно было посмотреть, какие команды можно использовать в этом боте.

---
###/courses
Является основной командой этого бота. 
Данная команда предназначена для вывода сообщени с курсами самых известных криптовалют.
Курсы выводять в трех сообщениях, это сделано для того того, что бы пользватель не путался в валютах.
Для которых написан вывод. Курсы вывдятся по отношению к рублю, доллару и биткойну.

---
###/balance
В проекте реализован некий виртуальный счет, сделано это для того что бы пользователь мог добавить, вычесть или удалить некое количество криптовалюты.
И эта команда созданна для того что бы пользователь мог вывести этот счет на экран, и что бы у него была возможность увидеть на какую сумму у него накопленно валюты 

---
###ADD
Команда ADD создана что бы добавить криктовалюту на виртульный счет.

Что бы вызвать эту команду необходимо просто в поле сообщения написать ADD NAME amount,
где NAME это сокращенное название криптовалюты, например, BTC, ETH и т.д. а amount это количество той валюты, которое вы хотите добавить.

---
###SUB
Команда SUB создана что бы вывести криктовалюту с виртульного счета.

Что бы вызвать эту команду необходимо просто в поле сообщения написать SUB NAME amount,
где NAME это сокращенное название криптовалюты, например, BTC, ETH и т.д. а amount это количество той валюты, которое вы хотите добавить.

---
###DEL
Команда DEL создана что бы удалить криктовалюту с виртульного счета.

Что бы вызвать эту команду необходимо просто в поле сообщения написать DEL NAME,
где NAME это сокращенное название криптовалюты, например, BTC, ETH и т.д.

## Алгоритм работы
После запуска разработанное программа подключается к созданному вами боту с помощью токена, который выдался вам при создании.
Далее нициализируется закрытый канал связи с вашим аккаунтом телеграмм, по которому будет приходить логированное сообщения об авторизации пользователя.
Также каждые 60 секунд происходит невилимое для вас обновление бота с целью поддержки связи с ботом.

## Функции используемые в ходе разработки
В ходе разработки приложения помимо основной функции main, было созданно еще 3 вспомогательные функции:

- getData
- getKurs
- getPriceUSD

---
### Функция main
Основная функция в которой описаны все команды и все действия со сторнними функциями.

Внутри этой функции была использован оператор switch, это замена всем известного оператора if else.
Только switch позволяет более "опрятно" написать код с перебором большого количества каких-либо команд. 

---
### Функции getData и getKurs
Эти функции предназначены для вывода сообщения, которое выводит курсы криптовалют.

Функции getData и getKurs очень тесно связаны. Так как функция getKurs собирает данные о курсах криптовалют с помощью API.
А функция getData занимается выводом этого сообщения и приводит его к приятному глазу ввиду. И разбивает эти курсы на группы.
Где выводятся курсы относительно рубля, доллара и биткоина, так как это одн из самых известных валют мира.

Данные функции запускаются при вызове команды /courses.

---
### Функция getPriceUSD
Данная функция создана для вывода и конвертации валют которые находятся на виртуальном счете.
И также проверяет что связи с биржай присутсвует и она стабильна.

Данная функция вызывается при отправке команды /balance.