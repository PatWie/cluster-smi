package messaging

import (
	"errors"
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
	var rcv_more bool

	m_id, err = soc.RecvBytes(0)
	if err != nil {
		return m, err
	}
	rcv_more, err = soc.GetRcvmore()
	if err != nil {
		return m, err
	}
	if !rcv_more {
		return m, errors.New("Expect m_empty but no more data")
	}
	m_empty, err = soc.RecvBytes(0)
	if err != nil {
		return m, err
	}
	rcv_more, err = soc.GetRcvmore()
	if err != nil {
		return m, err
	}
	if !rcv_more {
		return m, errors.New("Expect m_body but no more data")
	}
	m_body, err = soc.RecvBytes(0)
	if err != nil {
		return m, err
	}

	m = MultipartMessage{m_id, m_empty, m_body}

	return m, nil
}
