package main

import (
	"context"
	"proxy-data-filter/internal/app"
	"proxy-data-filter/internal/config"
	"proxy-data-filter/internal/logging"
	"proxy-data-filter/pkg/vld"
)

/**
Делаем отдельным проектом, который будет расположен где угодно
В проекте на PHP пишем генератор конфига по аттрибутам или как-либо еще
этот конфиг сохраняем в репе проекта
при деплое делаем copy конфига в директорию данного проекта и делаем вызов перезапуска конфига без остановки сервера
+ имеем команды для перезапуска сервера, для обновления из репозитория

? что если данный проект на другом физическом сервере. как доставить конфиг?
*/

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	logger := logging.NewLogger(cfg.Env)
	ctx = logging.ContextWithLogger(ctx, logger)

	logging.GetLogger(ctx).Infoln("Starting application")

	vld.InitValidator(ctx)

	a, err := app.NewApp(cfg)
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}

	logging.GetLogger(ctx).Println("Before Run")
	if err = a.Run(ctx); err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}
}
