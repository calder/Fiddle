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
    if len(b) == 0 { return Nil() }
    chop := (5 + FromByte(b[0]).To(3).Int()) % 8
    return FromRawBytes(b).FromTo(3, 8*len(b)-chop)
}

func FromRawBytes (b []byte) *Bits {
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

func FromChunks (chunks ...*Bits) *Bits {
    if len(chunks) == 0 { return Nil() }
    b := Nil()
    for i := range chunks[:len(chunks)-1] {
        b = b.Plus(createHeader(chunks[i].len)).Plus(chunks[i])
    }
    return b.Plus(chunks[len(chunks)-1])
}

func FromList (list []*Bits) *Bits {
    if len(list) == 0 { return Nil() }
    b := Nil()
    for i := range list[:len(list)] {
        b = b.Plus(createHeader(list[i].len)).Plus(list[i])
    }
    return b
}

func FromHex (s string) *Bits {
    b, e := hex.DecodeString(s)
    if e != nil { panic(e) }
    return FromBytes(b)
}

func FromRawHex (s string) *Bits {
    b, e := hex.DecodeString(s)
    if e != nil { panic(e) }
    return FromRawBytes(b)
}

func FromInt (x int) *Bits {
    s := ""
    for d := numBits(x)-1; d >= 0; d-- {
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

func (bits *Bits) PadLeft (length int) *Bits {
    if bits.len > length { return bits }
    return FromBin(strings.Repeat("0", length-bits.len)).Plus(bits)
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

func (bits *Bits) Bin () string {
    b := make([]byte, bits.len)
    for i := 0; i < bits.len; i++ {
        if bits.dat[i/8] & (1 << uint(7-i%8)) == 0 { b[i] = '0' } else { b[i] = '1' }
    }
    return string(b)
}

func (bits *Bits) Bytes () []byte {
    return FromInt(invRemainder(bits.len, 8)).PadLeft(3).Plus(bits).dat
}

func (bits *Bits) RawBytes () []byte {
    return bits.dat
}

func (bits *Bits) Hex () string {
    return hex.EncodeToString(bits.Bytes())
}

func (bits *Bits) RawHex () string {
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

func (bits *Bits) Chunks (num int) []*Bits {
    head := 0
    chunks := make([]*Bits, num)
    for i := 0; i < num-1; i++ {
        s, e, err := bits.readHeader(head)
        if err != nil { panic(err) }
        chunks[i] = bits.FromTo(s, e)
        head = e
    }
    chunks[num-1] = bits.From(head)
    return chunks
}

func (bits *Bits) List () (list []*Bits, err error) {
    list = []*Bits{}
    for head := 0; head < bits.len; {
        s, e, err := bits.readHeader(head)
        if err != nil { return nil, err }
        list = append(list, bits.FromTo(s, e))
        head = e
    }
    return list, nil
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

func invRemainder (x int, y int) int {
    return (y - x%y) % y
}

func numBits (x int) int {
    for y := 62; y >= 0; y-- {
        if (x >> uint(y)) & 1 == 1 { return y+1 }
    }
    return 0
}

func (bits *Bits) readHeader (head int) (start int, end int, err error) {
    if head+8 > bits.len { return 0, 0, errors.New("Decoding error: chunk header index "+strconv.Itoa(head+8)+" out of range") }

    hl := bits.FromTo(head, head+8).Int()
    if head+8+hl > bits.len { return 0, 0, errors.New("Decoding error: chunk start index "+strconv.Itoa(head+8+hl)+" out of range") }

    l := bits.FromTo(head+8, head+8+hl).Int()
    if head+8+hl+l > bits.len { return 0, 0, errors.New("Decoding error: chunk end index "+strconv.Itoa(head+8+hl+l)+" out of range") }

    return head+8+hl, head+8+hl+l, nil
}

func createHeader (length int) *Bits {
    if length < 0 { panic(errors.New("Encoding error: cannot create header with negative length")) }
    return FromInt(numBits(length)).PadLeft(8).Plus(FromInt(length))
}