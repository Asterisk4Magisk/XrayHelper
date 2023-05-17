package sharelinks

// ShareLink implement this interface, that node can be converted to xray OutoundJsonObject
type ShareLink interface {
	GetNodeInfo() string
	ToOutoundJsonWithTag(tag string) string
}
