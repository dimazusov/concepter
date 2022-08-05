package morph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const parseUrl = "http://%s:%s/parse/%s"

type parser struct {
	host string
	port string
}

func NewParser(host, port string) *parser {
	return &parser{
		host: host,
		port: port,
	}
}

func (m parser) Parse(word string) (wordForms []Form, err error) {
	resp, err := http.Get(fmt.Sprintf(parseUrl, m.host, m.port, word))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &wordForms)
	if err != nil {
		return nil, err
	}

	return wordForms, nil
}
