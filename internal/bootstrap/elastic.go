package bootstrap

import (
	"log"
	"net/http"
	"os"
	"skeleton/internal/elasticsearch"
	"skeleton/internal/variable"
	"time"

	"github.com/olivere/elastic/v7"
)

func InitElastic() (*elasticsearch.Elasticsearch, error) {
	conf := variable.Config.Get("Elastic").(map[string]any)
	if !(conf["enable"].(bool)) {
		return nil, nil
	}
	urls := conf["urls"].([]string)
	optins := []elastic.ClientOptionFunc{
		elastic.SetURL(urls...),
		elastic.SetBasicAuth(conf["user"].(string), conf["secret"].(string)),
		elastic.SetGzip(conf["gzip"].(bool)),
		// 是否转换请求地址，默认为true,当等于true时 请求http://ip:port/_nodes/http，将其返回的url作为请求路径
		elastic.SetSniff(conf["sniffer"].(bool)),
		elastic.SetHealthcheckInterval(time.Second * (time.Duration(conf["healthcheck"].(int64)))),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   conf["maxidle"].(int),
				ResponseHeaderTimeout: (time.Duration(conf["timeout"].(int64))) * time.Second,
			},
		}),
		elastic.SetErrorLog(log.New(os.Stderr, "ES-ERROR ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
	}
	return elasticsearch.NewElastic(optins...)
}
