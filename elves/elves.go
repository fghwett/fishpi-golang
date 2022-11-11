package elves

import (
	"crypto/md5"
	"fishpi/logger"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	addr = `https://fish.elves.online/mc/call/%s/%d/%s` // https://fish.elves.online/mc/call/{user}/{salt}/{sign} md5(user+token+salt)
)

type Elves struct {
	name  string
	token string

	logger logger.Logger
}

func NewElves(name, token string, logger logger.Logger) *Elves {
	e := &Elves{
		name:  name,
		token: token,

		logger: logger,
	}

	return e
}

func (e *Elves) HandleCall(data interface{}) {
	if err := e.call(); err != nil {
		e.logger.Logf("call stick error: %s", err)
	}
}

func (e *Elves) call() error {
	salt := time.Now().UnixMilli()
	sign := e.md5(fmt.Sprintf("%s%s%d", e.name, e.token, salt))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	u := fmt.Sprintf(addr, e.name, salt, sign)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return err
	}
	e.logger.Logf("call stick result: %s", string(body))

	return nil
}

func (e *Elves) md5(str string) string {
	newSig := md5.Sum([]byte(str))
	newArr := fmt.Sprintf("%x", newSig)
	return newArr
}
