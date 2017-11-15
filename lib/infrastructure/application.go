package infrastructure

import (
	"tricky/lib/infrastructure/dao"
	"tricky/lib/usecases"

	"github.com/sirupsen/logrus"
)

func StartApp() {
	logger := logrus.New()

	config := NewConfig()
	err := config.FromYAML()
	if err != nil {
		logger.Warnf("Couldn't read file: %s", err)
	}

	inMemoryRulesDAO := dao.InMemoryDAO{}
	rulesUseCase := usecases.NewRulesUseCase(&inMemoryRulesDAO)

	proxy, err := NewProxy(rulesUseCase, config.ProxyPort, config.CertPath, config.Debug)
	if err != nil {
		logger.Fatalf("Proxy init error: %s", err)
	}

	restApi := NewRestApi(rulesUseCase, config.ApiPort, config.Debug)

	go func() {
		logger.Infof("Start RestApi on %s", config.ApiPort)
		err := restApi.Run()
		logger.Fatalf("rest api run error: %s", err)
	}()

	logger.Infof("Start Proxy on %s", config.ProxyPort)
	err = proxy.Run()
	logger.Fatalf("proxy run error: %s", err)
}
