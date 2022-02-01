package two_miners

import (
	"bot/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type WorkerStats struct {
	LastBeat        int  `json:"lastBeat"`
	CurrentHashrate int  `json:"hr"`
	Offline         bool `json:"offline"`
	AverageHashrate int  `json:"hr2"`
	ReportHashRate  int  `json:"rhr"`
	SharesValid     int  `json:"sharesValid"`
	SharesInvalid   int  `json:"sharesInvalid"`
	SharesStale     int  `json:"sharesStale"`
}

type TwoMinersStats struct {
	Two4Hnumreward int `json:"24hnumreward"`
	Two4Hreward    int `json:"24hreward"`
	APIVersion     int `json:"apiVersion"`
	Config         struct {
		AllowedMaxPayout int64  `json:"allowedMaxPayout"`
		AllowedMinPayout int    `json:"allowedMinPayout"`
		DefaultMinPayout int    `json:"defaultMinPayout"`
		IPHint           string `json:"ipHint"`
		IPWorkerName     string `json:"ipWorkerName"`
		MinPayout        int    `json:"minPayout"`
		PaymentHubHint   string `json:"paymentHubHint"`
	} `json:"config"`
	CurrentHashrate int    `json:"currentHashrate"`
	CurrentLuck     string `json:"currentLuck"`
	Hashrate        int    `json:"hashrate"`
	PageSize        int    `json:"pageSize"`
	Payments        []struct {
		Amount    int    `json:"amount"`
		Timestamp int    `json:"timestamp"`
		Tx        string `json:"tx"`
		TxFee     int    `json:"txFee"`
	} `json:"payments"`
	PaymentsTotal int `json:"paymentsTotal"`
	Rewards       []struct {
		Blockheight int     `json:"blockheight"`
		Timestamp   int     `json:"timestamp"`
		Reward      int     `json:"reward"`
		Percent     float64 `json:"percent"`
		Immature    bool    `json:"immature"`
		Orphan      bool    `json:"orphan"`
		Uncle       bool    `json:"uncle"`
	} `json:"rewards"`
	RoundShares   int `json:"roundShares"`
	SharesInvalid int `json:"sharesInvalid"`
	SharesStale   int `json:"sharesStale"`
	SharesValid   int `json:"sharesValid"`
	Stats         struct {
		Balance   int `json:"balance"`
		Immature  int `json:"immature"`
		LastShare int `json:"lastShare"`
		Paid      int `json:"paid"`
		Pending   int `json:"pending"`
	} `json:"stats"`
	Sumrewards []struct {
		Inverval  int    `json:"inverval"`
		Reward    int    `json:"reward"`
		Numreward int    `json:"numreward"`
		Name      string `json:"name"`
		Offset    int    `json:"offset"`
	} `json:"sumrewards"`
	UpdatedAt      int64                  `json:"updatedAt"`
	Workers        map[string]WorkerStats `json:"workers"`
	WorkersOffline int                    `json:"workersOffline"`
	WorkersOnline  int                    `json:"workersOnline"`
	WorkersTotal   int                    `json:"workersTotal"`
}

func GetStats(cfg *config.Config) (*TwoMinersStats, error) {
	url := "https://eth.2miners.com/api/accounts/" + cfg.WalletAddress

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("cannot fetch data from 2miners.com")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	stats := TwoMinersStats{}
	if err = json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
