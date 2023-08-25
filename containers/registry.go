package containers

import (
	"time"
	"fmt"

	"github.com/coroot/coroot-node-agent/ebpftracer"
	"github.com/coroot/coroot-node-agent/flags"
	"github.com/coroot/coroot-node-agent/tracing"
	"github.com/prometheus/client_golang/prometheus"

	"k8s.io/klog/v2"
)

var (
	gcInterval = 10 * time.Minute
)

type Registry struct {
	reg prometheus.Registerer

	tracer *ebpftracer.Tracer
	events chan ebpftracer.Event

	pidConnMySQL       []uint32
	ConnectionById    map[uint64]*tracing.Connection
	preparedStatements map[string]string
}

func NewRegistry(reg prometheus.Registerer, kernelVersion string) (*Registry, error) {

	cs := &Registry{
		reg:             reg,
		events:          make(chan ebpftracer.Event, 10000),
		ConnectionById: map[uint64]*tracing.Connection{},

		pidConnMySQL:       []uint32{},
		preparedStatements: map[string]string{},
	}

	go cs.handleEvents(cs.events)
	t, err := ebpftracer.NewTracer(cs.events, kernelVersion, *flags.DisableL7Tracing)
	if err != nil {
		close(cs.events)
		return nil, err
	}
	cs.tracer = t

	return cs, nil
}

func (r *Registry) Close() {
	r.tracer.Close()
	close(r.events)
}

func (r *Registry) handleEvents(ch <-chan ebpftracer.Event) {
	gcTicker := time.NewTicker(gcInterval)
	defer gcTicker.Stop()
	for {
		select {

		case e, more := <-ch:
			fmt.Print("")
			// if e.Type != ebpftracer.EventTypeFileOpen {
			// if e.Type != ebpftracer.EventTypeFileOpen && e.Type != ebpftracer.EventTypeProcessStart && e.Type != ebpftracer.EventTypeProcessExit {
			// if e.Type == ebpftracer.EventTypeL7Request {
				// 	// // if e.Pid == 145575 || e.DstAddr.String() == "172.24.0.2:3306" {

				// klog.Info("---- type:", e.Type, "reason:", e.Reason, "pid:", e.Pid, "src:", e.SrcAddr, "dst:", e.DstAddr, "fd:", e.Fd, "ts:", e.Timestamp)
				fmt.Println("---- type:", e.Type, "reason:", e.Reason, "pid:", e.Pid, "src:", e.SrcAddr, "dst:", e.DstAddr, "fd:", e.Fd, "ts:", e.Timestamp, "l7:")
			// }

			if !more {
				return
			}

			// if slices.Contains(r.pidConnMySQL, e.Pid) {
			// 	// mysqlPid := r.containersById["/system.slice/mysql.service"].id
			// 	e.Pid = 1317
			// }

			switch e.Type {
			// case ebpftracer.EventTypeProcessStart:
			// 	c, seen := r.ConnectionById[e.Pid]
			// 	switch {
			// 	case c != nil || seen: // revalidating by cgroup
			// 		delete(r.ConnectionById, e.Pid)
			// 	}
			// 	r.ConnectionById[e.Pid] = &Connection{
			// 		Pid:       e.Pid,
			// 		SrcAddr:   e.SrcAddr,
			// 		DstAddr:   e.DstAddr,
			// 		Fd:        e.Fd,
			// 		Timestamp: e.Timestamp,
			// 	}
			// case ebpftracer.EventTypeProcessExit:
			// 	delete(r.ConnectionById, e.Pid)

			case ebpftracer.EventTypeConnectionOpen:
				// if e.DstAddr.String()[len(e.DstAddr.String())-4:] == "3306" {
					if c, seen := r.ConnectionById[e.Id]; c != nil || seen {
						delete(r.ConnectionById, e.Id)
					}
					r.ConnectionById[e.Id] = &tracing.Connection{
						Pid:       e.Pid,
						// SrcAddr:   e.SrcAddr,
						// DstAddr:   e.DstAddr,
						// Fd:        e.Fd,
						// Timestamp: e.Timestamp,
					}
				// }
			// case ebpftracer.EventTypeConnectionError:
			// 	if c := r.getOrCreateContainer(e.Pid); c != nil {
			// 		c.onConnectionOpen(e.Pid, e.Fd, e.SrcAddr, e.DstAddr, 0, true)
			// 	} else {
			// 		klog.Infoln("TCP connection error from unknown container", e)
			// 	}
			// case ebpftracer.EventTypeConnectionClose:
			// 	srcDst := AddrPair{src: e.SrcAddr, dst: e.DstAddr}
			// 	for _, c := range r.containersById {
			// 		if c.onConnectionClose(srcDst) {
			// 			break
			// 		}
			// 	}
			// case ebpftracer.EventTypeTCPRetransmit:
			// 	srcDst := AddrPair{src: e.SrcAddr, dst: e.DstAddr}
			// 	for _, c := range r.containersById {
			// 		if c.onRetransmit(srcDst) {
			// 			break
			// 		}
			// 	}
			case ebpftracer.EventTypeL7Request:
				if e.L7Request == nil {
					continue
				}
				if c, seen := r.ConnectionById[e.Id]; c != nil || seen {
					delete(r.ConnectionById, e.Id)
				}
				r.ConnectionById[e.Id] = &tracing.Connection{
					Pid:       e.Pid,
					// SrcAddr:   e.SrcAddr,
					// DstAddr:   e.DstAddr,
					// Fd:        e.Fd,
					// Timestamp: e.Timestamp,

				}
				// klog.Info("--------- EventTypeL7Request: ")
				// if c := r.ConnectionById[e.Pid]; c != nil {
					r.onL7Request(e.Id, e.Fd, e.Timestamp, e.L7Request)
				// } else {
				// 	klog.Infof("--------- no container found for pid:", e.Pid)
				// }
			}
		}
	}
}

func (r *Registry) onL7Request(Id uint64, fd uint64, timestamp uint64, re *ebpftracer.L7Request) {
	klog.Info("-------------- onL7Request")
	c := r.ConnectionById[Id]

	tracing.HandleL7Request("mysql.process", re, r.preparedStatements, c)
}
