package main

import (
	"context"
	"log"

	wemo "github.com/gecgooden/go.wemo"
	"github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace string = "wemo"
)

// TorrentCollector has a transmission.Client to create torrent metrics
type TorrentCollector struct {
	client *transmission.Client

	Status   *prometheus.Desc
	Added    *prometheus.Desc
	Files    *prometheus.Desc
	Finished *prometheus.Desc
	Done     *prometheus.Desc
	Ratio    *prometheus.Desc
	Download *prometheus.Desc
	Upload   *prometheus.Desc

	// TrackerStats
	Downloads *prometheus.Desc
	Leechers  *prometheus.Desc
	Seeders   *prometheus.Desc
}

// WemoCollector has a wemo.Device to create wemo metrics
type WemoCollector struct {
	device *wemo.Device

	OnFor          *prometheus.Desc
	OnToday        *prometheus.Desc
	OnTotal        *prometheus.Desc
	WifiStrength   *prometheus.Desc
	CurrentPower   *prometheus.Desc
	TodayPower     *prometheus.Desc
	TotalPower     *prometheus.Desc
	PowerThreshold *prometheus.Desc
}

// NewWemoCollector creates a new wemo device collector with the wemo.Device
func NewWemoCollector(device *wemo.Device) *WemoCollector {
	return &WemoCollector{
		device: device,

		OnFor: prometheus.NewDesc(
			namespace+"_on_for",
			"Time on in seconds since last boot",
			[]string{"device_name", "device_type"},
			nil,
		),
		OnToday: prometheus.NewDesc(
			namespace+"_on_today",
			"Time on in seconds today",
			[]string{"device_name", "device_type"},
			nil,
		),
		OnTotal: prometheus.NewDesc(
			namespace+"_on_total",
			"Time on in seconds in total for device",
			[]string{"device_name", "device_type"},
			nil,
		),
		WifiStrength: prometheus.NewDesc(
			namespace+"_wifi_strength",
			"Wifi Strength in RSSI",
			[]string{"device_name", "device_type"},
			nil,
		),
		CurrentPower: prometheus.NewDesc(
			namespace+"_current_power",
			"Current Power usage in mW",
			[]string{"device_name", "device_type"},
			nil,
		),
		TodayPower: prometheus.NewDesc(
			namespace+"_today_power",
			"Total Power today in mW",
			[]string{"device_name", "device_type"},
			nil,
		),
		TotalPower: prometheus.NewDesc(
			namespace+"_total_power",
			"Total Power in mW",
			[]string{"device_name", "device_type"},
			nil,
		),
		PowerThreshold: prometheus.NewDesc(
			namespace+"_power_threshold",
			"Power Threshold in mW",
			[]string{"device_name", "device_type"},
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface
func (wc *WemoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- wc.OnFor
	ch <- wc.OnToday
	ch <- wc.OnTotal
	ch <- wc.WifiStrength
	ch <- wc.CurrentPower
	ch <- wc.TodayPower
	ch <- wc.TotalPower
	ch <- wc.PowerThreshold
}

// Collect implements the prometheus.Collector interface
func (wc *WemoCollector) Collect(ch chan<- prometheus.Metric) {
	deviceInfo, err := wc.device.FetchDeviceInfo(context.Background())
	if err != nil {
		log.Printf("Failed to get Wemo Device Info: %v", err)
		return
	}

	insightData, err := wc.device.GetInsightParams()
	if err != nil {
		log.Printf("Failed to get Wemo Insight Info: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		wc.OnFor,
		prometheus.GaugeValue,
		float64(insightData.OnFor),
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.OnToday,
		prometheus.GaugeValue,
		float64(insightData.OnToday),
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.OnTotal,
		prometheus.CounterValue,
		float64(insightData.OnTotal),
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.WifiStrength,
		prometheus.GaugeValue,
		insightData.WifiStrength,
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.CurrentPower,
		prometheus.GaugeValue,
		insightData.CurrentPower,
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.TodayPower,
		prometheus.GaugeValue,
		insightData.TodayPower,
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.TotalPower,
		prometheus.GaugeValue,
		insightData.TotalPower,
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)

	ch <- prometheus.MustNewConstMetric(
		wc.PowerThreshold,
		prometheus.GaugeValue,
		insightData.PowerThreshold,
		deviceInfo.FriendlyName, deviceInfo.DeviceType,
	)
}
