package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

func Sha1(str string) string {
	h := sha1.New()
	bv := []byte(str)
	h.Write(bv)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func RandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func parse_message(message string) (*Message, error) {
	dec := json.NewDecoder(strings.NewReader(message))
	var m *Message
	for {
		if err := dec.Decode(&m); err == io.EOF {
			return nil, err
		} else if err != nil {
			return nil, err
		}

		return m, nil
	}

	return m, nil
}
