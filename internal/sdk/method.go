package sdk

var Method = struct {
	ID                   string
	Binary               string
	Parse                string
	Restore              string
	PrepareSource        string
	PrepareSemanticScope string
	CleanUp              string
	StructuralTransform  string
	SemanticTransformer  string
	RestoreTransform     string
}{
	ID:                   "ID",
	Binary:               "Binary",
	Parse:                "Parse",
	Restore:              "Restore",
	PrepareSource:        "PrepareSource",
	PrepareSemanticScope: "PrepareSemanticScope",
	CleanUp:              "CleanUp",
	StructuralTransform:  "StructuralTransform",
	SemanticTransformer:  "SemanticTransformer",
	RestoreTransform:     "RestoreTransform",
}
