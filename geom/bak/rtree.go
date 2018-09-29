package geom

import (
	"fmt"
	"math/bits"
)

const rTreeDims = 2
const rTreeM = 16 // must be 2, 4, 8, 16, 32, etc

// rTreeRect ...
type rTreeRect struct {
	min  [rTreeDims]float64
	max  [rTreeDims]float64
	data interface{}
}

func (r *rTreeRect) recalc() {
	n := r.data.(*rTreeNode)
	r.min = n.rects[0].min
	r.max = n.rects[0].max
	for i := 0; i < n.count; i++ {
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

// contains return struct when b is fully contained inside of n
func (r *rTreeRect) intersects(b *rTreeRect) bool {
	for i := 0; i < rTreeDims; i++ {
		if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
			return false
		}
	}
	return true
}

type overlapsResult int

const (
	not overlapsResult = iota
	intersects
	contains
)

// overlaps detects if r insersects or contains b.
// return not, intersects, contains
func (r *rTreeRect) overlaps(b *rTreeRect) overlapsResult {
	for i := 0; i < rTreeDims; i++ {
		if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
			return not
		}
		if r.min[i] > b.min[i] || b.max[i] > r.max[i] {
			i++
			for ; i < rTreeDims; i++ {
				if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
					return not
				}
			}
			return intersects
		}
	}
	return contains
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

// load ...
func (tr *rTree) load(rects []rTreeRect) {
	tr.root = buildRect(rects)
	tr.root.recalc()
}

func buildRect(rects []rTreeRect) rTreeRect {
	var rect rTreeRect
	build := breakUpRects(rects, nil, 0, bits.TrailingZeros(rTreeM))
	if len(build) > 0 {
		node := new(rTreeNode)
		for _, rects := range build {
			var child rTreeRect
			if len(rects) <= rTreeM {
				childNode := new(rTreeNode)
				copy(childNode.rects[:], rects)
				childNode.count = len(rects)
				child = rTreeRect{data: childNode}
			} else {
				child = buildRect(rects)
			}
			child.recalc()
			node.rects[node.count] = child
			node.count++
		}
		rect.data = node
		rect.recalc()
	}
	return rect
}

func breakUpRects(
	rects []rTreeRect,
	build [][]rTreeRect,
	depth, maxDepth int,
) [][]rTreeRect {
	//fmt.Printf("** axis: %v, depth: %v, maxDepth: %v, maxEntries: %v\n",
	//	axis, depth, maxDepth, maxEntries)
	if len(rects) == 0 {
		return build
	} else if len(rects) <= rTreeM || depth == maxDepth {
		build = append(build, rects)
		return build
	}
	sortByAxis(rects, largesAxis(rects))
	left := rects[:len(rects)/2]
	right := rects[len(rects)/2:]
	build = breakUpRects(left, build, depth+1, maxDepth)
	build = breakUpRects(right, build, depth+1, maxDepth)
	return build
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

// func (r *rTreeRect) search(
// 	rect rTreeRect, iter func(rect rTreeRect) bool,
// ) bool {
// 	if node, ok := r.data.(*rTreeNode); ok {
// 		for i := 0; i < node.count; i++ {
// 			if node.rects[i].intersects(&rect) {
// 				if _, ok := node.rects[i].data.(*rTreeNode); ok {
// 					if !node.rects[i].search(rect, iter) {
// 						return false
// 					}
// 				} else {
// 					if !iter(node.rects[i]) {
// 						return false
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return true
// }

func (r *rTreeRect) search(
	target rTreeRect, iter func(rect rTreeRect) bool,
) bool {
	n := r.data.(*rTreeNode)
	var branch bool
	for i := 0; i < n.count; i++ {
		if i == 0 {
			_, branch = n.rects[i].data.(*rTreeNode)
		}
		if branch {
			res := target.overlaps(&n.rects[i])
			if res == intersects {
				if !n.rects[i].search(target, iter) {
					return false
				}

			} else if res == contains {
				if !n.rects[i].scan(iter) {
					return false
				}
			}
		} else {
			if target.intersects(&n.rects[i]) {
				if !iter(n.rects[i]) {
					return false
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
	target rTreeRect, iter func(rect rTreeRect) bool,
) {
	if tr.root.data == nil {
		return
	}
	res := target.overlaps(&tr.root)
	if res == intersects {
		tr.root.search(target, iter)
	} else if res == contains {
		tr.root.scan(iter)
	}
}

func (r *rTreeRect) scan(iter func(rect rTreeRect) bool) bool {
	n := r.data.(*rTreeNode)
	var branch bool
	for i := 0; i < n.count; i++ {
		if i == 0 {
			_, branch = n.rects[i].data.(*rTreeNode)
		}
		if branch {
			if !n.rects[i].scan(iter) {
				return false
			}
		} else {
			if !iter(n.rects[i]) {
				return false
			}
		}
	}
	return true
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
