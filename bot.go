package main

import (
	"fmt"
	"github.com/SURETLY/GO-SDK"
	"math/rand"
	"os"
	"time"
)

const intSet = "012"

func main() {
	sur := suretly.NewDemo("59d25e8bcea0995959de2da9", "gobot123123123")

	println(curTime(), "Получаем лимиты на заявку...")
	// получили лимиты на заявку
	loan, err := sur.Options()
	if err.Msg != "" {
		os.Exit(1)
	}
	fmt.Println(loan)
	println(curTime(), "Принимаем заявку на «Микрозайм под поручительство» соответствующую лимитам...")
	time.Sleep(2 * time.Second)
	println(curTime(), "Идентифицируем Заемщика...")
	time.Sleep(2 * time.Second)

	// генерим внутренний uid заявки
	println(curTime(), "Генерим внутренний uid заявки...")
	time.Sleep(2 * time.Second)

	// отправляем данные для заявкки, получаем id заявки
	println(curTime(), "Отправляем Suretly данные договора займа...")
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
	fmt.Println(curTime(), "id новой заявки:", order.Id)
	if err.Msg != "" {
		fmt.Println("Ошибка OrderNew", err)
	}
	time.Sleep(1 * time.Second)

	// по id заявки проверяем статус
	fmt.Println(curTime(), "Проверяем статус новой заявки")
	time.Sleep(2 * time.Second)
	status, err := sur.OrderStatus(order.Id)
	fmt.Println(curTime(), "Статус новой заявки:", status)
	if err.Msg != "" {
		fmt.Println("Ошибка OrderStatus", err)
	}

	// и выгружаем договор по данной заявке
	fmt.Println(curTime(), "Получаем договор для заемщика")
	text, err := sur.ContractGet(order.Id)
	if err.Msg != "" {
		fmt.Println("Ошибка ContractGet", err)
		return
	}
	time.Sleep(2 * time.Second)
	println(text)
	time.Sleep(5 * time.Second)

	fmt.Println(curTime(), "Ожидаем подтверждение от заемщика")
	time.Sleep(3 * time.Second)

	// эмулируем случайным образом согласие заемщика
	success := suretly.StringWithCharset(1, intSet) == "2"
	if success {
		println(curTime(), "Заемщик подписал договор")
		err = sur.ContractAccept(order.Id)
		if err.Msg != "" {
			fmt.Println("Ошибка ContractAccept", err)
		}
		println(curTime(), "Идет поиск поручителей...")
	} else {
		println(curTime(), "Отказ заемщика")
		sur.OrderStop(order.Id)
		os.Exit(0)
	}

	// проверяем изменение статуса заявки
	for i := false; i != true; {
		status, err = sur.OrderStatus(order.Id)
		if err.Msg != "" {
			fmt.Println(curTime(), "Ошибка на стороне сервера", err)
			os.Exit(1)
		}
		time.Sleep(3 * time.Second)

		switch status.Status {
		case 2:
			println(curTime(), "Поиск поручителей остановлен заемщиком")
			os.Exit(0)
			break
		case 3:
			println(curTime(), "Заявка остановлена, по истечению времени, сумма не набрана")
			os.Exit(0)
			break
		case 4:
			println(curTime(), "Заявка успешно завершена, сумма набрана")
			i = true
			break
		}
	}

	// эмулируем случайным образом выдачу займа
	if success {
		println(curTime(), "Заявка оплачена и выдана")
		err = sur.OrderIssued(order.Id)
		if err.Msg != "" {
			fmt.Println(curTime(), "Ошибка OrderIssued", err)
		}
		time.Sleep(2 * time.Second)
	} else {
		println(curTime(), "Отказ заемщика")
		sur.OrderStop(order.Id)
		os.Exit(0)
	}
	println(curTime(), "Ожидание возврата займа")
	time.Sleep(5 * time.Second)

	switch suretly.StringWithCharset(1, intSet) {
	case "0":
		err = sur.OrderUnpaid(order.Id)
		println(curTime(), "Займ не выплачен")
		if err.Msg != "" {
			fmt.Println(curTime(), "Ошибка OrderUnpaid", err)
		}
		break
	case "1":
		err = sur.OrderPaid(order.Id)
		println(curTime(), "Займ выплачен полностью")
		if err.Msg != "" {
			fmt.Println(curTime(), "Ошибка OrderPaid", err)
		}
		break
	case "2":
		sum := rand.Float32() * loan.MaxSum / 2
		err = sur.OrderPartialPaid(order.Id, sum)
		fmt.Println(curTime(), "Займ выплачен частично", sum)
		if err.Msg != "" {
			fmt.Println(curTime(), "Ошибка OrderPartialPaid", err)
		}
		break
	}
}

func curTime() string {
	return time.Now().Format("15:04:05")
}
