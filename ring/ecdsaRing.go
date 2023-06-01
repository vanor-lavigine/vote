package ring

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"io"
	"math/big"
	"sync"
)

type PublicKeyRing struct {
	Ring []ecdsa.PublicKey
}

//
func (r *PublicKeyRing) LenList() int {
	return len(r.Ring)
}

func (r *PublicKeyRing) Bytes() (b []byte) {
	for _, pub := range r.Ring {
		b = append(b, pub.X.Bytes()...)
		b = append(b, pub.Y.Bytes()...)
	}
	return
}

var one = new(big.Int).SetInt64(1)

func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return
	}

	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}

func hashG(c elliptic.Curve, m []byte) (hx, hy *big.Int) {
	h := sha256.New()
	h.Write(m)
	d := h.Sum(nil)
	hx, hy = c.ScalarBaseMult(d) // g^H'()
	return
}

func hashAllq(mR []byte, ax, ay, bx, by []*big.Int) (hash *big.Int) {
	h := sha256.New()
	h.Write(mR)
	for i := 0; i < len(ax); i++ {
		h.Write(ax[i].Bytes())
		h.Write(ay[i].Bytes())
		h.Write(bx[i].Bytes())
		h.Write(by[i].Bytes())
	}
	hash = new(big.Int).SetBytes(h.Sum(nil))
	return
}

func hashAllqc(c elliptic.Curve, mR []byte, ax, ay, bx, by []*big.Int) (hash *big.Int) {
	h := sha256.New()
	h.Write(mR)
	for i := 0; i < len(ax); i++ {
		h.Write(ax[i].Bytes())
		h.Write(ay[i].Bytes())
		h.Write(bx[i].Bytes())
		h.Write(by[i].Bytes())
	}
	hash = hashToInt(h.Sum(nil), c)
	return
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

type RingSign struct {
	X, Y *big.Int
	C, T []*big.Int
}

func SignRing(rand io.Reader, sk ecdsa.PrivateKey, pkList PublicKeyRing, m []byte) (rs *RingSign, err error) {
	sList := pkList.LenList()
	ax := make([]*big.Int, sList, sList)
	ay := make([]*big.Int, sList, sList)
	bx := make([]*big.Int, sList, sList)
	by := make([]*big.Int, sList, sList)
	c := make([]*big.Int, sList, sList)
	t := make([]*big.Int, sList, sList)
	pub := sk.PublicKey
	curve := pub.Curve
	N := curve.Params().N
	mR := append(m, pkList.Bytes()...)
	hx, hy := hashG(curve, mR) // H(mR)
	var id int
	var wg sync.WaitGroup
	sum := new(big.Int).SetInt64(0)
	for j := 0; j < sList; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			c[j], err = randFieldElement(curve, rand)
			if err != nil {
				return
			}
			t[j], err = randFieldElement(curve, rand)
			if err != nil {
				return
			}

			if pkList.Ring[j] == pub {
				id = j
				rb := t[j].Bytes()
				ax[id], ay[id] = curve.ScalarBaseMult(rb)     // g^r
				bx[id], by[id] = curve.ScalarMult(hx, hy, rb) // H(mR)^r
			} else {
				ax1, ay1 := curve.ScalarBaseMult(t[j].Bytes())                                 // g^tj
				ax2, ay2 := curve.ScalarMult(pkList.Ring[j].X, pkList.Ring[j].Y, c[j].Bytes()) // yj^cj
				ax[j], ay[j] = curve.Add(ax1, ay1, ax2, ay2)

				w := new(big.Int)
				w.Mul(sk.D, c[j])
				w.Add(w, t[j])
				w.Mod(w, N)
				bx[j], by[j] = curve.ScalarMult(hx, hy, w.Bytes()) // H(mR)^(xi*cj+tj)
				// TODO may need to lock on sum object.
				sum.Add(sum, c[j]) // Sum needed in Step 3 of the algorithm
			}
		}(j)
	}
	wg.Wait()
	// Step 3, part 1: cid = H(m,R,{a,b}) - sum(cj) mod N
	hashmRab := hashAllq(mR, ax, ay, bx, by)
	// hashmRab := hashAllqc(curve, mR, ax, ay, bx, by)
	c[id].Sub(hashmRab, sum)
	c[id].Mod(c[id], N)

	// Step 3, part 2: tid = ri - cid * xi mod N
	cx := new(big.Int)
	cx.Mul(sk.D, c[id])
	t[id].Sub(t[id], cx) // here t[id] = ri (initialized inside the for-loop above)
	t[id].Mod(t[id], N)

	hsx, hsy := curve.ScalarMult(hx, hy, sk.D.Bytes()) // Step 4: H(mR)^xi
	return &RingSign{hsx, hsy, c, t}, nil

}
