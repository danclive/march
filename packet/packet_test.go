package packet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageId(t *testing.T) {
	packet := NewPacket()

	assert.Exactly(t, packet.Header.MessageId(), uint16(0))

	packet.Header.SetMessageId(12345)

	assert.Exactly(t, packet.Header.MessageId(), uint16(12345))
}

func TestEncodeDecode(t *testing.T) {
	packet := NewPacket()
	packet.Header.SetType(TypeAck)
	packet.Header.SetCode(1)
	packet.Header.SetContentType(1)
	packet.Header.SetExt(1)

	// log.Fatalf("%#v", packet)
	packet.Chan = "haha"
	packet.Body = []byte{1, 2, 3, 4, 5, 6}

	buffer := new(bytes.Buffer)

	err := packet.Encode(buffer)
	if err != nil {
		panic(err)
	}

	packet2, err := Decode(buffer.Bytes())
	if err != nil {
		panic(err)
	}

	assert.Exactly(t, packet, packet2)
}

func TestSetCompressAndSetCrypto(t *testing.T) {
	packet := NewPacket()

	assert.Exactly(t, packet.Header.Bytes()[5], uint8(0))

	packet.Header.SetCompress(CompressZstd)
	assert.Exactly(t, packet.Header.Compress(), CompressZstd)
	assert.Exactly(t, packet.Header.Crypto(), CryptoNone)

	packet.Header.SetCrypto(CryptoNoneAes256Gcm)
	assert.Exactly(t, packet.Header.Compress(), CompressZstd)
	assert.Exactly(t, packet.Header.Crypto(), CryptoNoneAes256Gcm)

	packet.Header.SetCompress(CompressNone)
	assert.Exactly(t, packet.Header.Compress(), CompressNone)
	assert.Exactly(t, packet.Header.Crypto(), CryptoNoneAes256Gcm)

	packet.Header.SetCrypto(CryptoNone)
	assert.Exactly(t, packet.Header.Compress(), CompressNone)
	assert.Exactly(t, packet.Header.Crypto(), CryptoNone)
}
