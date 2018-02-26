package main

import (
	"github.com/SURETLY/GO-SDK"
	"math/rand"
	"os"
	"time"
	. "log"
)

const intSet = "012"

func main() {
	sur := suretly.NewDemo("59d25e8bcea0995959de2da9", "gobot123123123")

	Print("Получаем лимиты на заявку...")
	// получили лимиты на заявку
	loan, err := sur.Options()
	if err.Msg != "" {
		os.Exit(1)
	}
	Print(loan)
	Print("Принимаем заявку на «Микрозайм под поручительство» соответствующую лимитам...")
	sleep(2)
	Print("Идентифицируем Заемщика...")
	sleep(2)

	// генерим внутренний uid заявки
	Print("Генерим внутренний uid заявки...")
	sleep(2)

	// отправляем данные для заявкки, получаем id заявки
	Print("Отправляем Suretly данные договора займа...")
	uid := suretly.StringWithCharset(16, suretly.Charset)
	newOrder := suretly.OrderNew{
		Uid:    uid,
		Public: true,
		Borrower: suretly.Borrower{
			Name: suretly.Name{
				First:  "Антон",
				Middle: "Викторович",
				Last:   "Фролов",
			},
			Gender: "1",
			Birth: suretly.Birth{
				Date:  623308357,
				Place: "г.Новосибирск",
			},
			Email:      "frolov_11123@mail.ru",
			Phone:      "+79231232766",
			Ip:         "109.226.15.42",
			ProfileUrl: "https://vk.com/frol_nsk",
			PhotoUrl:   "https://pp.userapi.com/c622420/v622420795/5368/BWdcNhJqFkc.jpg",
			Passport: suretly.Passport{
				Series:     "4431",
				Number:     "989922",
				IssueDate:  "25.07.2007",
				IssuePlace: "Советский, отдел полиции №10, Управление МВД России по г. Новосибирску",
				IssueCode:  "554-223",
			},
			Registration: suretly.Address{
				Country:  "Россия",
				Zip:      "630063",
				Area:     "Новосибирская область",
				City:     "Новосибирск",
				Street:   "Труженников",
				House:    "22",
				Building: "",
				Flat:     "24",
			},
			Residential: suretly.Address{
				Country:  "Россия",
				Zip:      "630063",
				Area:     "Новосибирская область",
				City:     "Новосибирск",
				Street:   "Труженников",
				House:    "22",
				Building: "",
				Flat:     "24",
			},
		},
		UserCreditScore: 678,
		LoanSum:         rand.Float32() * loan.MaxSum / 2,
		LoanTerm:        rand.Intn(loan.MaxTerm) / 2,
		LoanRate:        38.1,
		CurrencyCode:    "RUB",
		Callback:        "https://anyurl.com/callback",
	}
	order, err := sur.OrderNew(newOrder)
	Print("id новой заявки:", order.Id)
	if err.Msg != "" {
		Print("Ошибка OrderNew", err)
		return
	}
	sleep(1)

	// по id заявки проверяем статус
	Print("Проверяем статус новой заявки")
	sleep(2)
	status, err := sur.OrderStatus(order.Id)
	Print("Статус новой заявки:", status)
	if err.Msg != "" {
		Print("Ошибка OrderStatus", err)
	}

	// и выгружаем договор по данной заявке
	Print("Получаем договор для заемщика")
	text, err := sur.ContractGet(order.Id)
	if err.Msg != "" {
		Print("Ошибка ContractGet", err)
		return
	}
	sleep(2)
	println(text)
	sleep(5)

	Print("Ожидаем подтверждение от заемщика")
	sleep(3)

	// эмулируем случайным образом согласие заемщика
	success := rand.Intn(100) > 30
	if success {
		Print("Заемщик подписал договор")
		err = sur.ContractAccept(order.Id)
		if err.Msg != "" {
			Print("Ошибка ContractAccept", err)
		}
		Print("Идет поиск поручителей...")
	} else {
		Print("Отказ заемщика")
		sur.OrderStop(order.Id)
		os.Exit(0)
	}

	// проверяем изменение статуса заявки
	for i := false; i != true; {
		status, err = sur.OrderStatus(order.Id)
		if err.Msg != "" {
			Print("Ошибка на стороне сервера", err)
			os.Exit(1)
		}
		sleep(3)

		switch status.Status {
		case suretly.ORDER_STATUS_CANCELED:
			Print("Поиск поручителей остановлен заемщиком")
			os.Exit(0)
			break
		case suretly.ORDER_STATUS_TIMEOUT:
			Print("Заявка остановлена, по истечению времени, сумма не набрана")
			os.Exit(0)
			break
		case suretly.ORDER_STATUS_DONE:
			Print("Заявка успешно завершена, сумма набрана")
			i = true
			break
		}
	}

	// эмулируем случайным образом выдачу займа
	if success {
		Print("Заявка оплачена и выдана")
		err = sur.OrderIssued(order.Id)
		if err.Msg != "" {
			Print("Ошибка OrderIssued", err)
		}
		sleep(2)
	} else {
		Print("Отказ заемщика")
		sur.OrderStop(order.Id)
		os.Exit(0)
	}
	Print("Ожидание возврата займа")
	sleep(5)

	switch rand.Intn(2) {
	case 0:
		err = sur.OrderUnpaid(order.Id)
		Print("Займ не выплачен")
		if err.Msg != "" {
			Print("Ошибка OrderUnpaid", err)
		}
		break
	case 1:
		err = sur.OrderPaid(order.Id)
		Print("Займ выплачен полностью")
		if err.Msg != "" {
			Print("Ошибка OrderPaid", err)
		}
		break
	case 2:
		sum := rand.Float32() * loan.MaxSum / 2
		err = sur.OrderPartialPaid(order.Id, sum)
		Print("Займ выплачен частично", sum)
		if err.Msg != "" {
			Print("Ошибка OrderPartialPaid", err)
		}
		break
	}
}

func sleep(t int) {
	time.Sleep(time.Duration(t) * time.Second)
}
