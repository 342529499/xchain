package protos

import (
	"encoding/json"
)

// Bytes returns this transaction as an array of bytes.
type strArgs struct {
	Function string
	Args     []string
}

// UnmarshalJSON converts the string-based REST/JSON input to
// the []byte-based current ChaincodeInput structure.
func (c *XCodeInput) UnmarshalJSON(b []byte) error {
	sa := &strArgs{}
	err := json.Unmarshal(b, sa)
	if err != nil {
		return err
	}
	allArgs := sa.Args
	if sa.Function != "" {
		allArgs = append([]string{sa.Function}, sa.Args...)
	}
	c.Args = ToXcodeInputArgs(allArgs...)
	return nil
}

func ToXcodeInputArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}
