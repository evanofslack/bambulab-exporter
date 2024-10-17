package exporter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	monitor "github.com/evanofslack/bambulab-client/monitor"
	mqtt "github.com/evanofslack/bambulab-client/mqtt"
)

const (
	metricsPath     = "/metrics"
	shutdownTimeout = 5 * time.Second
)

type Exporter struct {
	registry *prom.Registry
	mon      *monitor.Monitor
	metrics  *metrics
	file     string
	server   *http.Server
	lastPercent int
}

func New(mon *monitor.Monitor, deviceId string) (*Exporter, error) {
	registry := prom.NewRegistry()
	metrics := newMetrics(deviceId)
	e := &Exporter{
		registry: registry,
		mon:      mon,
		metrics:  metrics,
	}
	e.register()
	return e, nil
}

func (e *Exporter) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.mon.Update:
			msg := e.mon.State
			if msg != nil && msg.Print != nil {
				e.record(*msg.Print)
			}
		}
	}
}

func (e *Exporter) Serve(port string) error {
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{}))
	addr := ":" + port
	e.server = &http.Server{Addr: addr, Handler: mux}
	fmt.Printf("serving prometheus metrics server at address %s\n", e.server.Addr)
	return e.server.ListenAndServe()
}

func (e *Exporter) Close() error {
	fmt.Println("starting prometheus metrics server close")
	defer fmt.Println("closed prometheus metrics server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return e.server.Shutdown(ctx)
}

func (e *Exporter) register() {
	e.registry.MustRegister(e.metrics.amsEnabled)
	e.registry.MustRegister(e.metrics.amsFilamentRemain)
	e.registry.MustRegister(e.metrics.amsHumid)
	e.registry.MustRegister(e.metrics.amsTemp)
	e.registry.MustRegister(e.metrics.cameraEnabled)
	e.registry.MustRegister(e.metrics.cameraTimelapseEnabled)
	e.registry.MustRegister(e.metrics.chamberLight)
	e.registry.MustRegister(e.metrics.fanSpeed)
	// e.registry.MustRegister(e.metrics.filamentWeightTotal)
	e.registry.MustRegister(e.metrics.gcodeState)
	e.registry.MustRegister(e.metrics.layerNumber)
	e.registry.MustRegister(e.metrics.layerTarget)
	e.registry.MustRegister(e.metrics.nozzleDiameter)
	e.registry.MustRegister(e.metrics.nozzleSpeedLevel)
	e.registry.MustRegister(e.metrics.nozzleSpeedMag)
	e.registry.MustRegister(e.metrics.nozzleTemp)
	e.registry.MustRegister(e.metrics.nozzleTargetTemp)
	e.registry.MustRegister(e.metrics.nozzleType)
	e.registry.MustRegister(e.metrics.printPercent)
	e.registry.MustRegister(e.metrics.printTimeRemain)
	e.registry.MustRegister(e.metrics.printsTotal)
	e.registry.MustRegister(e.metrics.wifiSignal)
}

// Get latest state and update prom metrics
func (e *Exporter) record(p mqtt.Print) {
	// First determine the name of the file printing
	if file := p.GcodeFile; file != nil {
		e.file = *file
	}
	e.recordAms(p.Ams)
	e.recordCamera(p.Ipcam)
	e.recordsLights(p.LightsReport)
	e.recordFans(p)
	e.recordGcode(p)
	e.recordLayer(p)
	e.recordNozzle(p)
	e.recordPrint(p)
	e.recordWifi(p)
}

func (e *Exporter) recordAms(a *mqtt.Ams) {
	if a == nil {
		return
	}
	if enabled := a.PowerOnFlag; enabled != nil {
		if *enabled {
			e.metrics.amsEnabled.Set(1)
		} else {
			e.metrics.amsEnabled.Set(0)
		}
	}
	if a.Ams != nil && len(*a.Ams) > 0 {
		innerSlice := *a.Ams
		inner := innerSlice[0]
		if humid := inner.Humidity; humid != nil {
			if n, err := strconv.ParseFloat(*humid, 64); err == nil {
				e.metrics.amsHumid.Set(n)
			}
		}
		if temp := inner.Temp; temp != nil {
			if n, err := strconv.ParseFloat(*temp, 64); err == nil {
				e.metrics.amsTemp.Set(n)
			}
		}
	}
}

func (e *Exporter) recordCamera(cam *mqtt.Ipcam) {
	if cam == nil {
		return
	}
	if record := cam.IpcamRecord; record != nil {
		if strings.ToLower(*record) == "enable" {
			e.metrics.cameraEnabled.Set(1)
		} else {
			e.metrics.cameraEnabled.Set(0)
		}
	}
	if timelapse := cam.Timelapse; timelapse != nil {
		if strings.ToLower(*timelapse) == "enable" {
			e.metrics.cameraTimelapseEnabled.Set(1)
		} else {
			e.metrics.cameraTimelapseEnabled.Set(0)
		}
	}
}

func (e *Exporter) recordsLights(lights *[]mqtt.LightsReport) {
	if lights == nil {
		return
	}
	if len(*lights) == 0 {
		return
	}
	for _, light := range *lights {
		if mode, node := light.Mode, light.Node; mode != nil && node != nil {
			if strings.ToLower(*node) == "chamber_light" {
				if strings.ToLower(*mode) == "on" {
					e.metrics.chamberLight.Set(1)
				} else {
					e.metrics.chamberLight.Set(1)
				}
			}
		}
	}
}

func (e *Exporter) recordFans(p mqtt.Print) {
	if aux := p.BigFan1Speed; aux != nil {
		if n, err := strconv.ParseFloat(*aux, 64); err == nil {
			e.metrics.fanSpeed.WithLabelValues("auxiliary").Set(n)
		}
	}
	if chamber := p.BigFan2Speed; chamber != nil {
		if n, err := strconv.ParseFloat(*chamber, 64); err == nil {
			e.metrics.fanSpeed.WithLabelValues("chamber").Set(n)
		}
	}
	if part := p.CoolingFanSpeed; part != nil {
		if n, err := strconv.ParseFloat(*part, 64); err == nil {
			e.metrics.fanSpeed.WithLabelValues("part").Set(n)
		}
	}
	if hotend := p.HeatbreakFanSpeed; hotend != nil {
		if n, err := strconv.ParseFloat(*hotend, 64); err == nil {
			e.metrics.fanSpeed.WithLabelValues("hotend").Set(n)
		}
	}
}

func (e *Exporter) recordGcode(p mqtt.Print) {
	if state := p.GcodeState; state != nil {
		e.metrics.gcodeState.WithLabelValues(*state).Set(1)
	}
}

func (e *Exporter) recordLayer(p mqtt.Print) {
	if layer := p.LayerNum; layer != nil {
		e.metrics.layerNumber.Set(float64(*layer))
	}
	if target := p.TotalLayerNum; target != nil {
		e.metrics.layerTarget.Set(float64(*target))
	}
}

func (e *Exporter) recordNozzle(p mqtt.Print) {
	if dia := p.NozzleDiameter; dia != nil {
		if n, err := strconv.ParseFloat(*dia, 64); err == nil {
			e.metrics.nozzleDiameter.Set(n)
		}
	}
	if lvl := p.SpdLvl; lvl != nil {
		e.metrics.nozzleSpeedLevel.Set(float64(*lvl))
	}
	if mag := p.SpdMag; mag != nil {
		e.metrics.nozzleSpeedMag.Set(float64(*mag))
	}
	if temp := p.NozzleTemper; temp != nil {
		e.metrics.nozzleTemp.Set(*temp)
	}
	if target := p.NozzleTargetTemper; target != nil {
		e.metrics.nozzleTargetTemp.Set(float64(*target))
	}
	if ty := p.NozzleType; ty != nil {
		e.metrics.nozzleType.WithLabelValues(*ty).Set(1)
	}
}

func (e *Exporter) recordPrint(p mqtt.Print) {
	percent := p.McPercent
	if percent != nil {
		e.metrics.printPercent.WithLabelValues(e.file).Set(float64(*percent))
	}
	if remain := p.McRemainingTime; remain != nil {
		e.metrics.printTimeRemain.WithLabelValues(e.file).Set(float64(*remain))
	}
	// Did we just complete a print?
	if percent != nil && *percent == 100 && e.lastPercent != 100 {
	    e.metrics.printsTotal.Inc()
	}
	if percent != nil {
        e.lastPercent = *percent
	}
}

func (e *Exporter) recordWifi(p mqtt.Print) {
	if wifi := p.WifiSignal; wifi != nil {
		w := strings.TrimSuffix(*wifi, "dBm")
		if n, err := strconv.ParseFloat(w, 64); err == nil {
			e.metrics.wifiSignal.Set(n)
		}
	}
}
