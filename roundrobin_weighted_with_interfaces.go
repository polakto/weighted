package weighted

// // rrWeighted is a wrapped weighted item that is used to implement LVS weighted round robin algorithm.
// type rrWeighted struct {
// 	Item   interface{}
// 	Weight int
// }

type IrrWeighted interface {
	Item() interface{}
	Weight() int
}

// RRW is a struct that contains weighted items implement LVS weighted round robin algorithm.
//
// http://kb.linuxvirtualitem.org/wiki/Weighted_Round-Robin_Scheduling
// http://zh.linuxvirtualitem.org/node/37
type RRWI struct {
	items []IrrWeighted
	n     int
	gcd   int
	maxW  int
	i     int
	cw    int
}

// Add a weighted item.
func (w *RRWI) Add(item IrrWeighted) {
	// weighted := &rrWeighted{Item: item, Weight: weight}
	if item.Weight() > 0 {
		if w.gcd == 0 {
			w.gcd = item.Weight()
			w.maxW = item.Weight()
			w.i = -1
			w.cw = 0
		} else {
			w.gcd = gcdI(w.gcd, item.Weight())
			if w.maxW < item.Weight() {
				w.maxW = item.Weight()
			}
		}
	}
	w.items = append(w.items, item)
	w.n++
}

// All returns all items.
func (w *RRWI) All() []interface{} {
	allItems := make([]interface{}, 0)
	for _, i := range w.items {
		allItems = append(allItems, i.Item())
	}
	return allItems
}

// RemoveAll removes all weighted items.
func (w *RRWI) RemoveAll() {
	w.items = w.items[:0]
	w.n = 0
	w.gcd = 0
	w.maxW = 0
	w.i = -1
	w.cw = 0
}

//Reset resets all current weights.
func (w *RRWI) Reset() {
	w.i = -1
	w.cw = 0
}

// Next returns next selected item.
func (w *RRWI) Next() interface{} {
	// fmt.Printf("selecting from state \n")
	// for k := range w.items {
	// 	fmt.Printf("%+v \n", w.items[k])
	// }

	if w.n == 0 {
		return nil
	}

	if w.n == 1 {
		return w.items[0].Item()
	}

	for {
		w.i = (w.i + 1) % w.n
		if w.i == 0 {
			w.cw = w.cw - w.gcd
			if w.cw <= 0 {
				w.cw = w.maxW
				if w.cw == 0 {
					return nil
				}
			}
		}

		if w.items[w.i].Weight() >= w.cw {
			return w.items[w.i].Item()
		}
	}
}

func gcdI(x, y int) int {
	var t int
	for {
		t = (x % y)
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}
