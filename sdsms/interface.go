package sdsms

import (
	"context"
)

// Interface 短信接口
type Interface interface {
	Send(ctx context.Context, req *SendRequest) error
}

// SendRequest 发送请求
type SendRequest struct {
	TemplateId string
	Messages   Messages
}

// Single 构建一个要发送的短信请求
func Single(templateId string, phone string, signName string, param any) *SendRequest {
	return &SendRequest{
		TemplateId: templateId,
		Messages: Messages{{
			Phone:    phone,
			SignName: signName,
			Param:    param,
		}},
	}
}
