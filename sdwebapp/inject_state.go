package sdwebapp

const akState = "sdwebapp.state"

func State(state any) Attribute {
	return Attr(akState, state)
}
