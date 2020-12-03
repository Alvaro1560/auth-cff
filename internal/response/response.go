package response

import (
	"encoding/json"
	"net/http"
)

// Model estructura para enviar la respuesta general al cliente de la api
type Model struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Type  string      `json:"type"`
	Msg   string      `json:"msg"`
}

// Set Asigna valores a todos los campos de la estructura
func (rm *Model) Set(err bool, data interface{}, Code int) {
	rm.Error = err
	rm.Data = data
	rm.Code = Code
}

func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
