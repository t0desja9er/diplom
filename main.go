package main

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Сборка данных с помощью json с сайта binance.com
type bnResp struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// Создание базы данных для хранения данных для баланса
type wallet map[string]float64

var bd = map[int64]wallet{}

// Основная функция которая выполняет основную логику
func main() {
	// Подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("5194755061:AAERQZ54XnxRwWHSJSmdEsDD8---PVi8gUQ")
	if err != nil {
		log.Panic(err)
	}
	// Логированное сообщение об авторизации пользователя в консоль
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Инициализируем канал, куда будут приходить обновления от API
	u := tgbotapi.NewUpdate(0)

	// Таймаут для поддержки связи с ботом
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Форма для перебора чего-либо
	for update := range updates {
		// При получении сообщения от пользователя
		if update.Message != nil {

			// Разделение сообщения на слова с помошью пробела
			var command = strings.Split(update.Message.Text, " ")

			switch command[0] {
			// Список команд которые используются в работе бота для вывода сообщений
			case "/start", "/help":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Commands:\n"+
					"/help - |Сообщение помощи\n"+
					"/courses - |Курсы\n"+
					"/balance - |Баланс\n"+
					"ADD - |Пополнить\n"+
					"SUB - |Вывести\n"+
					"DEL - |Удалить"))

			// Основной код команды "Курсы"
			case "/courses":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "|Курсы биржи Binance.com\n https://binance.com"))

				// Вывод курсов по отношению к рублям
				msg := ("|Курсы криптовалют по отношению к рюблю\n")
				msg += getData("RUB", "р.")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				time.Sleep(time.Second * 1)

				// Вывод курсов по отношению к доллару
				msg = ("|Курсы криптовалют по отношению к доллару\n")
				msg += getData("USD", "$.")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				time.Sleep(time.Second * 1)

				// Вывод курсов по отношению к биткоину
				msg = ("|Курсы криптовалют по отношению к биткоину\n")
				msg += getData("BTC", "btc.")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

			// Основной код команды "Балансы"
			case "/balance":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "|Ваш баланс:\n"))
				msg := " "
				var sum float64 = 0

				for key, value := range bd[update.Message.Chat.ID] {
					price, _ := getPriceUSD(key)
					sum += value * price
					msg += fmt.Sprintf("%s: %f [%.2f$]\n", key, value, value*price)
				}

				msg += fmt.Sprintf("|Всего на сумму: %.2f$\n", sum)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

			// Основной код команды "Внести"
			case "ADD":
				if len(command) != 3 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "|Ошибка, неверный формат команды\n"+
						"Верный формат команды: ADD NAME amount"))
					break
				}

				amount, err := strconv.ParseFloat(command[2], 64)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					break
				}

				if _, ok := bd[update.Message.Chat.ID]; !ok {
					bd[update.Message.Chat.ID] = wallet{}
				}

				bd[update.Message.Chat.ID][command[1]] += amount
				balanceText := fmt.Sprintf("%f", bd[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))

			// Основной код команды "Вывести"
			case "SUB":
				if len(command) != 3 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "|Ошибка, неверный формат команды\n"+
						"|Верный формат команды: SUB NAME amount"))
					break
				}

				amount, err := strconv.ParseFloat(command[2], 64)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					break
				}

				if _, ok := bd[update.Message.Chat.ID]; !ok {
					bd[update.Message.Chat.ID] = wallet{}
				}

				bd[update.Message.Chat.ID][command[1]] -= amount
				balanceText := fmt.Sprintf("%f", bd[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))

			// Основной код команды "Удалить"
			case "DEL":
				if len(command) != 2 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "|Ошибка, неверный формат команды\n"+
						"|Верный формат команды: SUB NAME"))
					break
				}

				delete(bd[update.Message.Chat.ID], command[1])

			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "|Ошибка, неизвестная команда\n"+
					"|Что бы узнать список команд введите: /help"))
			}
		}
	}
}

// Функция которая выводит курсы криповалют в сообщении
func getData(valuta, va string) (msg string) {
	i := 0
	price, err := getKurs(valuta)

	if err != nil {
		msg = "|Ошибка получения курсов валют"
	} else {
		for k, v := range price {
			msg += fmt.Sprintf("%s - %.8f %s\n", k, v, va)
			i++
			if i == 22 {
				return
			}
		}
	}
	return
}

// Функция которая связывает бота и сайт binance.com с курсами криптовалют
func getKurs(symbol string) (price map[string]float64, err error) {
	price = make(map[string]float64)
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price")

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("|Invalid code")
		return
	}

	var jsonResp []bnResp
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	if err != nil {
		return
	}

	for i := 0; i < len(jsonResp); i++ {
		if strings.HasPrefix(jsonResp[i].Symbol, symbol) || strings.HasSuffix(jsonResp[i].Symbol, symbol) {
			p, err := strconv.ParseFloat(jsonResp[i].Price, 64)

			if err != nil {
				continue
			}

			price[strings.Replace(jsonResp[i].Symbol, symbol, "", -1)] = p

		}
	}
	return
}

// Функция которая связывает бота и сайт binance.com с курсами для каждой валюты
func getPriceUSD(symbol string) (price float64, err error) {
	// Связь с сайтом binance.com для получения курса и записи в базу данных
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)) // <- Для перевода в доллары
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// Проверка того что связь с сайтом присутствует и она стабильна
	if resp.StatusCode != 200 {
		err = errors.New("|Invalid code")
		return
	}

	var jsonResp bnResp

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	if err != nil {
		return
	}

	price, err = strconv.ParseFloat(jsonResp.Price, 64)

	return
}
