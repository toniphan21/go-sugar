package sdk

type PrepareSourceParam struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type PrepareSemanticScopeParam[T any] struct {
	ID            string `json:"id"`
	SemanticScope T      `json:"semantic_scope"`
}

type CleanUpParam struct {
	SourceID string `json:"sourceId"`
	ScopeID  string `json:"scopeId"`
}

type TransformParam struct {
	SourceID string `json:"sourceId"`
	ScopeID  string `json:"scopeId"`
	Node     Node   `json:"node"`
}
