package binance

import (
	"bot/common"
	"bytes"
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func GetP2PUsdtHighestPrice() (float64, error) {
	url := "https://c2c.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	method := "POST"
	payload := strings.NewReader(`{"page":1,"rows":10,"payTypes":["MoMo"],"asset":"USDT","tradeType":"SELL","fiat":"VND","publisherType":null}`)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return 0, err
	}

	req.Header.Add("c2ctype", "c2c_merchant")
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	highestPrice := gjson.Get(string(body), "data.0.adv.price")
	if !highestPrice.Exists() {
		return 0, errors.New("failed to get highest price")
	}

	return highestPrice.Float(), nil
}

func GetP2PUsdtHighestPriceReport() (string, error) {
	price, err := GetP2PUsdtHighestPrice()
	if err != nil {
		return "", err
	}
	var tmpl *template.Template
	tmpl = template.Must(template.ParseFiles("templates/usdt-price.txt"))

	var report common.Report
	//init the loc
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	//set timezone
	now := time.Now().In(loc)
	report.DateTime = now.Format("02-01-2006 15:04:05 MST")

	report.UsdtPrice = price

	var output bytes.Buffer
	// Execute template
	if err = tmpl.Execute(&output, report); err != nil {
		return "", err
	}
	return output.String(), nil
}
