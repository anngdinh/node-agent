package tracing

import (
	"context"
	"os"
	"time"

	"github.com/coroot/coroot-node-agent/ebpftracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"inet.af/netaddr"
	"k8s.io/klog/v2"
)

var (
	tracer trace.Tracer
	space  = []byte{' '}
	crlf   = []byte{'\r', '\n'}
)

type Connection struct {
	Pid        uint32
	SrcAddr    netaddr.IPPort
	DstAddr    netaddr.IPPort
	Fd         uint64
	Timestamp  uint64
	Database   string
	Username   string
	Password   string
	PluginName string
}

func Init(machineId, hostname, version string) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		klog.Infoln("no OpenTelemetry collector endpoint configured")
		return
	}
	klog.Infoln("OpenTelemetry collector endpoint:", endpoint)

	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		klog.Exitln(err)
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("coroot-node-agent"),
			semconv.HostName(hostname),
			semconv.HostID(machineId),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	tracer = tracerProvider.Tracer("coroot-node-agent", trace.WithInstrumentationVersion(version))
}

func HandleL7Request(containerId string, src netaddr.IPPort, dest netaddr.IPPort, r *ebpftracer.L7Request, preparedStatements map[string]string, c *Connection) {
	if tracer == nil {
		return
	}
	end := time.Now()
	start := end.Add(-r.Duration)

	attrs := []attribute.KeyValue{
		semconv.ContainerID(containerId),
		semconv.NetPeerName(src.IP().String()),
		semconv.NetPeerPort(int(src.Port())),
		semconv.NetHostName(dest.IP().String()),
		semconv.NetHostPort(int(dest.Port())),
	}
	switch r.Protocol {
	case ebpftracer.L7ProtocolMysql:
		handleMysqlQuery(containerId, start, end, r, attrs, preparedStatements, c)
	default:
		return
	}
}
