package cryptoengine

import (
	"bytes"
	"errors"
	"github.com/sec51/convert/smallendian"
	"math"
)

// This struct encapsulate the ecnrypted message in a TCP packet, in an easily parseable format
// We assume the data is always encrypted
// Format:
// |version| => 8 bytes (uint64 total message length)
// |type| 	 => 4 bytes (int message version)
// |message| => N bytes ([]byte message)
type message struct {
	Version int    // version of the message, done to support backward compatibility
	Type    int    // message type - this can be ised on the receiver part to process different types
	Text    string // the encrypted message
}

// This struct represent the encrypted message which can be sent over the networl safely
// |lenght| => 8 bytes (uint64 total message length)
// |nonce| => 24 bytes ([]byte size)
// |message| => N bytes ([]byte message)
type EncryptedMessage struct {
	length uint64
	nonce  [nonceSize]byte
	data   []byte
}

// Create a new message with a clear text and the message type
// messageType: is an identifier to distinguish the messages on the receiver and parse them
// for example if zero is a JSON message and 1 is XML, then the received can parse different formats with different methods
func NewMessage(clearText string, messageType int) (message, error) {
	m := message{}
	if clearText == "" {
		return m, errors.New("Clear text cannot be empty")
	}

	m.Text = clearText //:= message{tcpVersion, messageType, clearText}
	m.Type = messageType
	m.Version = tcpVersion
	return m, nil
}

func (m message) toBytes() []byte {
	var buffer bytes.Buffer

	// version
	versionBytes := smallendian.ToInt(m.Version)
	buffer.Write(versionBytes[:])

	// type
	typeBytes := smallendian.ToInt(m.Type)
	buffer.Write(typeBytes[:])

	// message
	buffer.WriteString(m.Text)

	return buffer.Bytes()
}

// Parse the bytes coming from the network and extract
// |length| => 8
// |nonce|	=> nonce size
// |message| => message
func encryptedMessageFromBytes(data []byte) (EncryptedMessage, error) {

	var err error
	var lengthData [8]byte
	var nonceData [nonceSize]byte
	minimumDataSize := 8 + nonceSize
	m := EncryptedMessage{}

	// check if the data is smaller than 36 which is the minimum
	if data == nil {
		return m, MessageParsingError
	}

	if len(data) < minimumDataSize+1 {
		return m, MessageParsingError
	}

	lenght := data[:8]
	nonce := data[8 : 8+nonceSize] // 24 bytes
	message := data[minimumDataSize:]

	total := copy(lengthData[:], lenght)
	if total != 8 {
		return m, MessageParsingError
	}

	total = copy(nonceData[:], nonce)
	if total != nonceSize {
		return m, MessageParsingError
	}

	m.length = smallendian.FromUint64(lengthData)
	m.nonce = nonceData
	m.data = message
	return m, err

}

// This function separates the associated data once decrypted
func messageFromBytes(data []byte) (*message, error) {

	var err error
	var versionData [4]byte
	var typeData [4]byte
	minimumDataSize := 4 + 4
	m := new(message)

	// check if the data is smaller than 36 which is the minimum
	if data == nil {
		return nil, MessageParsingError
	}

	if len(data) < minimumDataSize+1 {
		return nil, MessageParsingError
	}

	version := data[:4]
	typeMsg := data[4:8]
	message := data[8:]

	total := copy(versionData[:], version)
	if total != 4 {
		return nil, MessageParsingError
	}

	total = copy(typeData[:], typeMsg)
	if total != 4 {
		return nil, MessageParsingError
	}

	m.Version = smallendian.FromInt(versionData)
	m.Type = smallendian.FromInt(typeData)
	m.Text = string(message)
	return m, err
}

// STRUCTURE
// 8  => |SIZE|
// 24 => |NONCE|
// N  => |DATA|
// |size| => 8 bytes (uint64 total message length)
// |type| 	 => 4 bytes (int message version)
// |message| => N bytes ([]byte message)
func (m EncryptedMessage) ToBytes() ([]byte, error) {
	if m.length > math.MaxUint64 {
		return nil, errors.New("The message exceeds the maximum allowed sized: uint64 MAX")
	}

	var buffer bytes.Buffer

	// length
	lengthBytes := smallendian.ToUint64(m.length)
	buffer.Write(lengthBytes[:])

	// nonce
	buffer.Write(m.nonce[:])

	// message
	buffer.Write(m.data)

	return buffer.Bytes(), nil

}
