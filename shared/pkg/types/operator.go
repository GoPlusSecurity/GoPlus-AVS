package types

// Operator 对自身主机的描述
type OperatorSocket struct {
	NodeClass string `json:"node_class"`
	URL       string `json:"url"`
}
