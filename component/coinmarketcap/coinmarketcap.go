package coinmarketcap

import (
"fmt"
"io/ioutil"
"log"
"net/http"
"net/url"
"os"
)

func Ahihi() {
	client := &http.Client{}
	req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "56080fd4-dbf5-4ab0-a6a7-45c4007b7e4f")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}