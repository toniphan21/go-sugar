package transport

import (
	"fmt"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/sdk"
)

type NodeSerializer interface {
	Serialize() (*sdk.Node, error)
}

var deserializers = map[string]func(in sdk.Node) (sugar.Node, error){}

func RegisterNodeDeserializer(typ string, fn func(in sdk.Node) (sugar.Node, error)) {
	deserializers[typ] = fn
}

func DeserializeNode(in sdk.Node) (sugar.Node, error) {
	fn, have := deserializers[in.Type]
	if !have {
		return nil, fmt.Errorf("unknown node type: %s", in.Type)
	}
	return fn(in)
}

func Registered() []string {
	var out []string
	for k := range deserializers {
		out = append(out, k)
	}
	return out
}
