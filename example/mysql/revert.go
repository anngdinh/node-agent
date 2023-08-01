package mysql

// import (
// 	"crypto/tls"
// 	"fmt"
// )

// func (mc *mysqlConn) UnwriteHandshakeResponsePacket(data []byte) ([]byte, string, error) {
// 	fmt.Println("writeHandshakeResponsePacket")
// 	// Adjust client flags based on server support
// 	clientFlags := clientProtocol41 |
// 		clientSecureConn |
// 		clientLongPassword |
// 		clientTransactions |
// 		clientLocalFiles |
// 		clientPluginAuth |
// 		clientMultiResults |
// 		clientConnectAttrs |
// 		mc.flags&clientLongFlag

// 	if mc.cfg.ClientFoundRows {
// 		clientFlags |= clientFoundRows
// 	}

// 	// To enable TLS / SSL
// 	if mc.cfg.TLS != nil {
// 		clientFlags |= clientSSL
// 	}

// 	if mc.cfg.MultiStatements {
// 		clientFlags |= clientMultiStatements
// 	}

// 	// encode length of the auth plugin data
// 	var authRespLEIBuf [9]byte
// 	authRespLen := len(authResp)
// 	authRespLEI := appendLengthEncodedInteger(authRespLEIBuf[:0], uint64(authRespLen))
// 	if len(authRespLEI) > 1 {
// 		// if the length can not be written in 1 byte, it must be written as a
// 		// length encoded integer
// 		clientFlags |= clientPluginAuthLenEncClientData
// 	}

// 	pktLen := 4 + 4 + 1 + 23 + len(mc.cfg.User) + 1 + len(authRespLEI) + len(authResp) + 21 + 1

// 	// To specify a db name
// 	if n := len(mc.cfg.DBName); n > 0 {
// 		clientFlags |= clientConnectWithDB
// 		pktLen += n + 1
// 	}

// 	// 1 byte to store length of all key-values
// 	// NOTE: Actually, this is length encoded integer.
// 	// But we support only len(connAttrBuf) < 251 for now because takeSmallBuffer
// 	// doesn't support buffer size more than 4096 bytes.
// 	// TODO(methane): Rewrite buffer management.
// 	pktLen += 1 + len(mc.connector.encodedAttributes)

// 	// Calculate packet length and get buffer with that size
// 	data, err := mc.buf.takeSmallBuffer(pktLen + 4)
// 	if err != nil {
// 		// cannot take the buffer. Something must be wrong with the connection
// 		mc.cfg.Logger.Print(err)
// 		return errBadConnNoWrite
// 	}

// 	// ClientFlags [32 bit]
// 	data[4] = byte(clientFlags)
// 	data[5] = byte(clientFlags >> 8)
// 	data[6] = byte(clientFlags >> 16)
// 	data[7] = byte(clientFlags >> 24)

// 	// MaxPacketSize [32 bit] (none)
// 	data[8] = 0x00
// 	data[9] = 0x00
// 	data[10] = 0x00
// 	data[11] = 0x00

// 	// Collation ID [1 byte]
// 	cname := mc.cfg.Collation
// 	if cname == "" {
// 		cname = defaultCollation
// 	}
// 	var found bool
// 	data[12], found = collations[cname]
// 	if !found {
// 		// Note possibility for false negatives:
// 		// could be triggered  although the collation is valid if the
// 		// collations map does not contain entries the server supports.
// 		return fmt.Errorf("unknown collation: %q", cname)
// 	}

// 	// Filler [23 bytes] (all 0x00)
// 	pos := 13
// 	for ; pos < 13+23; pos++ {
// 		data[pos] = 0
// 	}

// 	// SSL Connection Request Packet
// 	// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest
// 	if mc.cfg.TLS != nil {
// 		// Send TLS / SSL request packet
// 		if err := mc.writePacket(data[:(4+4+1+23)+4]); err != nil {
// 			return err
// 		}

// 		// Switch to TLS
// 		tlsConn := tls.Client(mc.netConn, mc.cfg.TLS)
// 		if err := tlsConn.Handshake(); err != nil {
// 			return err
// 		}
// 		mc.rawConn = mc.netConn
// 		mc.netConn = tlsConn
// 		mc.buf.nc = tlsConn
// 	}

// 	// User [null terminated string]
// 	if len(mc.cfg.User) > 0 {
// 		pos += copy(data[pos:], mc.cfg.User)
// 	}
// 	data[pos] = 0x00
// 	pos++

// 	// Auth Data [length encoded integer]
// 	pos += copy(data[pos:], authRespLEI)
// 	pos += copy(data[pos:], authResp)

// 	// Databasename [null terminated string]
// 	if len(mc.cfg.DBName) > 0 {
// 		pos += copy(data[pos:], mc.cfg.DBName)
// 		data[pos] = 0x00
// 		pos++
// 	}

// 	pos += copy(data[pos:], plugin)
// 	data[pos] = 0x00
// 	pos++

// 	// Connection Attributes
// 	data[pos] = byte(len(mc.connector.encodedAttributes))
// 	pos++
// 	pos += copy(data[pos:], []byte(mc.connector.encodedAttributes))

// 	// Send Auth packet
// 	return mc.writePacket(data[:pos])
// }
