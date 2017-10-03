package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"github.com/SURETLY/GO-SDK"
)

func main() {
	sur := gosdk.NewDemo("59d25e8bcea0995959de2da9", "gobot123123123")
	println(sur.AuthKeyGen())

	println("Получаем лимиты на заявку...")
	// получили лимиты на заявку
	loan, err := sur.Options()
	if err != nil {
		os.Exit(1)
	}
	println("Принимаем заявку на «Микрозайм под поручительство» соответствующую лимитам...")
	time.Sleep(2 * time.Second)
	println("Идентифицируем Заемщика...")
	time.Sleep(2 * time.Second)

	// генерим внутренний uid заявки
	println("Генерим внутренний uid заявки...")
	time.Sleep(2 * time.Second)

	// отправляем данные для заявкки, получаем id заявки
	println("Отправляем Suretly данные договора займа...")
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
	println("id новой заявки:", id.Id, err)
	time.Sleep(1 * time.Second)

	// по id заявки проверяем статус
	fmt.Println("Проверяем статус новой заявки")
	time.Sleep(2 * time.Second)
	status, err := sur.OrderStatus(id.Id)
	fmt.Println("Статус новой заявки:", status, err)

	// и выгружаем договор по данной заявке
	fmt.Println("Получаем договор для заемщика")
	text, _ := sur.ContractGet(id.Id)
	time.Sleep(2 * time.Second)
	println(text)
	time.Sleep(5 * time.Second)

	fmt.Println("Ожидаем подтверждение от заемщика")
	time.Sleep(3 * time.Second)

	// эмулируем случайным образом согласие заемщика
	if rand.Float32() > 0.5 {
		println("Заемщик подписал договор")
		err = sur.ContractAccept(id.Id)
		if err != nil {
			println("Ошибка ContractAccept", err)
		}
		println("Идет поиск поручителей...")
	} else {
		println("Отказ заемщика")
		sur.OrderStop(id.Id)
		os.Exit(0)
	}

	// проверяем изменение статуса заявки
	for i := 0; i != 1; {
		status, err = sur.OrderStatus(id.Id)
		if err != nil {
			println("Ошибка на стороне сервера", err)
			os.Exit(1)
		}
		time.Sleep(3 * time.Second)

		switch status.Status {
		case 2:
			println("Поиск поручителей остановлен заемщиком")
			os.Exit(0)
			break
		case 3:
			println("Заявка остановлена, по истечению времени, сумма не набрана")
			os.Exit(0)
			break
		case 4:
			println("Заявка успешно завершена, сумма набрана")
			i = 1
			break
		}
	}

	// эмулируем случайным образом выдачу займа
	if rand.Float32() > 0.5 {
		println("Заявка оплачена и выдана")
		err = sur.OrderIssued(id.Id)
		if err != nil {
			println("Ошибка OrderIssued", err)
		}
		time.Sleep(2 * time.Second)
	} else {
		println("Отказ заемщика")
		sur.OrderStop(id.Id)
		os.Exit(0)
	}
	println("Ожидание возврата займа")
	time.Sleep(5 * time.Second)

	switch rand.Int31n(3) {
	case 0:
		err = sur.OrderUnpaid(id.Id)
		println("Займ не выплачен")
		if err != nil {
			println("Ошибка OrderUnpaid", err)
		}
		break
	case 1:
		err = sur.OrderPaid(id.Id)
		println("Займ выплачен полностью")
		if err != nil {
			println("Ошибка OrderPaid", err)
		}
		break
	case 2:
		err = sur.OrderPartialPaid(id.Id, rand.Float32()*loan.MaxSum/2)
		println("Займ выплачен частично")
		if err != nil {
			println("Ошибка OrderPartialPaid", err)
		}
		break
	}
}
