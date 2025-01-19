package udpserver

import "bytes"

type Serializable interface {
	Serialize() ([]byte, error)
}

func SerializeItems(items []Serializable) ([]byte, error) {

	buf := bytes.Buffer{}

	for _, item := range items {
		serialized, err := item.Serialize()
		if err != nil {
			return nil, err
		}
		buf.Write(serialized)
	}

	return buf.Bytes(), nil
}
