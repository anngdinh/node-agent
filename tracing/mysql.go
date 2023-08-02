package tracing

import (
	// "encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"

	// "strings"
	"time"

	"github.com/coroot/coroot-node-agent/ebpftracer"
	"github.com/coroot/coroot-node-agent/tracing/mysql/proto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/klog/v2"
)

const (
	MysqlComQuery       = 3
	MysqlComStmtPrepare = 0x16
	MysqlComStmtExecute = 0x17
	MysqlComStmtClose   = 0x19
	mysqlMsgHeaderSize  = 4
)

func handleMysqlQuery(containerId string, start, end time.Time, r *ebpftracer.L7Request, attrs []attribute.KeyValue, preparedStatements map[string]string, c *Connection) {
	klog.Info("------------------ handleMysqlQuery")
	query := parseMysql(r.Payload[:], r.StatementId, preparedStatements, c)
	if query == "" {
		klog.Infof("---------- Can't: %s", query)
		return
	}

	klog.Infof("---------- query: %s", query)
	if c.Database != "" {
		attrs = append(attrs, attribute.KeyValue{Key: attribute.Key("db.name"), Value: attribute.StringValue(c.Database)})
	}
	if c.Username != "" {
		attrs = append(attrs, attribute.KeyValue{Key: attribute.Key("db.user"), Value: attribute.StringValue(c.Username)})
	}

	_, span := tracer.Start(nil, containerId+" | "+query, trace.WithTimestamp(start), trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(append(attrs, semconv.DBSystemMySQL, semconv.DBStatement(query))...)
	if r.Status == 500 {
		span.SetStatus(codes.Error, "")
	}
	span.End(trace.WithTimestamp(end))
}

func parseMysql(payload []byte, statementId uint32, preparedStatements map[string]string, c *Connection) string {
	klog.Info("------------------ parseMysql")
	payloadSize := len(payload)
	// klog.Info("---------- payloadSize:", payloadSize)

	// sEnc := base64.StdEncoding.EncodeToString(payload)
	// // klog.Info("---------- encode:", sEnc)

	// myString := string(payload[:])
	// klog.Info("---------- myString:", payload)
	// klog.Info("---------- myString:", myString)

	if payloadSize < mysqlMsgHeaderSize+5 {
		return ""
	}
	msgSize := int(payload[0]) | int(payload[1])<<8 | int(payload[2])<<16
	cmd := payload[4]
	readQuery := func() (query string) {
		to := mysqlMsgHeaderSize + msgSize
		partial := false
		if to > payloadSize-1 {
			to = payloadSize - 1
			partial = true
		}
		query = string(payload[mysqlMsgHeaderSize+1 : to])
		if partial {
			query += "..."
		}
		return query
	}
	readStatementId := func() string {
		return strconv.FormatUint(uint64(binary.LittleEndian.Uint32(payload[mysqlMsgHeaderSize+1:])), 10)
	}

	switch cmd {
	case MysqlComQuery:
		return readQuery()
	case MysqlComStmtExecute:
		statementIdStr := readStatementId()
		statement, ok := preparedStatements[statementIdStr]
		if !ok {
			statement = fmt.Sprintf(`EXECUTE %s /* unknown */`, statementIdStr)
		}
		return statement
	case MysqlComStmtPrepare:
		query := readQuery()
		statementIdStr := strconv.FormatUint(uint64(statementId), 10)
		preparedStatements[statementIdStr] = query
		return fmt.Sprintf("PREPARE %s FROM %s", statementIdStr, query)
	case MysqlComStmtClose:
		statementIdStr := readStatementId()
		delete(preparedStatements, statementIdStr)
	default:
		sum := 0
		for i := 13; i < 36; i++ {
			sum += int(payload[i])
		}
		if sum == 0 {
			a := proto.NewAuth()
			err := a.UnPack(payload[4:])
			klog.Info("---------- a:", a, err)

			c.Database = a.Database()
			c.Username = a.User()
		}
	}
	return ""
}
