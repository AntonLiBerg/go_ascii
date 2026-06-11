package aiapi

// AiAPI groups pure helper methods for the go_ascii module.
//
// Keep methods on this type deterministic and side-effect free: no I/O, no
// hidden global state, and no mutation of receiver state.
type AiAPI struct{}

// New returns an AiAPI value for calling pure helper methods.
func New() AiAPI {
	return AiAPI{}
}
