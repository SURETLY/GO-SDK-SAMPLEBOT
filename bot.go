package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"github.com/SURETLY/GO-SDK"
)

const intSet = "012"

func main() {
	sur := gosdk.NewDemo("59d25e8bcea0995959de2da9", "gobot123123123")

	println(time.Now().Format("15:04:05"), "Получаем лимиты на заявку...")
	// получили лимиты на заявку
	loan, err := sur.Options()
	if err.Msg != "" {
		os.Exit(1)
	}
	fmt.Println(loan)
	println(time.Now().Format("15:04:05"), "Принимаем заявку на «Микрозайм под поручительство» соответствующую лимитам...")
	time.Sleep(2 * time.Second)
	println(time.Now().Format("15:04:05"), "Идентифицируем Заемщика...")
	time.Sleep(2 * time.Second)

	// генерим внутренний uid заявки
	println(time.Now().Format("15:04:05"), "Генерим внутренний uid заявки...")
	time.Sleep(2 * time.Second)

	// отправляем данные для заявкки, получаем id заявки
	println(time.Now().Format("15:04:05"), "Отправляем Suretly данные договора займа...")
	uid := gosdk.StringWithCharset(16, gosdk.Charset)
	newOrder := gosdk.OrderNew{
		Uid:    uid,
		Public: true,
		Borrower: gosdk.Borrower{
			Name: gosdk.Name{
				First:  "Антон",
				Middle: "Викторович",
				Last:   "Фролов",
			},
			Gender: "1",
			Birth: gosdk.Birth{
				Date:  623308357,
				Place: "г.Новосибирск",
			},
			Email:      "frolov_11123@mail.ru",
			Phone:      "+79231232766",
			Ip:         "109.226.15.42",
			ProfileUrl: "https://vk.com/frol_nsk",
			PhotoUrl:   "https://pp.userapi.com/c622420/v622420795/5368/BWdcNhJqFkc.jpg",
			Passport: gosdk.Passport{
				Series:     "4431",
				Number:     "989922",
				IssueDate:  "25.07.2007",
				IssuePlace: "Советский, отдел полиции №10, Управление МВД России по г. Новосибирску",
				IssueCode:  "554-223",
			},
			Registration: gosdk.Address{
				Country:  "Россия",
				Zip:      "630063",
				Area:     "Новосибирская область",
				City:     "Новосибирск",
				Street:   "Труженников",
				House:    "22",
				Building: "",
				Flat:     "24",
			},
			Residential: gosdk.Address{
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
		LoanSum:         loan.MaxSum / 2,
		LoanTerm:        loan.MaxTerm / 2,
		LoanRate:        38.1,
		CurrencyCode:    "RUB",
		Callback:        "callback",
	}
	id, err := sur.OrderNew(newOrder)
	fmt.Println(time.Now().Format("15:04:05"), "id новой заявки:", id.Id)
	if err.Msg != "" {
		fmt.Println("Ошибка OrderNew", err)
	}
	time.Sleep(1 * time.Second)

	// по id заявки проверяем статус
	fmt.Println(time.Now().Format("15:04:05"), "Проверяем статус новой заявки")
	time.Sleep(2 * time.Second)
	status, err := sur.OrderStatus(id.Id)
	fmt.Println(time.Now().Format("15:04:05"), "Статус новой заявки:", status)
	if err.Msg != "" {
		fmt.Println("Ошибка OrderStatus", err)
	}

	// и выгружаем договор по данной заявке
	fmt.Println(time.Now().Format("15:04:05"), "Получаем договор для заемщика")
	text, err := sur.ContractGet(id.Id)
	if err.Msg != "" {
		fmt.Println("Ошибка ContractGet", err)
	}
	time.Sleep(2 * time.Second)
	println(text)
	time.Sleep(5 * time.Second)

	fmt.Println(time.Now().Format("15:04:05"), "Ожидаем подтверждение от заемщика")
	time.Sleep(3 * time.Second)

	// эмулируем случайным образом согласие заемщика
	success := gosdk.StringWithCharset(1, intSet) == "2"
	if success {
		println(time.Now().Format("15:04:05"), "Заемщик подписал договор")
		err = sur.ContractAccept(id.Id)
		if err.Msg != "" {
			fmt.Println("Ошибка ContractAccept", err)
		}
		println(time.Now().Format("15:04:05"), "Идет поиск поручителей...")
	} else {
		println(time.Now().Format("15:04:05"), "Отказ заемщика")
		sur.OrderStop(id.Id)
		os.Exit(0)
	}

	// проверяем изменение статуса заявки
	for i := false; i != true; {
		status, err = sur.OrderStatus(id.Id)
		if err.Msg != "" {
			fmt.Println(time.Now().Format("15:04:05"), "Ошибка на стороне сервера", err)
			os.Exit(1)
		}
		time.Sleep(3 * time.Second)

		switch status.Status {
		case 2:
			println(time.Now().Format("15:04:05"), "Поиск поручителей остановлен заемщиком")
			os.Exit(0)
			break
		case 3:
			println(time.Now().Format("15:04:05"), "Заявка остановлена, по истечению времени, сумма не набрана")
			os.Exit(0)
			break
		case 4:
			println(time.Now().Format("15:04:05"), "Заявка успешно завершена, сумма набрана")
			i = true
			break
		}
	}

	// эмулируем случайным образом выдачу займа
	if success {
		println(time.Now().Format("15:04:05"), "Заявка оплачена и выдана")
		err = sur.OrderIssued(id.Id)
		if err.Msg != "" {
			fmt.Println(time.Now().Format("15:04:05"), "Ошибка OrderIssued", err)
		}
		time.Sleep(2 * time.Second)
	} else {
		println(time.Now().Format("15:04:05"), "Отказ заемщика")
		sur.OrderStop(id.Id)
		os.Exit(0)
	}
	println(time.Now().Format("15:04:05"), "Ожидание возврата займа")
	time.Sleep(5 * time.Second)

	switch gosdk.StringWithCharset(1, intSet) {
	case "0":
		err = sur.OrderUnpaid(id.Id)
		println(time.Now().Format("15:04:05"), "Займ не выплачен")
		if err.Msg != "" {
			fmt.Println(time.Now().Format("15:04:05"), "Ошибка OrderUnpaid", err)
		}
		break
	case "1":
		err = sur.OrderPaid(id.Id)
		println(time.Now().Format("15:04:05"), "Займ выплачен полностью")
		if err.Msg != "" {
			fmt.Println(time.Now().Format("15:04:05"), "Ошибка OrderPaid", err)
		}
		break
	case "2":
		sum := rand.Float32() * loan.MaxSum / 2
		err = sur.OrderPartialPaid(id.Id, sum)
		fmt.Println(time.Now().Format("15:04:05"), "Займ выплачен частично", sum)
		if err.Msg != "" {
			fmt.Println(time.Now().Format("15:04:05"), "Ошибка OrderPartialPaid", err)
		}
		break
	}
}
