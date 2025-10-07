package gob

import (
	"bytes"
	"encoding/gob"
	"time"
)

func encode(gl GameLog) ([]byte, error) {
	w := &bytes.Buffer{}
	if err := gob.NewEncoder(w).Encode(gl); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func decode(data []byte) (GameLog, error) {
	gl := GameLog{}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&gl); err != nil {
		return gl, err
	}

	return gl, nil
}

type GameLog struct {
	CurrentTime time.Time
	Message     string
	Username    string
}
