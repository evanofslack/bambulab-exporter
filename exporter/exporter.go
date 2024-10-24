package exporter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	monitor "github.com/evanofslack/bambulab-client/monitor"
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
	go e.handleEvent(ctx)
	go e.handleUpdate(ctx)
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
	e.registry.MustRegister(e.metrics.amsPowered)
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

func (e *Exporter) handleEvent(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.mon.PrintStarted:
			e.metrics.printsTotal.WithLabelValues("start").Inc()
		case <-e.mon.PrintFinished:
			e.metrics.printsTotal.WithLabelValues("finish").Inc()
		case <-e.mon.PrintCancelled:
			e.metrics.printsTotal.WithLabelValues("cancel").Inc()
		case <-e.mon.PrintFailed:
			e.metrics.printsTotal.WithLabelValues("fail").Inc()
		}
	}
}

func (e *Exporter) handleUpdate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.mon.Update:
			state := e.mon.CurrentState()
			e.record(state)
		}
	}
}

// Get latest state and update prom metrics
func (e *Exporter) record(s monitor.State) {
	// First determine the name of the file printing
	if file, err := s.Gcode.File.Take(); err != nil {
		e.file = file
	}
	e.recordAms(s.Ams)
	e.recordCamera(s.Camera)
	e.recordsLights(s.Lights)
	e.recordFans(s.Fans)
	e.recordGcode(s.Gcode)
	e.recordNozzle(s.Nozzle)
	e.recordSpeed(s.Speed)
	e.recordPrint(s.CurrentPrint)
	e.recordWifi(s)
}

func (e *Exporter) recordAms(a monitor.Ams) {
	a.Enabled.IfSome(func(v bool) {
		if v {
			e.metrics.amsEnabled.Set(1)
		} else {
			e.metrics.amsEnabled.Set(0)
		}
	})
	a.Powered.IfSome(func(v bool) {
		if v {
			e.metrics.amsPowered.Set(1)
		} else {
			e.metrics.amsPowered.Set(0)
		}
	})
	for i, unit := range a.Units {
		unit.Humidity.IfSome(func(v float64) {
			e.metrics.amsHumid.WithLabelValues(fmt.Sprintf("%d", i)).Set(v)
		})
		unit.Temperature.IfSome(func(v float64) {
			e.metrics.amsTemp.WithLabelValues(fmt.Sprintf("%d", i)).Set(v)
		})
	}
}

func (e *Exporter) recordCamera(cam monitor.Camera) {
	cam.Recording.IfSome(func(v bool) {
		if v {
			e.metrics.cameraEnabled.Set(1)
		} else {
			e.metrics.cameraEnabled.Set(0)
		}
	})
	cam.Timelapse.IfSome(func(v bool) {
		if v {
			e.metrics.cameraTimelapseEnabled.Set(1)
		} else {
			e.metrics.cameraTimelapseEnabled.Set(0)
		}
	})
}

func (e *Exporter) recordsLights(lights monitor.Lights) {
	lights.Chamber.IfSome(func(v bool) {
		if v {
			e.metrics.chamberLight.Set(1)
		} else {
			e.metrics.chamberLight.Set(0)
		}
	})
}

func (e *Exporter) recordFans(fans monitor.Fans) {
	fans.Auxilliary.IfSome(func(v float64) {
		e.metrics.fanSpeed.WithLabelValues("auxiliary").Set(v)
	})
	fans.Chamber.IfSome(func(v float64) {
		e.metrics.fanSpeed.WithLabelValues("chamber").Set(v)
	})
	fans.Part.IfSome(func(v float64) {
		e.metrics.fanSpeed.WithLabelValues("part").Set(v)
	})
	fans.Hotend.IfSome(func(v float64) {
		e.metrics.fanSpeed.WithLabelValues("hotend").Set(v)
	})
}

func (e *Exporter) recordGcode(gcode monitor.Gcode) {
	gcode.State.IfSome(func(v string) {
		e.metrics.gcodeState.WithLabelValues(v).Set(1)
	})
}

func (e *Exporter) recordNozzle(nozzle monitor.Nozzle) {
	nozzle.Diameter.IfSome(func(v float64) {
		e.metrics.nozzleDiameter.Set(v)
	})
	nozzle.Temperature.IfSome(func(v float64) {
		e.metrics.nozzleTemp.Set(v)
	})
	nozzle.TemperatureTarget.IfSome(func(v int) {
		e.metrics.nozzleTargetTemp.Set(float64(v))
	})
	nozzle.Type.IfSome(func(v string) {
		e.metrics.nozzleType.WithLabelValues(v).Set(1)
	})
}

func (e *Exporter) recordSpeed(speed monitor.Speed) {
	speed.Magnitude.IfSome(func(v int) {
		e.metrics.nozzleSpeedMag.Set(float64(v))
	})
	speed.LevelName.IfSome(func(v string) {
		e.metrics.nozzleSpeedLevel.WithLabelValues(v).Set(1)
	})
}

func (e *Exporter) recordPrint(curr monitor.CurrentPrint) {
	curr.LayerNumber.IfSome(func(v int) {
		e.metrics.layerNumber.Set(float64(v))
	})
	curr.LayerNumberTarget.IfSome(func(v int) {
		e.metrics.layerTarget.Set(float64(v))
	})
	curr.Percent.IfSome(func(v int) {
		e.metrics.printPercent.WithLabelValues(e.file).Set(float64(v))
	})
	curr.TimeRemaining.IfSome(func(v int) {
		e.metrics.printTimeRemain.WithLabelValues(e.file).Set(float64(v))
	})
}

func (e *Exporter) recordWifi(s monitor.State) {
	s.Wifi.IfSome(func(v float64) {
		e.metrics.wifiSignal.Set(v)
	})
}
