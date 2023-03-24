package types

type (
	ValueValidatorFn func(value any) error

	// ParamSetPair is used for associating paramsubspace key and field of param
	// structs.
	ParamSetPair struct {
		Key         []byte
		Value       any
		ValidatorFn ValueValidatorFn
	}

	// ParamSetPairs Slice of KeyFieldPair
	ParamSetPairs []ParamSetPair

	// ParamSet defines an interface for structs containing parameters for a module
	ParamSet interface {
		ParamSetPairs() ParamSetPairs
	}
)

// NewParamSetPair creates a new ParamSetPair instance.
func NewParamSetPair(key []byte, value any, vfn ValueValidatorFn) ParamSetPair {
	return ParamSetPair{key, value, vfn}
}
