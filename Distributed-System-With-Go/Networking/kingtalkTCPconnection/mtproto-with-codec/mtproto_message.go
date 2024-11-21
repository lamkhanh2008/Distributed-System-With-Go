package mtproto

import "fmt"

type MessageBase interface {
	Encode() []byte
	Decode(b []byte) error
}

func NewMTPRawMessage(authKeyId int64, quickAckId int32, connType int) *MTPRawMessage {
	return &MTPRawMessage{
		connType:   connType,
		authKeyId:  authKeyId,
		quickAckId: quickAckId,
	}
}

type MTPRawMessage struct {
	connType   int
	authKeyId  int64 // obtained by decompressing the original data
	quickAckId int32 // EncryptedMessage, there may be

	// Raw data
	Payload []byte
}

func (m *MTPRawMessage) String() string {
	return fmt.Sprintf("{conn_type: %d, auth_key_id: %d, quick_ack_id: %d, payload_len: %d}",
		m.connType,
		m.authKeyId,
		m.quickAckId,
		len(m.Payload))
}

func (m *MTPRawMessage) ConnType() int {
	return m.connType
}

func (m *MTPRawMessage) AuthKeyId() int64 {
	return m.authKeyId
}

func (m *MTPRawMessage) QuickAckId() int32 {
	return m.quickAckId
}

////////////////////////////////////////////////////////////////////////////
func (m *MTPRawMessage) Encode() []byte {
	return m.Payload
}

func (m *MTPRawMessage) Decode(b []byte) error {
	m.Payload = b
	return nil
}
