package sdsms

import (
	"encoding/json"
	"github.com/samber/lo"
)

// Message 短信消息
type Message struct {
	Phone    string
	SignName string
	Param    any
	Content  string
}

// Messages 多个短信消息
type Messages []Message

// ParamToMap 将此message的param转换为map形式
func (msg Message) ParamToMap() (map[string]string, bool) {
	param := msg.Param
	if param == nil {
		return nil, true
	} else if p, ok := param.(map[string]string); ok {
		return p, true
	} else {
		j, err := json.Marshal(param)
		if err != nil {

			return nil, false
		}
		var paramMap map[string]string
		err = json.Unmarshal(j, &paramMap)
		if err != nil {
			return nil, false
		}
		return paramMap, true
	}
}

// SetSignName 设置所有消息的签名
func (msgs Messages) SetSignName(signName string) {
	for i := 0; i < len(msgs); i++ {
		msgs[i].SignName = signName
	}
}

// Phones 获取所有消息的手机号
func (msgs Messages) Phones() []string {
	return lo.Map(msgs, func(msg Message, _ int) string {
		return msg.Phone
	})
}

// SignNames 获取所有消息的签名
func (msgs Messages) SignNames() []string {
	return lo.Map(msgs, func(msg Message, _ int) string {
		return msg.SignName
	})
}

// Params 获取所有消息的参数
func (msgs Messages) Params() []any {
	return lo.Map(msgs, func(msg Message, _ int) any {
		return msg.Param
	})
}

// ParamsToMap 将所有消息的param转换为map形式
func (msgs Messages) ParamsToMap() ([]map[string]string, bool) {
	var params []map[string]string
	for _, msg := range msgs {
		param, ok := msg.ParamToMap()
		if !ok {
			return nil, false
		}
		params = append(params, param)
	}
	return params, true
}
