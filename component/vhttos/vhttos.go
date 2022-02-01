package vhttos

import (
	"bot/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Miner struct {
	ID    string `json:"id"`
	Mac   string `json:"mac"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Data  struct {
		GpuTemp         map[string]string `json:"gpuTemp"`
		GpuVddGfx       interface{}       `json:"gpuVddGfx"`
		GpuMvdd         interface{}       `json:"gpuMvdd"`
		GpuMvddci       interface{}       `json:"gpuMvddci"`
		SysLoad5        string            `json:"sysLoad5"`
		SysCPUModel     string            `json:"sysCpuModel"`
		ConsoleShort    string            `json:"consoleShort"`
		GpuModel        map[string]string `json:"gpuModel"`
		ImgVersionOs    string            `json:"imgVersionOs"`
		Driver          string            `json:"driver"`
		SysBios         string            `json:"sysBios"`
		GpuPwrLimit     map[string]string `json:"gpuPwrLimit"`
		GpuMemTemp      interface{}       `json:"gpuMemTemp"`
		GpuVramSize     map[string]string `json:"gpuVramSize"`
		GpuPwrMax       map[string]string `json:"gpuPwrMax"`
		GpuCoreClk      map[string]string `json:"gpuCoreClk"`
		GpuCount        string            `json:"gpuCount"`
		Rej             string            `json:"rej"`
		SysPwr          string            `json:"sysPwr"`
		GpuMemClk       map[string]string `json:"gpuMemClk"`
		GpuFan          map[string]string `json:"gpuFan"`
		GpuHash         interface{}       `json:"gpuHash"`
		GpuPciBus       interface{}       `json:"gpuPciBus"`
		GpuVramType     interface{}       `json:"gpuVramType"`
		GpuAsicTemp     interface{}       `json:"gpuAsicTemp"`
		SysMbo          string            `json:"sysMbo"`
		SysRAMSize      string            `json:"sysRamSize"`
		Hash            string            `json:"hash"`
		Acc             string            `json:"acc"`
		SysHdd          string            `json:"sysHdd"`
		GpuPwrCur       map[string]string `json:"gpuPwrCur"`
		Uptime          string            `json:"uptime"`
		GpuManufacturer map[string]string `json:"gpuManufacturer"`
		GpuPwrMin       interface{}       `json:"gpuPwrMin"`
		IPLAN           string            `json:"ipLAN"`
		IPPublic        string            `json:"ip_public"`
		GpuVramChip     interface{}       `json:"gpuVramChip"`
		GpuBiosVer      map[string]string `json:"gpuBiosVer"`
		Kernel          string            `json:"kernel"`
	} `json:"data"`
	OsSeries             string      `json:"os_series"`
	OsVersion            string      `json:"os_version"`
	Status               bool        `json:"status"`
	LastTimeUpdate       int         `json:"last_time_update"`
	IDGroupConfig        string      `json:"id_group_config"`
	UseSelfOc            bool        `json:"use_self_oc"`
	SelfOc               interface{} `json:"self_oc"`
	AdditionMinerOptions string      `json:"addition_miner_options"`
	IDGroupOc            string      `json:"id_group_oc"`
	IDGroupSchedule      string      `json:"id_group_schedule"`
	CountRestart         int         `json:"count_restart"`
	UpdatedAt            time.Time   `json:"updated_at"`
	CreatedAt            time.Time   `json:"created_at"`
	Online               bool        `json:"online"`
}

func GetMiners(cfg *config.Config) ([]Miner, error) {
	url := "https://vhttos.com/api/v2/rig/list"
	accessToken := cfg.AccessToken

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("cannot fetch data from VHTTOS")
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

	minerCount := gjson.Get(string(body), "data.miners.#")
	if !minerCount.Exists() {
		return nil, errors.New("failed to count miners")
	}

	var miners []Miner

	var i int64
	for i = 0; i < minerCount.Int(); i++ {
		miner := Miner{}
		jsonValue := gjson.Get(string(body), fmt.Sprintf("data.miners.%d", i))
		if !jsonValue.Exists() {
			log.Fatalln("failed to get JSON value (miner)")
		}
		if err := json.Unmarshal([]byte(jsonValue.String()), &miner); err != nil {
			return nil, err
		}
		miners = append(miners, miner)
	}

	return miners, nil
}
