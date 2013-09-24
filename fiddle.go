package fiddle

import "bytes"
import "encoding/hex"

/********************
***   Bits Type   ***
********************/

type Bits struct {
    dat    []byte
    length int
}

/***********************
***   Constructors   ***
***********************/

func FromByte (dat byte) *Bits {
    return &Bits{[]byte{dat}, 8}
}

func FromBytes (dat []byte) *Bits {
    return &Bits{dat, 8*len(dat)}
}

func FromUnicode (dat string) *Bits {
    b := []byte(dat)
    return &Bits{b, 8*len(b)}
}

/*************************
***   Common Methods   ***
*************************/

func (bits *Bits) Len () int {
    return bits.length
}

/*************************
***   Splice Methods   ***
*************************/

func (bits *Bits) To (start int) *Bits {
    return &Bits{bits.dat[:start/8], bits.length-start}
}

func (bits *Bits) From (start int) *Bits {
    return &Bits{bits.dat[start/8:], bits.length-start}
}

func (bits *Bits) FromTo (start int, end int) *Bits {
    return &Bits{bits.dat[start/8:end/8], end-start}
}

/********************
***   Operators   ***
********************/

func (bits *Bits) Equal (other *Bits) bool {
    return bytes.Equal(bits.dat, other.dat)
}

func (bits *Bits) Plus (other *Bits) *Bits {
    return &Bits{append(bits.dat, other.dat...), bits.length + other.length}
}

/*****************************
***   Conversion Methods   ***
*****************************/

func (bits *Bits) Byte () byte {
    return bits.dat[0]
}

func (bits *Bits) Bytes () []byte {
    return bits.dat
}

func (bits *Bits) Hex () string {
    return hex.EncodeToString(bits.dat)
}

func (bits *Bits) Unicode () string {
    return string(bits.dat)
}