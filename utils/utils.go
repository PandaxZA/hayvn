package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Netflix/go-env"
)

type Environment struct {
	HOST_NAME          string `env:"HOST_NAME,default=http://localhost"`
	HOST_PORT          int    `env:"HOST_PORT,default=8081"`
	ENVIRONMENT_NAME   string `env:"ENVIRONMENT_NAME,default=development"`
	RESPONSE_URL       string `env:"RESPONSE_URL,default=http://localhost"`
	RESPONSE_PORT      int    `env:"RESPONSE_PORT,default=8082"`
	RATE_LIMIT_SECONDS int    `env:"RATE_LIMIT_SECONDS,default=10"`

	IS_DEBUG bool `env:"IS_DEBUG,default=false"`
	IS_LOCAL bool `env:"IS_LOCAL,default=false"`
}

var CurrentEnvironment Environment

func MakeEnvironment() Environment {
	// Setup env variables string

	_, err := env.UnmarshalFromEnviron(&CurrentEnvironment)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}

	return CurrentEnvironment

}

func WWWFormUrlEncodedToURLValues(r *http.Request) (url.Values, error) {
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	params, err := url.ParseQuery(string(buf))
	if err != nil {
		return nil, err
	}
	urlVals := make(url.Values)
	for key, val := range params {
		if len(val[0]) > 0 {
			urlVals[key] = []string{val[0]}
		}
	}
	return urlVals, nil
}
