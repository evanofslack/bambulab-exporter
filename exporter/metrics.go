package exporter

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "bambulab"
)

type metrics struct {
	amsEnabled             prom.Gauge
	amsFilamentRemain      *prom.GaugeVec
	amsHumid               prom.Gauge
	amsTemp                prom.Gauge
	cameraEnabled          prom.Gauge
	cameraTimelapseEnabled prom.Gauge
	chamberLight           prom.Gauge
	fanSpeed               *prom.GaugeVec
	// filamentWeightTotal    *prom.CounterVec
	gcodeState             *prom.GaugeVec
	layerNumber            prom.Gauge
	layerTarget            prom.Gauge
	nozzleDiameter         prom.Gauge
	nozzleSpeedLevel       prom.Gauge
	nozzleSpeedMag         prom.Gauge
	nozzleTargetTemp       prom.Gauge
	nozzleTemp             prom.Gauge
	nozzleType             *prom.GaugeVec
	printPercent           *prom.GaugeVec
	printTimeRemain        *prom.GaugeVec
	printsTotal            prom.Counter
	wifiSignal             prom.Gauge
}

func newMetrics(deviceId string) *metrics {
	constLabels := prom.Labels{"device": deviceId}
	m := &metrics{
		amsEnabled: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "ams_enabled_state",
				ConstLabels: constLabels,
			}),
		amsFilamentRemain: prom.NewGaugeVec(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "ams_filament_remaining",
				ConstLabels: constLabels,
			},
			[]string{"id", "material", "color"},
		),
		amsHumid: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "ams_humidity",
				ConstLabels: constLabels,
			}),
		amsTemp: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "ams_temperature",
				ConstLabels: constLabels,
			}),
		cameraEnabled: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "camera_enabled_state",
				ConstLabels: constLabels,
			}),
		cameraTimelapseEnabled: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "camera_timelapse_state",
				ConstLabels: constLabels,
			}),
		chamberLight: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "chamber_light_state",
				ConstLabels: constLabels,
			}),
		fanSpeed: prom.NewGaugeVec(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "fan_speed_percent",
				ConstLabels: constLabels,
			},
			[]string{"fan"},
		),
		// filamentWeightTotal: prom.NewCounterVec(
		// 	prom.CounterOpts{
		// 		Namespace:   namespace,
		// 		Name:        "filament_extruded_total_grams",
		// 		ConstLabels: constLabels,
		// 	},
		// 	[]string{"id", "material", "color"},
		// ),
		gcodeState: prom.NewGaugeVec(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "gcode_state",
				ConstLabels: constLabels,
			},
			[]string{"state"},
		),
		layerNumber: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "layer_number",
				ConstLabels: constLabels,
			}),
		layerTarget: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "layer_number_target",
				ConstLabels: constLabels,
			}),
		nozzleDiameter: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "nozzle_diameter",
				ConstLabels: constLabels,
			}),
		nozzleSpeedLevel: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "nozzle_speed_level",
				ConstLabels: constLabels,
			}),
		nozzleSpeedMag: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "nozzle_speed_magnitude",
				ConstLabels: constLabels,
			}),
		nozzleTemp: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "nozzle_temperature",
				ConstLabels: constLabels,
			}),
		nozzleTargetTemp: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "nozzle_temperature_target",
				ConstLabels: constLabels,
			}),
		nozzleType: prom.NewGaugeVec(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "nozzle_type_state",
				ConstLabels: constLabels,
			},
			[]string{"type"},
		),

		printPercent: prom.NewGaugeVec(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "print_percent",
				ConstLabels: constLabels,
			},
			[]string{"model"},
		),
		printTimeRemain: prom.NewGaugeVec(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "print_time_remaining_minutes",
				ConstLabels: constLabels,
			},
			[]string{"model"},
		),
		printsTotal: prom.NewCounter(
			prom.CounterOpts{
				Namespace:   namespace,
				Name:        "print_completed_total",
				ConstLabels: constLabels,
			}),
		wifiSignal: prom.NewGauge(
			prom.GaugeOpts{
				Namespace:   namespace,
				Name:        "wifi_signal",
				ConstLabels: constLabels,
			}),
	}
	return m
}
