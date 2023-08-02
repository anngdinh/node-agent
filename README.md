export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4318/v1/traces
export ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.20

// detected a new container" pid=1387 cg="/system.slice/mysql.service" id=/system.slice/mysql.service



panic: Something in this program imports go4.org/unsafe/assume-no-moving-gc to declare that it assumes a non-moving garbage collector, but your version of go4.org/unsafe/assume-no-moving-gc hasn't been updated to assert that it's safe against the go1.20 runtime. If you want to risk it, run with environment variable ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.20 set. Notably, if go1.20 adds a moving garbage collector, this program is unsafe to use.

goroutine 1 [running]:
go4.org/unsafe/assume-no-moving-gc.init.0()
        /root/go/pkg/mod/go4.org/unsafe/assume-no-moving-gc@v0.0.0-20220617031537-928513b29760/untested.go:25 +0x1ba
exit status 2