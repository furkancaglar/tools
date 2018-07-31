package protocol

type MUGSOFT struct {
	signature [2]byte
	GameType  uint16
	CMD       uint16
	Data      []byte
	DataLen   uint
}

const (
	//m__sig_start m__ stands for meta
	m__sig_start = 'm'
	//m__sig_end
	m__sig_end = 'g'
	//m__meta_len
	m__meta_len = 5
	//m__auth_key_len
	m__auth_key_len = 36
	//pos__signature_start pos__ stands for position
	pos__signature_start = iota
	//pos__signature_end position of second signature byte
	pos__signature_end
	//pos__game_id position
	pos__game_id
	//pos__cmd is command position
	pos__cmd
	//pos__data_len defines the position of the byte which has the value how many more bytes are on the way
	pos__data_len
)
const (
	//CMD_ERROR is the command 0 which stads for error
	CMD_ERROR uint16 = iota
	//CMD_HANDSHAKE handshake command
	CMD_HANDSHAKE
	//CMD_KEYCHECK checks the key sent after handshake
	CMD_KEYCHECK
	//CMD_NEWGAME
	CMD_NEWGAME
	//CMD_NEWBALL
	CMD_NEWBALL
	//CMD_WINNING
	CMD_WINNING
	//CMD_ENDGAME
	CMD_ENDGAME
)

func (p *MUGSOFT) Unmarshal(data []byte) ERRCODE {
	if !check__sig(data[:2]) {
		return ERR_INVALID_SIG
	}
	return ERR_SUCCES
}
func check__sig(sig []byte) (isSigCorrect bool) {
	isSigCorrect = true
	if len(sig) < 2 || m__sig_start != sig[0] || m__sig_end != sig[1] {
		isSigCorrect = false
	}
	return
}
