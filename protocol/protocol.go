package protocol

import (
	"github.com/mugsoft/tools"
	"io"
	"fmt"
)

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
	m__sig_end  = 'g'
	m__sig__len = 2
	//m__meta_len
	m__meta_len = 7
	//m__auth_key_len
	m__auth_key_len = 36
)
const (
	//pos__signature_start pos__ stands for position
	pos__signature_start = iota
	//pos__signature_end position of second signature byte
	pos__signature_end
	//pos__game_id position
	pos__game_id
	//pos__cmd is command position
	pos__cmd = 4
	//pos__data_len defines the position of the byte which has the value how many more bytes are on the way
	pos__data_len = 6
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

func (p *MUGSOFT) Parse(data []byte) ERRCODE {
	if !check__sig(data[:2]) {
		return ERR_INVALID_SIG
	}
	var sig [2]byte
	sig[0] = data[pos__signature_start]
	sig[1] = data[pos__signature_end]
	p.signature = sig
	p.GameType = uint16(tools.LE2Int(data[pos__game_id:pos__cmd]))
	p.CMD = uint16(tools.LE2Int(data[pos__cmd:pos__data_len]))
	p.DataLen = tools.LE2Int(data[pos__data_len : pos__data_len+4])
	p.Data = data[pos__data_len+4:]

	if p.DataLen != uint(len(p.Data)) {
		return ERR_DATA_LEN
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

//Bytes it turns struct to slice of bytes according to `Mugsoft Protocol`
func (p *MUGSOFT) Bytes() []byte {
	var data = make([]byte, m__sig__len)

	data[pos__signature_start] = p.signature[0]
	data[pos__signature_end] = p.signature[1]

	game__type := tools.Int2LE(uint(p.GameType))
	data = append(data, game__type[:2]...)

	cmd := tools.Int2LE(uint(p.CMD))
	data = append(data, cmd[:2]...)

	data__len := tools.Int2LE(p.DataLen)
	data = append(data, data__len[:]...)
	data = append(data, p.Data...)

	return data
}

//Scan it reads from `io.Reader` and fills the struct
func (p *MUGSOFT) Scan(reader io.Reader) error {
	var data []byte
	_, err := reader.Read(data)
	if nil != err {
		return err
	}

	err_code := p.Parse(data)
	if 0 != err_code {
		return fmt.Errorf("Parse error code %v", err_code)
	}

	return nil
}
