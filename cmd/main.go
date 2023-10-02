package main

import (
	"fmt"
	"github.com/BoryslavGlov/logrusx"
	"github.com/subosito/gotenv"
	"log"
	"os"
	"tonListener/internal/config"
	streamClient "tonListener/internal/stream-client"
	tonClient "tonListener/internal/ton-client"
	"tonListener/pkg/repository"
	client "tonListener/pkg/telegram"
	"tonListener/service"
)

func main() {
	if err := gotenv.Load(); err != nil {
		log.Println(err)
		return
	}

	logx, err := logrusx.New("TonTransactions")
	if err != nil {
		log.Fatal(err)
	}

	logx.Info("TON SERVICE HAS STARTED")

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	api := tonClient.NewApi(conf, logx)
	stream := streamClient.NewClient(conf)
	telegram := client.NewTelegramApi(conf.TelegramToken)
	db, err := repository.NewDB(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))

	sql := repository.NewSQL(db)

	tonListener := service.NewTonListener(api, stream, telegram, sql, logx, conf.MainAddress, conf.NotificationDest)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from panic:\n", err)
		}
	}()

	tonListener.Start()
	select {}

}
