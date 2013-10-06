package fiddle

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

func TestChunks (t *testing.T) {
    for i := 0; i < 1000; i++ {
        c     := []*Bits{randBits(), randBits(), randBits()}
        c2    := FromChunks(c)
        c3, e := c2.Chunks(3)

        if  e != nil ||
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

func TestList (t *testing.T) {
    for i := 0; i < 1000; i++ {
        c     := []*Bits{randBits(), randBits(), randBits()}
        c2    := FromList(c)
        c3, e := c2.List()

        if  e != nil ||
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