package handler

import "proxy-data-filter/internal/config"

var appConf *config.Config

func InitHandler(cfg *config.Config) {
	appConf = cfg
}
