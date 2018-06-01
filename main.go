// Package classification GO2HAL API.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Marc Arndt<marc@marcarndt.com> http://www.marcarndt.com
//     Title: go2hal API
//
//     basePath: /api
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"github.com/weAutomateEverything/bankHal/bankCallout"
	"github.com/weAutomateEverything/bankHal/bankSkynet"
	"github.com/weAutomateEverything/bankHal/bankldapService"
	go2hal2 "github.com/weAutomateEverything/go2hal/go2hal"
)

func main() {

	go2hal := go2hal2.NewGo2Hal()

	firstcallStore := bankCallout.NewMongoStore(go2hal.DB)
	bankLdapStore := bankldapService.NewMongoStore(go2hal.DB)

	//A auth service for our bank. If you dont want auth use auth.alwaysTrueAuthService()
	go2hal.AuthService = bankldapService.NewService(bankLdapStore)
	go2hal.FirstCallService = bankCallout.NewService(firstcallStore, go2hal.TelegramService, go2hal.TelegramStore)

	skynetService := bankSkynet.NewService(go2hal.AlertService, go2hal.ChefStore, go2hal.CalloutService)

	go2hal.TelegramService.RegisterCommand(bankSkynet.NewRebuildCHefNodeCommand(go2hal.TelegramStore, go2hal.ChefStore, go2hal.TelegramService,
		go2hal.AlertService))
	go2hal.TelegramService.RegisterCommand(bankSkynet.NewRebuildNodeCommand(go2hal.AlertService, skynetService, go2hal.TelegramStore, go2hal.TelegramService))
	go2hal.TelegramService.RegisterCommand(bankldapService.NewRegisterCommand(go2hal.TelegramService, bankLdapStore))
	go2hal.TelegramService.RegisterCommand(bankldapService.NewTokenCommand(go2hal.TelegramService, bankLdapStore))

	go2hal.TelegramService.RegisterCommandLet(bankSkynet.NewRebuildChefNodeEnvironmentReplyCommandlet(go2hal.TelegramService,
		skynetService, go2hal.ChefService, go2hal.TelegramStore))
	go2hal.TelegramService.RegisterCommandLet(bankSkynet.NewRebuildChefNodeExecute(skynetService, go2hal.AlertService, go2hal.TelegramStore, go2hal.TelegramService))
	go2hal.TelegramService.RegisterCommandLet(bankSkynet.NewRebuildChefNodeRecipeReplyCommandlet(go2hal.ChefStore, go2hal.AlertService,
		go2hal.TelegramService))

	go2hal.Mux.Handle("/api/skynet/", bankSkynet.MakeHandler(skynetService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/bankcallout/firstcall", bankCallout.MakeHandler(go2hal.FirstCallService.(bankCallout.Service), go2hal.HTTPLogger, go2hal.MachineLearningService))

	go2hal.Start()
}
