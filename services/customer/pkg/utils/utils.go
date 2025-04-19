package utils

import (
	"time"
)

func TranslateStatusToRussian(status string) string {
	statusMap := map[string]string{
		"created":                           "Заказ создан",
		"waiting for courier":               "Ищем свободного курьера",
		"canceled because no couriers":      "Заказ отменен. К сожалению, сейчас нет свободных курьеров",
		"preparing":                         "Заказ собирают",
		"order collected":                   "Заказ собран и ожидает курьера",
		"order collected with some changes": "Заказ собран с некоторыми изменениями. Ожидаем курьера",
		"declined because no products left": "Заказ отменен. К сожалению, все заказанные вами продукты раскупили",
		"order taken from shop":             "Курьер забрал заказ и направляется к вам",
		"order delivered":                   "Заказ доставлен",
		"order declined":                    "Заказ отменен",
		"declined by courier":               "Курьер отменил выполнение заказа",
	}
	return statusMap[status]
}

func TruncateTime(inputTime time.Time) string {
	return inputTime.Format("2006-01-02 15:04:05")
}
