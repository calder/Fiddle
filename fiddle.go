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

func (bits *Bits) RawString () string {
    return strings.Replace(bits.String(), " ", "", -1)
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
    // Error checking
    println(bits.len)
    if start < 0 || start > bits.len { panic(errors.New("Start index "+strconv.Itoa(start)+" out of range.")) }
    if end < start || end > bits.len { panic(errors.New("End index "+strconv.Itoa(end)+" out of range.")) }
    if start == end { return &Bits{make([]byte,0), 0} }

    // Byte splicing
    dat := make([]byte, (end-start+7)/8)
    copy(dat, bits.dat[start/8:(end+7)/8])

    // Bit shifting
    shift := uint(start % 8)
    if shift > 0 {
        for i := 0; i < len(dat)-1; i++ {
            dat[i] = (dat[i] << shift) | (dat[i+1] >> (8-shift))
        }
        dat[len(dat)-1] <<= shift
    }

    // Bit chopping
    chop := uint(8 - (end-start) % 8)
    if chop != 8 {
        dat[len(dat)-1] = dat[len(dat)-1] & (byte(0xFF) << chop)
    }

    return &Bits{dat, end-start}
}

/********************
***   Operators   ***
********************/

func (bits *Bits) Equal (other *Bits) bool {
    return bytes.Equal(bits.dat, other.dat) && bits.len == other.len
}

func (bits *Bits) Plus (other *Bits) *Bits {
    // Byte splicing
    b := &Bits{append(bits.dat, other.dat...), bits.len + other.len}

    // Bit shifting
    shift := uint(8 - bits.len%8)
    if shift != 8 {
        b.dat[len(bits.dat)-1] |= b.dat[len(bits.dat)] >> (8-shift)
        for i := len(bits.dat); i < len(b.dat)-1; i++ {
            b.dat[i] = (b.dat[i] << shift) | (b.dat[i+1] >> (8-shift))
        }
        b.dat[len(b.dat)-1] <<= shift
        if b.len + other.len < 8 { b.dat = b.dat[:len(b.dat)-1] }
    }

    return b
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
    s := bits.RawString()
    if s == "" { s = "0" }
    x, e := strconv.ParseInt(s, 2,32)
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
        s, e, err := bits.chunkBounds(head)
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

func (bits *Bits) chunkBounds (head int) (start int, end int, err error) {
    if head+3 > bits.len { return 0, 0, errors.New("Decoding error: chunk header index "+strconv.Itoa(head+3)+" out of range.") }

    hl := 1 >> 1 << uint(bits.FromTo(head, head+3).Int())
    if head+3+hl > bits.len { return 0, 0, errors.New("Decoding error: chunk start index "+strconv.Itoa(head+3+hl)+" out of range.") }

    println("asdf", hl)
    l := bits.FromTo(head+3, head+3+hl).Int()
    if head+3+hl+l > bits.len { return 0, 0, errors.New("Decoding error: chunk end index "+strconv.Itoa(head+3+hl+hl)+" out of range.") }

    println("asdf", hl, l)
    return head+3+hl, head+3+hl+l, nil
}