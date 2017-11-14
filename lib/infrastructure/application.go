package infrastructure

import (
	"tricky/lib/infrastructure/dao"
	"tricky/lib/usecases"
)

func StartApp() {
	debug := true
	proxyPort := ":3030"
	apiPort := ":5050"
	certPath := "cert/cert"

	inMemoryRulesDAO := dao.InMemoryDAO{}
	rulesUseCase := usecases.NewRulesUseCase(&inMemoryRulesDAO)

	proxy, err := NewProxy(rulesUseCase, proxyPort, certPath, debug)
	if err != nil {
		panic(err)
	}

	restApi := NewRestApi(rulesUseCase, apiPort)

	go func() {
		restApi.Run()
	}()

	proxy.Run()
}
