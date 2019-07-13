package xmetrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type MetricsIn struct {
	fx.In

	Viper *viper.Viper
}

type MetricsOut struct {
	fx.Out

	Registerer prometheus.Registerer
	Gatherer   prometheus.Gatherer
	Registry   Registry
	Handler    Handler
}

func Provide(configKey string, ho promhttp.HandlerOpts) func(MetricsIn) (MetricsOut, error) {
	return func(in MetricsIn) (MetricsOut, error) {
		var o Options
		if err := in.Viper.UnmarshalKey(configKey, &o); err != nil {
			return MetricsOut{}, err
		}

		registry, err := New(o)
		if err != nil {
			return MetricsOut{}, err
		}

		return MetricsOut{
			Registerer: registry,
			Gatherer:   registry,
			Registry:   registry,
			Handler:    promhttp.HandlerFor(registry, ho),
		}, nil
	}
}

func ProvideCounter(o prometheus.CounterOpts, labelNames ...string) fx.Annotated {
	return fx.Annotated{
		Name: o.Name,
		Target: func(r Registry) (metrics.Counter, error) {
			return r.NewCounter(o, labelNames)
		},
	}
}

func ProvideGauge(o prometheus.GaugeOpts, labelNames ...string) fx.Annotated {
	return fx.Annotated{
		Name: o.Name,
		Target: func(r Registry) (metrics.Gauge, error) {
			return r.NewGauge(o, labelNames)
		},
	}
}

func ProvideHistogram(o prometheus.HistogramOpts, labelNames ...string) fx.Annotated {
	return fx.Annotated{
		Name: o.Name,
		Target: func(r Registry) (metrics.Histogram, error) {
			return r.NewHistogram(o, labelNames)
		},
	}
}

func ProvideSummary(o prometheus.SummaryOpts, labelNames ...string) fx.Annotated {
	return fx.Annotated{
		Name: o.Name,
		Target: func(r Registry) (metrics.Histogram, error) {
			return r.NewSummary(o, labelNames)
		},
	}
}
