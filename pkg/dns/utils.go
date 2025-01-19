package dns

import "errors"

var (
	ErrInvalidLabel = errors.New("invalid DNS label")
)

func ValidateDNSName(encoded []byte) error {
	i := 0
	for i < len(encoded) {
		length := int(encoded[i])
		if length == 0 {
			break
		}

		if length < 1 || length > 63 || i+1+length > len(encoded) {
			return ErrInvalidLabel
		}

		label := encoded[i+1 : i+1+length]
		if err := ValidateDNSLabel(label); err != nil {
			return ErrInvalidLabel
		}

		i += 1 + length
	}

	if i != len(encoded)-1 || encoded[len(encoded)-1] != 0 {
		return ErrInvalidLabel
	}

	return nil
}

func ValidateDNSLabel(label []byte) error {
	if len(label) < 1 || len(label) > 63 {
		return ErrInvalidLabel
	}

	// Check if the label starts or ends with a hyphen
	if label[0] == '-' || label[len(label)-1] == '-' {
		return ErrInvalidLabel
	}

	for _, b := range label {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '-') {
			return ErrInvalidLabel
		}
	}

	return nil
}
