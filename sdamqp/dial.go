package sdamqp

import (
	"github.com/gaorx/stardust6/sderr"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Address 地址
type Address struct {
	Url string `json:"url" toml:"url"`
}

// DialConn 获取Conn
func DialConn(addr Address) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr.Url)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return conn, nil
}

// DialChan 获取Conn和Channel
func DialChan(addr Address) (*ChannelConn, error) {
	conn, err := DialConn(addr)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return &ChannelConn{Chan: channel, Conn: conn}, nil
}
