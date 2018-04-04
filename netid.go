package lorawan

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
)

// NetID represents the NetID.
type NetID [3]byte

// NwkID returns the NwkID bits of the NetID.
func (n NetID) NwkID() byte {
	return n[2] & 127 // 7 lsb
}

// String implements fmt.Stringer.
func (n NetID) String() string {
	return hex.EncodeToString(n[:])
}

// MarshalText implements encoding.TextMarshaler.
func (n NetID) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *NetID) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(b) != len(n) {
		return fmt.Errorf("lorawan: exactly %d bytes are expected", len(n))
	}
	copy(n[:], b)
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (n NetID) MarshalBinary() ([]byte, error) {
	out := make([]byte, len(n))

	// little endian
	for i, v := range n {
		out[len(n)-1-i] = v
	}

	return out, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (n *NetID) UnmarshalBinary(data []byte) error {
	if len(data) != len(n) {
		return fmt.Errorf("lorawan: %d bytes of data are expected", len(n))
	}

	for i, v := range data {
		// little endian
		n[len(n)-1-i] = v
	}

	return nil
}

// Value implements driver.Valuer.
func (n NetID) Value() (driver.Value, error) {
	return n[:], nil
}

// Scan implements sql.Scanner.
func (n *NetID) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("lorawan: []byte type expected")
	}
	if len(b) != len(n) {
		return fmt.Errorf("lorawan: []byte must have length %d", len(n))
	}
	copy(n[:], b)
	return nil
}
