package protocol

type ERRCODE byte

const (
	ERR_SUCCES ERRCODE = iota
	ERR_INVALID_SIG
	//ERR_DATA_LEN if `len(data) is not equal to MUGSOFT.DataLen
	ERR_DATA_LEN
	ERR_GAME_TYPE
)
