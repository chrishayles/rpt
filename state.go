package rpt

type InternalStateChange struct {
	NewState string
}

func newInternalState(state string) *InternalStateChange {
	return &InternalStateChange{
		NewState: state,
	}
}
