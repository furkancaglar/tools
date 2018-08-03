package protocol

import (
	"github.com/mugsoft/tools"
	"io"
	"fmt"
)

type (
	//GameID is constant that defines the games
	GAME_TYPE uint16

	COMMAND uint16
)

type MUGSOFT struct {
	signature [2]byte
	GameType  GAME_TYPE
	CMD       COMMAND
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
	pos__cmd = 4
	//pos__data_len defines the position of the byte which has the value how many more bytes are on the way
	pos__data_len = 6
)
const (
	//CMD_ERROR is the command 0 which stads for error
	CMD_ERROR COMMAND = iota
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
const (
	_ GAME_TYPE = iota
	GAME_TOMBALA
	GAME_KENO
)

var game__map = map[GAME_TYPE]string{
	GAME_TOMBALA: "TOMBALA",
	GAME_KENO:    "KENO",
}

var command__map = map[COMMAND]string{
	CMD_ERROR:     "ERROR",
	CMD_HANDSHAKE: "HANDSHAKE",
	CMD_KEYCHECK:  "KEY CHECK",
	CMD_NEWGAME:   "NEW GAME",
	CMD_NEWBALL:   "NEW BALL",
	CMD_WINNING:   "WINNING",
	CMD_ENDGAME:   "END GAME",
	CMD_CARD:      "CARD",
	CMD_PING:      "PING",
}

func (p *MUGSOFT) Parse(data []byte) ERRCODE {
	if !check__sig(data[:2]) {
		return ERR_INVALID_SIG
	}
	var sig [2]byte
	sig[0] = data[pos__signature_start]
	sig[1] = data[pos__signature_end]
	p.signature = sig

	var game__type = GAME_TYPE(tools.LE2Int(data[pos__game_id:pos__cmd]))
	if !check__game__type(game__type) {
		return ERR_GAME_TYPE
	}

	p.GameType = game__type

	var cmd = COMMAND(tools.LE2Int(data[pos__cmd:pos__data_len]))
	if !check__cmd(cmd) {
		return ERR_COMMAND
	}

	p.CMD = cmd
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

func check__game__type(game__type GAME_TYPE) bool {
	_, ok := game__map[game__type]
	if !ok {
		return false
	}

	return true
}

func check__cmd(cmd COMMAND) bool {
	_, ok := command__map[cmd]
	if !ok {
		return false
	}
	return true
}

//Bytes it turns struct to slice of bytes according to `Mugsoft Protocol`
func (p *MUGSOFT) Bytes() []byte {
	p.signature = [2]byte{'m', 'g'}
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

	total__data = append(total__data, meta...)

	var remainning__data__len = int(tools.LE2Int(meta[pos__data_len : pos__data_len+4]))
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
	if 0 != err__code {
		return fmt.Errorf("parse error code %v", err__code)
	}

	return nil
}
