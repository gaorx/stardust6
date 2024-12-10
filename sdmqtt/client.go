package sdmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gaorx/stardust6/sderr"
)

// Client MQTT客户端
type Client struct {
	mqtt.Client
}

// ClientOptions MQTT客户端选项
type ClientOptions = mqtt.ClientOptions

// Dial 连接MQTT服务器
func Dial(opts *ClientOptions) *Client {
	c := mqtt.NewClient(opts)
	return &Client{c}
}

// ConnectSync 连接MQTT服务器(同步方式)
func (c *Client) ConnectSync() error {
	token := c.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrapf(err, "connect MQTT sync error")
	}
	return nil
}

// SubscribeSync 订阅主题(同步方式)
func (c *Client) SubscribeSync(topic string, qos byte, callback mqtt.MessageHandler) error {
	token := c.Subscribe(topic, qos, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrapf(err, "subscribe MQTT sync error")
	}
	return nil
}

// UnsubscribeSync 取消订阅主题(同步方式)
func (c *Client) UnsubscribeSync(topics ...string) error {
	token := c.Unsubscribe(topics...)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrapf(err, "unsubscribe MQTT sync error")
	}
	return nil
}

// PublishSync 发布消息(同步方式)
func (c *Client) PublishSync(topic string, qos byte, retained bool, payload any) error {
	token := c.Publish(topic, qos, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrapf(err, "publish MQTT sync error")
	}
	return nil
}
