package model

type KeyInfo struct {
	DataType string
	DataSize int
}

type KeyInfoer interface {
	KeyInfo(key string) (KeyInfo, error)
}
