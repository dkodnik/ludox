package patch

import (
	"errors"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// EOF is the end of the IPS patch
const EOF = 0x454f46

func ipsAllocTargetData(patch, source []byte) ([]byte, error) {
	offset := 5
	targetLength := len(source)

	for {
		if offset > len(patch)-3 {
			break
		}

		address := int(patch[offset]) << 16
		offset++
		address |= int(patch[offset]) << 8
		offset++
		address |= int(patch[offset]) << 0
		offset++

		if address == EOF {
			if offset == len(patch) {
				prov := make([]byte, targetLength)
				return prov, nil
			} else if offset == len(patch)-3 {
				size := int(patch[offset]) << 16
				offset++
				size |= int(patch[offset]) << 8
				offset++
				size |= int(patch[offset]) << 0
				offset++
				targetLength = size
				prov := make([]byte, targetLength)
				return prov, nil
			}
		}

		if offset > len(patch)-2 {
			break
		}

		length := int(patch[offset]) << 8
		offset++
		length |= int(patch[offset]) << 0
		offset++

		if length > 0 /* Copy */ {
			if offset > len(patch)-int(length) {
				break
			}

			for ; length > 0; length-- {
				address++
				offset++
			}
		} else /* RLE */ {
			if offset > len(patch)-3 {
				break
			}

			length := int(patch[offset]) << 8
			offset++
			length |= int(patch[offset]) << 0
			offset++

			if length == 0 /* Illegal */ {
				break
			}

			for ; length > 0; length-- {
				address++
			}

			offset++
		}

		if address > targetLength {
			targetLength = address
		}
	}

	txtI18n := l10n.T9(&i18n.Message{ID: "InvalidPatch", Other: "invalid patch"})
	return nil, errors.New(txtI18n)
}

func applyIPS(patch, source []byte) (*[]byte, error) {
	if len(patch) < 8 {
		txtI18n := l10n.T9(&i18n.Message{ID: "PatchTooSmall", Other: "patch too small"})
		return nil, errors.New(txtI18n)
	}

	if string(patch[0:5]) != "PATCH" {
		txtI18n := l10n.T9(&i18n.Message{ID: "InvalidPatchHeader", Other: "invalid patch header"})
		return nil, errors.New(txtI18n)
	}

	targetData, err := ipsAllocTargetData(patch, source)
	if err != nil {
		return nil, err
	}

	copy(targetData, source)

	offset := 5
	for {
		if offset > len(patch)-3 {
			break
		}

		address := int(patch[offset]) << 16
		offset++
		address |= int(patch[offset]) << 8
		offset++
		address |= int(patch[offset]) << 0
		offset++

		if address == EOF {
			if offset == len(patch) {
				return &targetData, nil
			} else if offset == len(patch)-3 {
				size := int(patch[offset]) << 16
				offset++
				size |= int(patch[offset]) << 8
				offset++
				size |= int(patch[offset]) << 0
				offset++
				return &targetData, nil
			}
		}

		if offset > len(patch)-2 {
			break
		}

		length := int(patch[offset]) << 8
		offset++
		length |= int(patch[offset]) << 0
		offset++

		if length > 0 /* Copy */ {
			if offset > len(patch)-length {
				break
			}

			for ; length > 0; length-- {
				targetData[address] = patch[offset]
				address++
				offset++
			}
		} else /* RLE */ {
			if offset > len(patch)-3 {
				break
			}

			length = int(patch[offset]) << 8
			offset++
			length |= int(patch[offset]) << 0
			offset++

			if length == 0 /* Illegal */ {
				break
			}

			for ; length > 0; length-- {
				targetData[address] = patch[offset]
				address++
			}

			offset++
		}
	}

	txtI18n := l10n.T9(&i18n.Message{ID: "InvalidPatch", Other: "invalid patch"})
	return nil, errors.New(txtI18n)
}
