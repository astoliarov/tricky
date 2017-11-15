package infrastructure

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"tricky/lib/usecases"
)

type RuleSerializer struct {
	Key string `json:"key" binding:"required"`
	Url string `json:"url"  binding:"required,url"`
}

type ResponseSerializer struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func makeApiResponse(c *gin.Context, status int, data interface{}, err error) {
	c.JSON(status, ResponseSerializer{
		Status: status,
		Data:   data,
		Error:  err,
	})
}

type RestApi struct {
	router       *gin.Engine
	rulesUseCase *usecases.RulesUseCase

	port string
}

func NewRestApi(useCase *usecases.RulesUseCase, port string, debug bool) *RestApi {
	api := RestApi{}

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	api.router = gin.Default()
	api.rulesUseCase = useCase
	api.port = port

	rulesRoutesGroup := api.router.Group("/api/v1/rules")
	rulesRoutesGroup.GET("/", api.AllRules)
	rulesRoutesGroup.POST("/", api.NewRule)

	return &api
}

func (api *RestApi) Run() error{
	return http.ListenAndServe(api.port, api.router)
}

func (api *RestApi) AllRules(c *gin.Context) {
	var serializedRules []RuleSerializer

	rules, err := api.rulesUseCase.GetAllRules()
	if err != nil {
		makeApiResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	for i := range rules {
		rule := rules[i]
		serializedRules = append(serializedRules, RuleSerializer{Key: rule.Key, Url: rule.Url.String()})
	}

	data := struct {
		Count int              `json:"count"`
		Items []RuleSerializer `json:"items"`
	}{
		Count: len(serializedRules),
		Items: serializedRules,
	}

	makeApiResponse(c, http.StatusOK, data, err)
}

func (api *RestApi) NewRule(c *gin.Context) {
	var serializer RuleSerializer

	err := c.BindJSON(&serializer)
	if err != nil {
		makeApiResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	parsedUrl, err := url.Parse(serializer.Url)
	if err != nil {
		makeApiResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = api.rulesUseCase.AddRule(serializer.Key, parsedUrl)
	if err != nil {
		makeApiResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	makeApiResponse(c, http.StatusCreated, serializer, nil)
	return
}
