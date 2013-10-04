package fiddle

import "bytes"
import "encoding/hex"
import "errors"
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

func Nil () *Bits {
    return &Bits{make([]byte,0), 0}
}

func Zero () *Bits {
    return &Bits{[]byte{0x00}, 1}
}

func One () *Bits {
    return &Bits{[]byte{0x80}, 1}
}

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

func FromChoppedBytes (b []byte) *Bits {
    if len(b) == 0 { panic(errors.New("Decoding error: first byte of chopped bytes must contain chop length")) }
    if len(b) == 1 && b[0] != 0 { panic(errors.New("Decoding error: chop length exceeds content length")) }
    if b[0] > 7 { panic(errors.New("Decoding error: chop length exceeds 8")) }
    return &Bits{b, 8*len(b)-int(b[0])}
}

func FromChunks (chunks []*Bits) *Bits {
    if len(chunks) == 0 { return Nil() }
    b := Nil()
    for i := range chunks[:len(chunks)-1] {
        b = b.Plus(createHeader(chunks[i].len)).Plus(chunks[i])
    }
    return b.Plus(chunks[len(chunks)-1])
}

func FromHex (s string) *Bits {
    b, e := hex.DecodeString(s)
    if e != nil { panic(e) }
    return FromBytes(b)
}

func FromInt (x int) *Bits {
    s := ""
    for d := log2(x); d >= 0; d-- {
        if (x >> uint(d)) % 2 == 0 { s += "0" } else { s += "1" }
    }
    return FromBin(s)
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
        for j := uint(0); 8*i+int(j) < min(bits.len, 8*(i+1)); j++ {
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

func (bits *Bits) PadLeft (length int) *Bits {
    if bits.len > length { return bits }
    return FromBin(strings.Repeat("0", length-bits.len)).Plus(bits)
}

func (bits *Bits) Bin () string {
    b := make([]byte, bits.len)
    for i := 0; i < bits.len; i++ {
        if bits.dat[i/8] & (1 << uint(7-i%8)) == 0 { b[i] = '0' } else { b[i] = '1' }
    }
    return string(b)
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
    return FromBin(bits.Bin()[start:end])
}

/********************
***   Operators   ***
********************/

func (bits *Bits) Equal (other *Bits) bool {
    return bytes.Equal(bits.dat, other.dat) && bits.len == other.len
}

func (bits *Bits) Plus (other *Bits) *Bits {
    return FromBin(bits.Bin() + other.Bin())
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

func (bits *Bits) Int () int {
    s := bits.Bin()
    if s == "" { s = "0" }
    x, e := strconv.ParseInt(s, 2, 64)
    if e != nil { panic(e) }
    return int(x)
}

func (bits *Bits) Unicode () string {
    return string(bits.dat)
}

/***************************
***   Decoding Methods   ***
***************************/

func (bits *Bits) Chunks (num int) (chunks []*Bits, err error) {
    head := 0
    chunks = make([]*Bits, num)
    for i := 0; i < num-1; i++ {
        s, e, err := bits.readHeader(head)
        if err != nil { return nil, err }
        chunks[i] = bits.FromTo(s, e)
        head = e
    }
    chunks[num-1] = bits.From(head)
    return chunks, nil
}

/******************
***   Private   ***
******************/

func min (x int, y int) int {
    if x < y { return x } else { return y }
}

func max (x int, y int) int {
    if x > y { return x } else { return y }
}

func ceil2 (x int) int {
    for y := uint(0); y < 63; y++ {
        if 1 << y >= x { return 1 << y }
    }
    return -1
}

func log2 (x int) int {
    for y := 62; y >= 0; y-- {
        if (x >> uint(y)) % 2 == 1 { return y }
    }
    return -1
}

func numBits (x int) int {
    for y := 62; y >= 0; y-- {
        if (x >> uint(y)) & 1 == 1 { return y+1 }
    }
    return -1
}

func (bits *Bits) readHeader (head int) (start int, end int, err error) {
    if head+4 > bits.len { return 0, 0, errors.New("Decoding error: chunk header index "+strconv.Itoa(head+4)+" out of range") }

    hl := 1 << uint(bits.FromTo(head, head+4).Int()) >> 1
    if head+4+hl > bits.len { return 0, 0, errors.New("Decoding error: chunk start index "+strconv.Itoa(head+4+hl)+" out of range") }

    l := bits.FromTo(head+4, head+4+hl).Int()
    if head+4+hl+l > bits.len { return 0, 0, errors.New("Decoding error: chunk end index "+strconv.Itoa(head+4+hl+l)+" out of range") }

    return head+4+hl, head+4+hl+l, nil
}

func createHeader (length int) *Bits {
    if length < 0 { panic(errors.New("Encoding error: negative length")) }

    // The number of bits needed to encode the length
    headerLength := numBits(length)

    // The number of bits which will actually be used to encode the length
    paddedHeaderLength := ceil2(headerLength)

    // The log-plus-1 encoding (where 0 means 0 length) of the header length
    headerSize := log2(paddedHeaderLength) + 1

    return FromInt(headerSize).PadLeft(4).Plus(FromInt(length).PadLeft(paddedHeaderLength))
}