package protocol

type ERRCODE byte

const (
	ERR_SUCCES ERRCODE = iota
	ERR_INVALID_SIG
)
