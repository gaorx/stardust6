package sderr

const (
	keyCode      = "code"
	keyPublicMsg = "pubmsg"
)

// code

func WithCode(code string) Builder {
	return newBuilder().WithCode(code)
}

func (b Builder) WithCode(code string) Builder {
	return b.With(keyCode, code)
}

func Code(err error) string {
	if err == nil {
		return ""
	}
	if code, ok := Attrs(err)[keyCode].(string); ok {
		return code
	}
	return ""
}

// public msg

func WithPublicMsg(msg string) Builder {
	return newBuilder().WithPublicMsg(msg)
}

func (b Builder) WithPublicMsg(msg string) Builder {
	return b.With(keyPublicMsg, msg)
}

func PublicMsg(err error) string {
	if err == nil {
		return ""
	}
	unwrappedErrs := UnwrapNested(err)
	for _, unwrappedErr := range unwrappedErrs {
		if e, ok := Probe(unwrappedErr); ok {
			if pubMsg, ok := e.attrs[keyPublicMsg].(string); ok {
				if pubMsg != "" {
					return pubMsg
				}
			}
		}
	}
	return err.Error()
}
