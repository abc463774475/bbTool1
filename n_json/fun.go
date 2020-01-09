package n_json

import (
	"bbTool/n_log"
	"encoding/json"
	"reflect"
)

func Unmarshal(data []byte, i interface{}) error {
	if err := json.Unmarshal(data, i); err != nil {
		n_log.Erro_special(3,"json unmarshal  %v\n%v\n%v", err, string(data), reflect.TypeOf(i))
		return err
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(v); err != nil {
		n_log.Erro_special(3,"json marshal not right  ok ?  %v  :  %v", err, v)
		n_log.Panic("pppppp")
		return data, err
	}
	return data, err
}
