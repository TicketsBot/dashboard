package chatreplica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSHandshakeTimeout: time.Second * 3, // We're not using TLS anyway
	},
	Timeout: time.Second * 3,
}

func Render(payload Payload) ([]byte, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
        return nil, err
    }

	res, err := client.Post(config.Conf.Bot.RenderServiceUrl, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
        return nil, err
    }

	fmt.Println(string(encoded))

	if res.StatusCode != 200 {
        return nil, fmt.Errorf("render service returned status code %d", res.StatusCode)
    }

	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
        return nil, err
    }

	return bytes, nil
}
