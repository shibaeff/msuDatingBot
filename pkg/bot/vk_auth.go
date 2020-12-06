package bot

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	linkStub = "https://oauth.vk.com/authorize?client_id=7679100&scope=327682&&display=page&" +
		"response_type=code&v=5.126&state=123456&redirect_uri=" +
		"http://161.156.162.178:30000/check?tg_id=%d"
	pollingStub = "http://161.156.162.178:30000/answer?tg_id=%d"
	message     = "<a href=\"%s\">Подтвердите свою принадлежность к МГУ, пройдя по ссылке</a>"

	notKnown    = 3
	failedAuth  = 0
	successAuth = 1
)

func prepareLink(id int64) string {
	return fmt.Sprintf(linkStub, id)
}

func prepareMessage(link string) string {
	return fmt.Sprintf(message, link)
}

func getStatus(link string) int64 {
	resp, err := http.Get(link)
	if err != nil {
		return notKnown
	}
	for {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return notKnown
		}
		bodyString := string(bodyBytes)
		status := int64(notKnown)
		fmt.Sscanf(bodyString, "%d", &status)
		return status
	}
}

func pollAuthAPI(id int64) (verdict int64) {
	pollingLink := fmt.Sprintf(pollingStub, id)
	return getStatus(pollingLink)
}
