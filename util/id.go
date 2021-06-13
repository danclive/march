package util

import "github.com/danclive/nson-go"

func RandomID() string {
	return nson.NewMessageId().Hex()
}
