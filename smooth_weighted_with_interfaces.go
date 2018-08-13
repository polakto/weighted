package weighted

// smoothWeighted is a wrapped weighted item.
// type smoothWeighted struct {
// 	Item            interface{}
// 	Weight          int
// 	CurrentWeight   int
// 	EffectiveWeight int
// }

type iSmoothWeighted interface {
	Item() interface{}
	Weight() int
	CurrentWeight() int
	EffectiveWeight() int
	SetEffectiveWeight(int)
	SetCurrentWeight(int)
}

/*
SW (Smooth Weighted) is a struct that contains weighted items and provides methods to select a weighted item.
It is used for the smooth weighted round-robin balancing algorithm. This algorithm is implemented in Nginx:
https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35.

Algorithm is as follows: on each peer selection we increase current_weight
of each eligible peer by its weight, select peer with greatest current_weight
and reduce its current_weight by total number of weight points distributed
among peers.

In case of { 5, 1, 1 } weights this gives the following sequence of
current_weight's: (a, a, b, a, c, a, a)

*/
type SWI struct {
	items []iSmoothWeighted
	n     int
}

// func (w *smoothWeighted) fail() {
// 	w.EffectiveWeight -= w.Weight
// 	if w.EffectiveWeight < 0 {
// 		w.EffectiveWeight = 0
// 	}
// }

// Add a weighted server.
func (w *SWI) Add(item iSmoothWeighted) {
	w.items = append(w.items, item)
	w.n++
}

// RemoveAll removes all weighted items.
func (w *SWI) RemoveAll() {
	w.items = w.items[:0]
	w.n = 0
}

//Reset resets all current weights.
func (w *SWI) Reset() {
	for _, s := range w.items {
		s.SetEffectiveWeight(s.Weight())
		s.SetCurrentWeight(0)
	}
}

// All returns all items.
func (w *SWI) All() []interface{} {
	allItems := make([]interface{}, 0)
	for _, i := range w.items {
		allItems = append(allItems, i.Item())
	}
	return allItems
}

// Next returns next selected server.
func (w *SWI) Next() interface{} {
	i := w.nextWeighted()
	if i == nil {
		return nil
	}
	return i.Item()
}

// nextWeighted returns next selected weighted object.
func (w *SWI) nextWeighted() iSmoothWeighted {
	if w.n == 0 {
		return nil
	}
	if w.n == 1 {
		return w.items[0]
	}

	return nextSmoothWeightedInterfaces(w.items)
}

//https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
func nextSmoothWeightedInterfaces(items []iSmoothWeighted) (best iSmoothWeighted) {
	total := 0

	for i := 0; i < len(items); i++ {
		w := items[i]

		if w == nil {
			continue
		}

		// w.CurrentWeight += w.EffectiveWeight
		// total += w.EffectiveWeight
		// if w.EffectiveWeight < w.Weight {
		// 	w.EffectiveWeight++
		// }

		w.SetCurrentWeight(w.CurrentWeight() + w.EffectiveWeight())
		total += w.EffectiveWeight()
		if w.EffectiveWeight() < w.Weight() {
			w.SetEffectiveWeight(w.EffectiveWeight() + 1)
		}

		if best == nil || w.CurrentWeight() > best.CurrentWeight() {
			best = w
		}

	}

	if best == nil {
		return nil
	}

	// best.CurrentWeight -= total
	best.SetCurrentWeight(best.CurrentWeight() - total)
	return best
}
