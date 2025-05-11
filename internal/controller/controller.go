package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"proxy-data-filter/internal/config"
	"proxy-data-filter/internal/handler"
	"proxy-data-filter/pkg/vld"
)

type Init struct {
	cfg    *config.Config
	router *httprouter.Router
}

func New(cfg *config.Config, router *httprouter.Router) *Init {
	return &Init{
		cfg:    cfg,
		router: router,
	}
}

var conf = `[
{"path":"/online/staff/one","method":"POST","existsBody":true,"bodyType":"json","body":{"name":["string","required","max=200"]}},
{"path":"/online/companies/jQu","method":"POST","existsBody":true,"bodyType":"json","body":{"name":["string","required","max=200"]}}
]`

type Conf struct {
	Path       string              `json:"path"`
	Method     string              `json:"method"`
	ExistsBody bool                `json:"existsBody"`
	BodyType   string              `json:"bodyType"`
	Body       map[string][]string `json:"body"`
}

//var conf = []map[string]interface{}{
//	{
//		"path":       "/online/companies/jQu",
//		"method":     "POST",
//		"existsBody": true,
//		"bodyType":   "json",
//		"body":       map[string]interface{}{"name": "string,required,max=200"},
//	},
//}

var configs []Conf

func (controller *Init) SetRoutes(ctx context.Context) error {
	handler.InitHandler(controller.cfg)

	err := json.Unmarshal([]byte(conf), &configs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", configs)

	for id, item := range configs {
		controller.router.HandlerFunc(item.Method, item.Path, BuildHandler(Handle, id))
	}

	return nil
}

func BuildHandler(h http.HandlerFunc, ruleIDX int) http.HandlerFunc {
	h = SetRuleId(h, ruleIDX)
	return h
}

func SetRuleId(next http.HandlerFunc, ruleIDX int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "ruleID", ruleIDX)
		next(w, r.WithContext(ctx))
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	ruleID := r.Context().Value("ruleID").(int)
	fmt.Println("ruleID:", ruleID)

	rule := configs[ruleID]
	for _, body := range rule.Body {

	}

	if err := json.NewEncoder(w).Encode(r.Body); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func getValidationMap() error {
	c := map[string]interface{}{
		"Name":  "Qiang Xue",
		"Email": "q",
		"Address": map[string]interface{}{
			"Street": "123",
			"City":   "Unknown",
			"State":  "Virginia",
			"Zip":    "12345",
		},
	}

	err := vld.Validate.Var(c["Name"], "required,max=5")
	if err != nil {
		fmt.Printf("Name%s", vld.TextFromFirstError(err, "ru"))
		return nil
	}
	return nil
}

//func GetTranslator(lang string) (trans ut.Translator) {
//	translator, _ := translator.GetTranslator(lang)
//	return translator
//}
