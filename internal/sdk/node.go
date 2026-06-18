package sdk

import "encoding/json"

type Lex struct {
	Tok    int    `json:"tok"`
	Pos    int    `json:"pos"`
	Lit    string `json:"lit"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Offset int    `json:"offset"`
}

type Node struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Pos     Lex             `json:"pos"`
	End     Lex             `json:"end"`
	Payload json.RawMessage `json:"payload"`
}
