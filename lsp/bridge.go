package lsp

func DoNothing(envelope Envelope) (Envelope, error) {

	return envelope, nil
}

var _ EnvelopeHook = DoNothing
