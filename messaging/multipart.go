package messaging

import (
	"github.com/pebbe/zmq4"
)

type MultipartMessage struct {
	Id    []byte
	Empty []byte
	Body  []byte
}

func SendMultipartMessage(soc *zmq4.Socket, m *MultipartMessage) (err error) {
	_, err = soc.SendBytes(m.Id, zmq4.SNDMORE)
	if err != nil {
		return err
	}
	_, err = soc.SendBytes(m.Empty, zmq4.SNDMORE)
	if err != nil {
		return err
	}
	_, err = soc.SendBytes(m.Body, 0)
	if err != nil {
		return err
	}
	return nil
}

func ReceiveMultipartMessage(soc *zmq4.Socket) (m MultipartMessage, err error) {

	var m_id []byte
	var m_empty []byte
	var m_body []byte

	m_id, err = soc.RecvBytes(0)
	if err != nil {
		return m, err
	}
	m_empty, err = soc.RecvBytes(0)
	if err != nil {
		return m, err
	}
	m_body, err = soc.RecvBytes(0)
	if err != nil {
		return m, err
	}

	m = MultipartMessage{m_id, m_empty, m_body}

	return m, nil
}
