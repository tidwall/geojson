package rtree

// const rTreeDims = 2
// const rTreeM = 16 // must be 4,8,16,32,etc

// // rTreeRect ...
// type rTreeRect struct {
// 	min  [rTreeDims]float64
// 	package geom

// import (
// 	"math"
// )

// const rTreeM = 16

// type rTreeNode struct {
// 	rect     Rect                // bounding rectangle
// 	children [rTreeM]interface{} // either a node or an index to item
// 	count    int                 // number of children
// }

// type rTree struct {
// 	items []Rect
// 	root  *rTreeNode
// }

// func newRTree(items []Rect) *rTree {
// 	tr := new(rTree)
// 	if len(items) == 0 {
// 		return tr
// 	}
// 	// load with OMT
// 	tr.items = items
// 	treeHeight := int(math.Ceil(math.Log(float64(len(tr.items))) /
// 		math.Log(rTreeM)))
// 	rootMaxEntries := int(math.Ceil(float64(len(tr.items)) /
// 		math.Pow(rTreeM, float64(treeHeight-1))))
// 	tr.root = tr.buildNodes(0, len(tr.items)-1, treeHeight, rootMaxEntries)
// 	return tr
// }

// func sortByAxis(items []Rect, axis int) {
// 	if len(items) < 2 {
// 		return
// 	}
// 	left, right := 0, len(items)-1
// 	pivotIndex := len(items) / 2
// 	items[pivotIndex], items[right] = items[right], items[pivotIndex]
// 	for i := range items {
// 		var less bool
// 		if axis == 0 {
// 			less = items[i].Min.X < items[right].Min.X
// 		} else if axis == 1 {
// 			less = items[i].Min.Y < items[right].Min.Y
// 		}
// 		if less {
// 			items[i], items[left] = items[left], items[i]
// 			left++
// 		}
// 	}
// 	items[left], items[right] = items[right], items[left]
// 	sortByAxis(items[:left], axis)
// 	sortByAxis(items[left+1:], axis)
// }

// func (tr *rTree) buildNodes(left, right, height, maxEntries int) *rTreeNode {
// 	var num = right - left + 1
// 	if num <= maxEntries {
// 		newNode := new(rTreeNode)
// 		items := tr.items[left : left+num]
// 		for i, item := range items {
// 			index := left + i
// 			newNode.children[newNode.count] = index
// 			if newNode.count == 0 {
// 				newNode.rect = item
// 			} else {
// 				newNode.rect = rectUnion(newNode.rect, item)
// 			}
// 			newNode.count++
// 		}
// 		return newNode
// 	}
// 	sortByAxis(tr.items[left:left+num], 0)
// 	nodeSize := (num + (maxEntries - 1)) / maxEntries
// 	subSortLength := nodeSize * int(math.Ceil(math.Sqrt(float64(maxEntries))))
// 	newNode := new(rTreeNode)
// 	for subCounter := left; subCounter <= right; subCounter += subSortLength {
// 		subRight := int(math.Min(float64(subCounter+subSortLength-1),
// 			float64(right)))
// 		sortByAxis(tr.items[subCounter:subCounter+(subRight-subCounter+1)], 1)
// 		for nodeCounter := subCounter; nodeCounter <= subRight; nodeCounter += nodeSize {
// 			child := tr.buildNodes(
// 				nodeCounter,
// 				int(math.Min(float64(nodeCounter+nodeSize-1), float64(subRight))),
// 				height-1,
// 				rTreeM)
// 			newNode.children[newNode.count] = child
// 			if newNode.count == 0 {
// 				newNode.rect = child.rect
// 			} else {
// 				newNode.rect = rectUnion(newNode.rect, child.rect)
// 			}
// 			newNode.count++
// 		}
// 	}
// 	return newNode
// }

// // type overlapsResult int

// // const (
// // 	not overlapsResult = iota
// // 	intersects
// // 	contains
// // )

// // // overlaps detects if r insersects or contains b.
// // // return not, intersects, contains
// // func (r *rTreeNode) overlaps(node *box) overlapsResult {
// // 	for i := 0; i < dims; i++ {
// // 		if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
// // 			return not
// // 		}
// // 		if r.min[i] > b.min[i] || b.max[i] > r.max[i] {
// // 			i++
// // 			for ; i < dims; i++ {
// // 				if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
// // 					return not
// // 				}
// // 			}
// // 			return intersects
// // 		}
// // 	}
// // 	return contains
// // }

// func (tr *rTree) search(rect Rect, iter func(rect Rect, index int) bool) {
// 	if tr.root == nil {
// 		return
// 	}
// 	if rect.IntersectsRect(tr.root.rect) {
// 		tr.searchNode(tr.root, rect, iter)
// 	}
// }

// func (tr *rTree) searchNode(
// 	node *rTreeNode, rect Rect, iter func(rect Rect, index int) bool,
// ) bool {
// 	for i := 0; i < node.count; i++ {
// 		switch v := node.children[i].(type) {
// 		case int:
// 			if tr.items[v].IntersectsRect(rect) {
// 				if !iter(tr.items[v], v) {
// 					return false
// 				}
// 			}
// 		case *rTreeNode:
// 			if v.rect.IntersectsRect(rect) {
// 				if !tr.searchNode(v, rect, iter) {
// 					return false
// 				}
// 			}
// 		}
// 	}
// 	return true
// }
// max  [rTreeDims]float64
// 	data interface{}
// }

// func (r *rTreeRect) recalc() {
// 	n := r.data.(*rTreeNode)
// 	r.min = n.rects[0].min
// 	r.max = n.rects[0].min
// 	for i := 1; i < n.count; i++ {
// 		r.expand(&n.rects[i])
// 	}
// }

// func (r *rTreeRect) expand(b *rTreeRect) {
// 	for i := 0; i < rTreeDims; i++ {
// 		if b.min[i] < r.min[i] {
// 			r.min[i] = b.min[i]
// 		}
// 		if b.max[i] > r.max[i] {
// 			r.max[i] = b.max[i]
// 		}
// 	}
// }

// type rTreeNode struct {
// 	rects [rTreeM]rTreeRect
// 	count int
// }

// // rTree ...
// type rTree struct {
// 	root rTreeRect
// }

// // load ...
// func (tr *rTree) load(rects []rTreeRect) {
// 	axis := largesAxis(rects)
// 	tr.root = buildRect(rects, axis)
// }

// func largesAxis(rects []rTreeRect) int {
// 	if len(rects) == 0 {
// 		return 0
// 	}
// 	rect := rects[0]
// 	for i := 1; i < len(rects); i++ {
// 		rect.expand(&rects[i])
// 	}
// 	axisSize := rect.max[0] - rect.min[0]
// 	axis := 0
// 	for i := 1; i < rTreeDims; i++ {
// 		sz := rect.max[i] - rect.min[i]
// 		if sz > axisSize {
// 			axisSize = sz
// 			axis = i
// 		}
// 	}
// 	return axis
// }

// func sortByAxis(rects []rTreeRect, axis int) {
// 	if len(rects) < 2 {
// 		return
// 	}
// 	left, right := 0, len(rects)-1
// 	pivotIndex := len(rects) / 2
// 	rects[pivotIndex], rects[right] = rects[right], rects[pivotIndex]
// 	for i := range rects {
// 		if rects[i].min[axis] < rects[right].min[axis] {
// 			rects[i], rects[left] = rects[left], rects[i]
// 			left++
// 		}
// 	}
// 	rects[left], rects[right] = rects[right], rects[left]
// 	sortByAxis(rects[:left], axis)
// 	sortByAxis(rects[left+1:], axis)
// }

// func buildRect(rects []rTreeRect, axis int) rTreeRect {
// 	if len(rects) == 0 {
// 		return rTreeRect{}
// 	} else if len(rects) == 1 {
// 		return rects[0]
// 	} else if len(rects) <= 16 {
// 		node := new(rTreeNode)
// 		copy(node.rects[:len(rects)], rects)
// 		node.count = len(rects)
// 		rect := rTreeRect{data: node}
// 		rect.recalc()
// 		return rect
// 	}
// 	sortByAxis(rects, axis%rTreeDims)
// 	node := new(rTreeNode)
// 	parts := len(rects) / rTreeM
// 	for i := 0; i < rTreeM; i++ {
// 		if i == rTreeM-1 {
// 			node.rects[node.count] = buildRect(rects, axis+1)
// 		} else {
// 			node.rects[node.count] = buildRect(rects[:parts], axis+1)
// 		}
// 		node.count++
// 		rects = rects[parts:]
// 	}
// 	var rect rTreeRect
// 	rect.data = node
// 	rect.recalc()
// 	return rect
// }
