package service

var receivedMap = map[string]string{"ru": "Вы получили {amount} {currency}", "en": "You received {amount} {currency}"}

var transferredMap = map[string]string{"ru": "Bы отправили: {floatAmount} {currency} Адрес получателя: {destination}",
	"en": "Bы отправили: {floatAmount} {currency} Адрес получателя: {destination}",
}
