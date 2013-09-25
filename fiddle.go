package fiddle

import "bytes"
import "encoding/hex"
import "strconv"
import "strings"

/********************
***   Bits Type   ***
********************/

type Bits struct {
    dat []byte
    len int
}

/***********************
***   Constructors   ***
***********************/

func FromByte (b byte) *Bits {
    return &Bits{[]byte{b}, 8}
}

func FromBytes (b []byte) *Bits {
    return &Bits{b, 8*len(b)}
}

func FromBin (s string) *Bits {
    s = strings.Replace(s, " ", "", -1)
    l := (len(s)+7) / 8
    b := &Bits{make([]byte, l), len(s)}
    for i := 0; i < len(s); i++ {
        if s[i] == '1' { b.dat[i/8] |= 1 << uint(7-i%8) }
    }
    return b
}

func FromHex (s string) *Bits {
    b, e := hex.DecodeString(s)
    if e != nil { panic(e) }
    return FromBytes(b)
}

func FromUnicode (s string) *Bits {
    b := []byte(s)
    return &Bits{b, 8*len(b)}
}

/*************************
***   Common Methods   ***
*************************/

func (bits *Bits) Len () int {
    return bits.len
}

func (bits *Bits) String () string {
    s := ""
    for i := 0; i < len(bits.dat); i++ {
        if i > 0 { s += " " }
        for j := uint(0); i*8+int(j)<bits.len; j++ {
            if bits.dat[i] >> (7-j) & 1 == 0 { s += "0" } else { s += "1" }
        }
    }
    return s
}

func (bits *Bits) HexString () string {
    chop := "-" + strconv.Itoa(8 - bits.len%8)
    if chop == "-8" { chop = "" }
    return hex.EncodeToString(bits.dat) + chop
}

/*************************
***   Splice Methods   ***
*************************/

func (bits *Bits) To (end int) *Bits {
    return bits.FromTo(0, end)
}

func (bits *Bits) From (start int) *Bits {
    return bits.FromTo(start, bits.len)
}

func (bits *Bits) FromTo (start int, end int) *Bits {
    // Byte splicing
    start = min(max(start, 0    ), bits.len)
    end   = min(max(end  , start), bits.len)
    b := &Bits{bits.dat[start/8:(end+7)/8], end-start}

    // Bit shifting
    shift := uint(start % 8)
    if shift > 0 {
        for i := 0; i < len(b.dat)-1; i++ {
            b.dat[i] = (b.dat[i] << shift) | (b.dat[i+1] >> uint(8-shift))
        }
        b.dat[len(b.dat)-1] = b.dat[len(b.dat)-1] << shift
    }

    // Bit chopping
    chop := uint(8 - (end-start) % 8)
    if chop > 0 {
        b.dat[len(b.dat)-1] = b.dat[len(b.dat)-1] & (byte(0xFF) << chop)
    }

    return b
}

/********************
***   Operators   ***
********************/

func (bits *Bits) Equal (other *Bits) bool {
    return bytes.Equal(bits.dat, other.dat) && bits.len == other.len
}

func (bits *Bits) Plus (other *Bits) *Bits {
    // // Byte splicing
    // b := &Bits{append(bits.dat, other.dat...), bits.len + other.len}

    // // Bit shifting
    // shift := uint(start % 8)
    // if shift > 0 {
    //     for i := len(); i < len(b.dat)-1; i++ {
    //         b.dat[i] = (b.dat[i] << shift) | (b.dat[i+1] >> uint(8-shift))
    //     }
    //     b.dat[len(b.dat)-1] = b.dat[len(b.dat)-1] << shift
    // }
    return &Bits{append(bits.dat, other.dat...), bits.len + other.len}
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

/***************************
***   Decoding Methods   ***
***************************/

// func (bits *Bits) Chunks (num int) chunks []*Bits, err error {
    
// }

/******************
***   Private   ***
******************/

func min (x int, y int) int {
    if x < y { return x } else { return y }
}

func max (x int, y int) int {
    if x > y { return x } else { return y }
}

func parseBinByte (b string) byte {
    x, e := strconv.ParseUint(b, 2, 8)
    if e != nil { panic(e) }
    return byte(x)
}