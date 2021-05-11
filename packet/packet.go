package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Packet struct {
	Header Header
	Chan   string
	Body   []byte
}

type Header struct {
	bytes []byte // 8
}

const Version byte = 1

const (
	TypeNon byte = iota
	TypeCon
	TypeAck
	TypeRst
)

const (
	CodeOK byte = iota
)

const (
	ContentTypeNone byte = iota
	ContentTypeText
	ContentTypeJson
	ContentTypeNson
	ContentTypeBin
)

const (
	CompressNone byte = iota
	CompressZstd
	CompressGzip
)

const (
	CryptoNone byte = iota
	CryptoNoneAes128Gcm
	CryptoNoneAes256Gcm
	CryptoNoneChaCha20Poly1305
)

func NewPacket() *Packet {
	packet := &Packet{
		Header: NewHeader(),
		Body:   make([]byte, 0),
	}

	return packet
}

func NewHeader() Header {
	header := Header{
		bytes: make([]byte, 8),
	}

	header.bytes[2] = Version

	return header
}

func NewHeaderFromBytes(bytes []byte) (Header, error) {
	header := Header{
		bytes,
	}

	if header.bytes[2] > Version {
		return header, fmt.Errorf("version %v unsupported", header.bytes[2])
	}

	if header.Type() > TypeRst {
		return header, fmt.Errorf("type %v unsupported", header.Type())
	}

	if header.Compress() > CompressGzip {
		return header, fmt.Errorf("compress %v unsupported", header.Compress())
	}

	if header.Crypto() > CryptoNoneChaCha20Poly1305 {
		return header, fmt.Errorf("crypto %v unsupported", header.Crypto())
	}

	if header.ContentType() > ContentTypeBin {
		return header, fmt.Errorf("content type %v unsupported", header.ContentType())
	}

	return header, nil
}

func (h *Header) Bytes() []byte {
	return h.bytes
}

func (h *Header) MessageId() uint16 {
	return binary.LittleEndian.Uint16(h.bytes[:2])
}

func (h *Header) SetMessageId(messageId uint16) {
	binary.LittleEndian.PutUint16(h.bytes, messageId)
}

func (h *Header) Type() byte {
	return h.bytes[3]
}

func (h *Header) SetType(t byte) error {
	if t > TypeRst {
		return fmt.Errorf("type %v unsupported", t)
	}

	h.bytes[3] = t

	return nil
}

func (h *Header) Code() byte {
	return h.bytes[4]
}

func (h *Header) SetCode(c byte) {
	h.bytes[4] = c
}

func (h *Header) Compress() byte {
	c := h.bytes[5]
	return c >> 4
}

func (h *Header) SetCompress(c byte) error {
	if c > CompressGzip {
		return fmt.Errorf("compress %v unsupported", c)
	}

	h.bytes[5] &= 0b00001111
	h.bytes[5] |= c << 4

	return nil
}

func (h *Header) Crypto() byte {
	c := h.bytes[5]
	c &= 0b00001111
	return c
}

func (h *Header) SetCrypto(c byte) error {
	if c > CryptoNoneChaCha20Poly1305 {
		return fmt.Errorf("crypto %v unsupported", c)
	}

	h.bytes[5] &= 0b11110000
	h.bytes[5] |= c

	return nil
}

func (h *Header) ContentType() byte {
	return h.bytes[6]
}

func (h *Header) SetContentType(c byte) error {
	if c > ContentTypeBin {
		return fmt.Errorf("content type %v unsupported", c)
	}

	h.bytes[6] = c

	return nil
}

func (h *Header) Ext() byte {
	return h.bytes[7]
}

func (h *Header) SetExt(ext byte) {
	h.bytes[7] = ext
}

func (h *Header) Copy() Header {
	header := Header{
		bytes: make([]byte, 8),
	}

	copy(header.bytes, h.bytes)

	return header
}

func Decode(bytes []byte) (*Packet, error) {
	if len(bytes) < 9 {
		return nil, errors.New("The bytes are too short")
	}

	header, err := NewHeaderFromBytes(bytes[:8])
	if err != nil {
		return nil, err
	}

	n := 8
	for {
		if bytes[n] != 0 {
			n += 1
		} else {
			break
		}
	}

	packet := &Packet{
		Header: header,
		Chan:   string(bytes[8:n]),
		Body:   bytes[n+1:],
	}

	return packet, nil
}

func (p *Packet) Encode(buffer *bytes.Buffer) error {
	buffer.Write(p.Header.bytes)
	buffer.WriteString(p.Chan)
	buffer.WriteByte(0)

	buffer.Write(p.Body)

	return nil
}
