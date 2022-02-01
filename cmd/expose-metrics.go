package cmd

import (
	two_miners "bot/component/2miners"
	"errors"

	"bot/component/vhttos"
	"bot/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

func exposeMetrics() error {
	cfg, err := config.GetConfigs()
	if err != nil {
		return err
	}

	if cfg.AccessToken == "" {
		return errors.New("cannot get ACCESS_TOKEN")
	}

	if cfg.WalletAddress == "" {
		return errors.New("cannot get WALLET_ADDRESS")
	}

	gpuTempGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_temperature_celsius",
			Help: "Temperature of VGAs.",
		},
		[]string{
			"miner",
			"gpu",
		},
	)
	minerWorkerCurrentHashrateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "miner_worker_current_hashrate",
			Help: "Current hashrate of miner worker (last 30 minutes of work).",
		},
		[]string{
			"miner",
		},
	)
	minerWorkerAverageHashrateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "miner_worker_average_hashrate",
			Help: "Average hashrate of miner worker (last 6 hours of worker).",
		},
		[]string{
			"miner",
		},
	)
	minerWorkerReportHashrateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "miner_worker_report_hashrate",
			Help: "Report hashrate of miner worker.",
		},
		[]string{
			"miner",
		},
	)

	prometheus.MustRegister(gpuTempGauge)
	prometheus.MustRegister(minerWorkerCurrentHashrateGauge)
	prometheus.MustRegister(minerWorkerAverageHashrateGauge)
	prometheus.MustRegister(minerWorkerReportHashrateGauge)

	go func() {
		for {
			miners, err := vhttos.GetMiners(cfg)
			if err != nil {
				log.Fatalln(err)
			}

			for _, miner := range miners {
				for key, val := range miner.Data.GpuTemp {
					gpu := key
					celsius, err := strconv.ParseFloat(val, 64)
					if err != nil {
						log.Fatalln("cannot convert string to celsius (float64)")
					}

					if miner.Online {
						gpuTempGauge.With(prometheus.Labels{"miner": miner.Name, "gpu": gpu}).Set(celsius)
					} else {
						gpuTempGauge.With(prometheus.Labels{"miner": miner.Name, "gpu": gpu}).Set(0)
					}
				}
			}

			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			stats, err := two_miners.GetStats(cfg)
			if err != nil {
				log.Fatalln(err)
			}

			for key, worker := range stats.Workers {
				currentHashrate := float64(worker.CurrentHashrate)
				averageHashrate := float64(worker.AverageHashrate)
				reportHashrate := float64(worker.ReportHashRate)

				minerWorkerCurrentHashrateGauge.With(prometheus.Labels{"miner": key}).Set(currentHashrate)
				minerWorkerAverageHashrateGauge.With(prometheus.Labels{"miner": key}).Set(averageHashrate)
				minerWorkerReportHashrateGauge.With(prometheus.Labels{"miner": key}).Set(reportHashrate)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
	return nil
}
