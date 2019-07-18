package protocol

import (
	"github.com/furkancaglar/tools"
	"io"
	"fmt"
	"unsafe"
)

type MUGSOFT struct {
	signature [2]byte
	Type      uint16
	CMD       uint16
	Data      []byte
	DataLen   uint32
}

const (
	//m__sig_start m__ stands for meta
	m__sig_start = 'm'
	//m__sig_end
	m__sig_end  = 'g'
	m__sig__len = 2
	//m__meta_len
	m__meta_len = 10
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
	pos__cmd = iota + 1
	_
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
	//CMD_CARD
	CMD_CARD
	//CMD_PING
	CMD_PING
)

func (p *MUGSOFT) Parse(data []byte) ERRCODE {
	if 0 == len(data) || nil == data {
		return ERR_NIL_DATA
	}
	if m__meta_len > len(data) {
		return ERR_META_LEN
	}
	if !check__sig(data[:2]) {
		return ERR_INVALID_SIG
	}
	var sig [2]byte
	sig[0] = data[pos__signature_start]
	sig[1] = data[pos__signature_end]
	p.signature = sig
	p.Type = uint16(tools.LE2Int(data[pos__game_id:pos__cmd]))
	p.CMD = uint16(tools.LE2Int(data[pos__cmd:pos__data_len]))
	p.DataLen = tools.LE2Int(data[pos__data_len : pos__data_len+unsafe.Sizeof(p.DataLen)])
	p.Data = data[pos__data_len+unsafe.Sizeof(p.DataLen):]

	if p.DataLen != uint32(len(p.Data)) {
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
	p.signature = [2]byte{'m', 'g'}
	var data = make([]byte, m__sig__len)

	data[pos__signature_start] = p.signature[0]
	data[pos__signature_end] = p.signature[1]

	game__type := tools.Int2LE(uint32(p.Type))
	data = append(data, game__type[:unsafe.Sizeof(p.Type)]...)

	cmd := tools.Int2LE(uint32(p.CMD))
	data = append(data, cmd[:unsafe.Sizeof(p.CMD)]...)

	data__len := tools.Int2LE(p.DataLen)
	data = append(data, data__len[:]...)
	data = append(data, p.Data...)

	return data
}

//Scan it reads from `io.Reader` and fills the struct
func (p *MUGSOFT) Scan(reader io.Reader) error {
	var total__data []byte
	var meta = make([]byte, m__meta_len)
	n, err := reader.Read(meta)
	if nil != err {
		return err
	}
	if n != m__meta_len {
		return fmt.Errorf("meta length is not rigth")
	}

	if !check__sig(meta[:2]) {
		return fmt.Errorf("signature error")
	}

	total__data = meta

	var remainning__data__len = int(tools.LE2Int(meta[pos__data_len : pos__data_len+unsafe.Sizeof(p.DataLen)]))
	var remainning__data = make([]byte, remainning__data__len)

consume__remaining:

	n, err = reader.Read(remainning__data[:remainning__data__len])
	remainning__data__len -= n
	if nil != err {
		return err
	}
	total__data = append(total__data, remainning__data[:n]...)
	if remainning__data__len > 0 {

		goto consume__remaining
	}

	err__code := p.Parse(total__data)
	if ERR_SUCCES != err__code {
		return fmt.Errorf("parse error code %v", err__code)
	}

	return nil
}
