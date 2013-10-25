package fiddle

import "math/big"
import "math/rand"
import "testing"
import "time"

func randBits () *Bits {
    b := FromInt(rand.Int())
    return b.To(rand.Intn(b.Len()+1))
}

func init () {
    rand.Seed(time.Now().UnixNano())
}

func TestInt (t *testing.T) {
    for i := 0; i < 1000; i++ {
        x  := rand.Int()
        b  := FromInt(x)
        x2 := b.Int()

        if x2 != x {
            t.Log("Original:", x)
            t.Log("Encoded: ", b)
            t.Log("Decoded: ", x2)
            t.FailNow()
        }
    }
}

func TestBigInt (t *testing.T) {
    for i := 0; i < 1000; i++ {
        x  := big.NewInt(int64(rand.Int())); x.Mul(x,x); x.Mul(x,x)
        b  := FromBigInt(x)
        x2 := b.BigInt()

        if x2.Cmp(x) != 0 {
            t.Log("Original:", x)
            t.Log("Encoded: ", b)
            t.Log("Decoded: ", x2)
            t.FailNow()
        }
    }
}

func TestPlus (t *testing.T) {
    for i := 0; i < 1000; i++ {
        x  := randBits()
        y  := randBits()
        z  := x.Plus(y)
        z2 := FromBin(x.Bin() + y.Bin())

        if !z2.Equal(z) {
            t.Log("First:   ", x.Bin())
            t.Log("Second:  ", y.Bin())
            t.Log("Expected:", z.Bin())
            t.Log("Got:     ", z2.Bin())
            t.FailNow()
        }
    }
}

func TestFromTo (t *testing.T) {
    for i := 0; i < 1000; i++ {
        x  := randBits()
        s  := 0
        l  := 0
        if x.Len() > 0 {
            s = rand.Intn(x.Len())
            l = rand.Intn(x.Len()-s)
        }
        y  := x.Bin()[s:s+l]
        y2 := x.FromTo(s, s+l).Bin()
        
        if y2 != y {
            t.Log("Original:", x)
            t.Log("Start:   ", s)
            t.Log("End:     ", s+l)
            t.Log("Expected:", y)
            t.Log("Got:     ", y2)
            t.FailNow()
        }
    }
}

func TestBytes (t *testing.T) {
    for i := 0; i < 1000; i++ {
        x  := randBits()
        b  := x.Bytes()
        x2 := FromBytes(b)

        if !x2.Equal(x) {
            t.Log("Original:", x)
            t.Log("Encoded: ", b)
            t.Log("Decoded: ", x2)
            t.FailNow()
        }
    }
}

func TestHex (t *testing.T) {
    for i := 0; i < 1000; i++ {
        b  := randBits()
        b2 := FromHex(b.Hex())

        if !b2.Equal(b) {
            t.Log("Original:", b)
            t.Log("Encoded: ", b.Hex())
            t.Log("Decoded: ", b2)
            t.FailNow()
        }
    }
}

func TestRawHex (t *testing.T) {
    for i := 0; i < 1000; i++ {
        h  := randBits().RawHex()
        h2 := FromRawHex(h).RawHex()

        if h2 != h {
            t.Log("Original:", h)
            t.Log("Encoded: ", FromRawHex(h))
            t.Log("Decoded: ", h2)
            t.FailNow()
        }
    }
}

func TestChunks (t *testing.T) {
    for i := 0; i < 1000; i++ {
        c     := []*Bits{randBits(), randBits(), randBits()}
        b     := FromChunks(c...)
        c2 := b.Chunks(3)

        if !c2[0].Equal(c[0]) ||
           !c2[1].Equal(c[1]) ||
           !c2[2].Equal(c[2]) {
            t.Log("Chunk 0:  ", c[0].Bin())
            t.Log("Chunk 1:  ", c[1].Bin())
            t.Log("Chunk 2:  ", c[2].Bin())
            t.Log("Encoded:  ", b.Bin())
            t.Log("Decoded 0:", c2[0].Bin())
            t.Log("Decoded 1:", c2[1].Bin())
            t.Log("Decoded 2:", c2[2].Bin())
            t.FailNow()
        }
    }
}

func TestList (t *testing.T) {
    for i := 0; i < 1000; i++ {
        c     := []*Bits{randBits(), randBits(), randBits()}
        c2    := FromList(c)
        c3, e := c2.List()

        if e != nil ||
        !c3[0].Equal(c[0]) ||
        !c3[1].Equal(c[1]) ||
        !c3[2].Equal(c[2]) {
            t.Log("Chunk 0:  ", c[0].Bin())
            t.Log("Chunk 1:  ", c[1].Bin())
            t.Log("Chunk 2:  ", c[2].Bin())
            t.Log("Encoded:  ", c2.Bin())
            t.Log("Error:    ", e)
            t.Log("Decoded 0:", c3[0].Bin())
            t.Log("Decoded 1:", c3[1].Bin())
            t.Log("Decoded 2:", c3[2].Bin())
            t.FailNow()
        }
    }
}