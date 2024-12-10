package sdamqp

import (
	"github.com/gaorx/stardust6/sderr"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ChannelConn AMQP channel 和 connection
type ChannelConn struct {
	Chan *amqp.Channel
	Conn *amqp.Connection
}

// Close 关闭
func (cc *ChannelConn) Close() error {
	var chanErr, connErr error
	if cc.Chan != nil {
		chanErr = cc.Chan.Close()
		cc.Chan = nil
	}
	if cc.Conn != nil {
		connErr = cc.Conn.Close()
		cc.Conn = nil
	}
	if chanErr != nil {
		return sderr.Wrapf(chanErr, "close AMQP channel error")
	} else if connErr != nil {
		return sderr.Wrapf(connErr, "close AMQP connection error")
	} else {
		return nil
	}
}
