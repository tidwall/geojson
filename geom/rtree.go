package geom

import "fmt"

const rTreeDims = 2
const rTreeM = 16 // must be 4,8,16,32,etc

// rTreeRect ...
type rTreeRect struct {
	min  [rTreeDims]float64
	max  [rTreeDims]float64
	data interface{}
}

func (r *rTreeRect) recalc() {
	n := r.data.(*rTreeNode)
	r.min = n.rects[0].min
	r.max = n.rects[0].min
	for i := 1; i < n.count; i++ {
		r.expand(&n.rects[i])
	}
}

func (r *rTreeRect) expand(b *rTreeRect) {
	for i := 0; i < rTreeDims; i++ {
		if b.min[i] < r.min[i] {
			r.min[i] = b.min[i]
		}
		if b.max[i] > r.max[i] {
			r.max[i] = b.max[i]
		}
	}
}

type rTreeNode struct {
	rects [rTreeM]rTreeRect
	count int
}

// rTree ...
type rTree struct {
	root rTreeRect
}

// load ...
func (tr *rTree) load(rects []rTreeRect) {
	axis := largesAxis(rects)
	tr.root = buildRect(rects, axis)
}

// contains return struct when b is fully contained inside of n
func (r *rTreeRect) intersects(b *rTreeRect) bool {
	for i := 0; i < rTreeDims; i++ {
		if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
			return false
		}
	}
	return true
}

func largesAxis(rects []rTreeRect) int {
	if len(rects) == 0 {
		return 0
	}
	rect := rects[0]
	for i := 1; i < len(rects); i++ {
		rect.expand(&rects[i])
	}
	axisSize := rect.max[0] - rect.min[0]
	axis := 0
	for i := 1; i < rTreeDims; i++ {
		sz := rect.max[i] - rect.min[i]
		if sz > axisSize {
			axisSize = sz
			axis = i
		}
	}
	return axis
}

func sortByAxis(rects []rTreeRect, axis int) {
	if len(rects) < 2 {
		return
	}
	left, right := 0, len(rects)-1
	pivotIndex := len(rects) / 2
	rects[pivotIndex], rects[right] = rects[right], rects[pivotIndex]
	for i := range rects {
		if rects[i].min[axis] < rects[right].min[axis] {
			rects[i], rects[left] = rects[left], rects[i]
			left++
		}
	}
	rects[left], rects[right] = rects[right], rects[left]
	sortByAxis(rects[:left], axis)
	sortByAxis(rects[left+1:], axis)
}

func buildRect(rects []rTreeRect, axis int) rTreeRect {
	if len(rects) == 0 {
		return rTreeRect{}
	} else if len(rects) == 1 {
		return rects[0]
	} else if len(rects) <= rTreeM {
		node := new(rTreeNode)
		copy(node.rects[:len(rects)], rects)
		node.count = len(rects)
		rect := rTreeRect{data: node}
		rect.recalc()
		return rect
	}

	sortByAxis(rects, axis%rTreeDims)

	node := new(rTreeNode)
	parts := len(rects) / rTreeM
	for i := 0; i < rTreeM; i++ {
		if i == rTreeM-1 {
			node.rects[node.count] = buildRect(rects, axis+1)
		} else {
			node.rects[node.count] = buildRect(rects[:parts], axis+1)
		}
		node.count++
		rects = rects[parts:]
	}
	var rect rTreeRect
	rect.data = node
	rect.recalc()
	return rect
}

func (r *rTreeRect) traverse(
	height int, iter func(rect rTreeRect, height int) bool,
) bool {
	if !iter(*r, height) {
		return false
	}
	if node, ok := r.data.(*rTreeNode); ok {
		for i := 0; i < node.count; i++ {
			if !node.rects[i].traverse(height+1, iter) {
				return false
			}
		}
	}
	return true
}

func (r *rTreeRect) search(
	rect rTreeRect, iter func(rect rTreeRect) bool,
) bool {
	if node, ok := r.data.(*rTreeNode); ok {
		for i := 0; i < node.count; i++ {
			if node.rects[i].intersects(&rect) {
				if _, ok := node.rects[i].data.(*rTreeNode); ok {
					if !node.rects[i].search(rect, iter) {
						return false
					}
				} else {
					if !iter(node.rects[i]) {
						return false
					}
				}
			}
		}
	}
	return true
}

func (tr *rTree) traverse(iter func(rect rTreeRect, height int) bool) {
	if tr.root.data == nil {
		return
	}
	tr.root.traverse(0, iter)
}

func (tr *rTree) search(
	rect rTreeRect, iter func(rect rTreeRect) bool,
) {
	if tr.root.data == nil {
		return
	}
	tr.root.search(rect, iter)
}

const svgScale = 4.0

var strokes = [...]string{"black", "#cccc00", "green", "red", "purple"}

func (tr *rTree) svg() string {
	var out string
	out += fmt.Sprintf("<svg viewBox=\"%.0f %.0f %.0f %.0f\" "+
		"xmlns =\"http://www.w3.org/2000/svg\">\n",
		-190.0*svgScale, -100.0*svgScale,
		380.0*svgScale, 190.0*svgScale)
	out += fmt.Sprintf("<g transform=\"scale(1,-1)\">\n")
	var outb []byte
	tr.traverse(func(rect rTreeRect, height int) bool {
		outb = append(outb, svg(rect, height)...)
		return true
	})
	out += string(outb)
	out += fmt.Sprintf("</g>\n")
	out += fmt.Sprintf("</svg>\n")
	return out
}

func svg(rect rTreeRect, height int) string {
	min, max := rect.min, rect.max
	var out string
	point := true
	for i := 0; i < 2; i++ {
		if min[i] != max[i] {
			point = false
			break
		}
	}
	if point { // is point
		out += fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%0.f\" width=\"%0.f\" height=\"%0.f\" "+
				"stroke=\"%s\" fill=\"purple\" "+
				"fill-opacity=\"0\" stroke-opacity=\"1\" "+
				"rx=\"15\" ry=\"15\"/>\n",
			(min[0])*svgScale,
			(min[1])*svgScale,
			(max[0]-min[0]+1/svgScale)*svgScale,
			(max[1]-min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])
	} else { // is rect
		out += fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%0.f\" width=\"%0.f\" height=\"%0.f\" "+
				"stroke=\"%s\" fill=\"purple\" "+
				"fill-opacity=\"0\" stroke-opacity=\"1\"/>\n",
			(min[0])*svgScale,
			(min[1])*svgScale,
			(max[0]-min[0]+1/svgScale)*svgScale,
			(max[1]-min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])
	}
	return out
}

// package geom

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

// func rectUnion(rect, other Rect) Rect {
// 	if other.Min.X < rect.Min.X {
// 		rect.Min.X = other.Min.X
// 	} else if other.Max.X > rect.Max.X {
// 		rect.Max.X = other.Max.X
// 	}
// 	if other.Min.Y < rect.Min.Y {
// 		rect.Min.Y = other.Min.Y
// 	} else if other.Max.Y > rect.Max.Y {
// 		rect.Max.Y = other.Max.Y
// 	}
// 	return rect
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
