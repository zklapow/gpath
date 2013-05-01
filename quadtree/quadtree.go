package quadtree

import "errors"

const (
    MAXLEVELS = 5
    MAXOBJECTS = 10
)

type Rect struct {
    X       int
    Y       int
    Width   int
    Height  int
}

type QuadTree struct {
    Bounds  Rect
    Level   int

    Nodes   Trees
    Objects Rects
}

type Trees []*QuadTree
func (l Trees) Count() int {
    c := 0
    for _, val := range l {
        if val != nil {c++}
    }
    return c
}

type Rects []*Rect
func (l Rects) Count() int {
    c := 0
    for _, val := range l {
        if val != nil {c++}
    }
    return c
}

func New(level int, bounds Rect) (ret *QuadTree) {
    ret = &QuadTree{Bounds: bounds, Level: level}

    ret.Nodes = make(Trees, 4)
    ret.Objects = make(Rects, MAXOBJECTS)

    return ret
}


func (node *QuadTree) Clear() {
    // Make a new empty list
    node.Objects = make([]*Rect, MAXOBJECTS)

    // If there are no object then there
    // can't be any child nodes
    node.Nodes = make([]*QuadTree, 4)
}

func (node *QuadTree) Split() {
    nw := node.Bounds.Width/2
    nh := node.Bounds.Height/2

    x := node.Bounds.X
    y := node.Bounds.Y

    node.Nodes[0] = New(node.Level+1, Rect{X: x+nw, Y: y,    Width: nw, Height: nh})
    node.Nodes[1] = New(node.Level+1, Rect{X: x,    Y: y,    Width: nw, Height: nh})
    node.Nodes[2] = New(node.Level+1, Rect{X: x,    Y: y+nh, Width: nw, Height: nh})
    node.Nodes[3] = New(node.Level+1, Rect{X: x+nw, Y: y+nh, Width: nw, Height: nh})
}

func (node *QuadTree) index(obj *Rect) (index int) {
    // Calculate what quadrant the object should be placed in
    // -1 indicates the object does not fit entirely within
    // a child node, thus it belong to this node
    index = -1

    xmid := node.Bounds.X + (node.Bounds.Width / 2)
    ymid := node.Bounds.Y + (node.Bounds.Height / 2)

    top := (obj.Y < ymid && obj.Y + obj.Height < ymid)
    bot := (obj.Y > ymid)

    // Left half
    if obj.X < xmid && (obj.X + obj.Width) < xmid {
        if top {
            index = 1
        } else if bot {
            index = 2
        }
    // Right half
    }else if obj.X > xmid {
        if top {
            index = 0
        } else if bot {
            index = 3
        }
    }

    return index
}

func (node *QuadTree) Insert(obj *Rect) (err error) {
    if node.Nodes.Count() != 0 {
        index := node.index(obj)

        // This object belongs to a child
        if index != -1 {
            node.Nodes[index].Insert(obj)
            return nil
        }
    }
        
    // Check if its time to split this node
    if node.Objects.Count() == MAXOBJECTS {
        // Create a new list for any object which need to stay here
        tmp := make(Rects, MAXOBJECTS)
        node.Split()

        for i := 0; i < node.Objects.Count(); i++ {
            if node.Objects[i] != nil {
                index := node.index(node.Objects[i])
                if index != -1 {
                    node.Nodes[index].Insert(node.Objects[i])
                } else {
                    tmp[tmp.Count()] = node.Objects[i]
                }
            }
        }

        if tmp.Count() == MAXOBJECTS {
            // If the new list is the same size we have a problem
            node.Objects = tmp
            return errors.New("Quadtree appears to have been overloaded!")
        } else {
            // We can now figure out where to put our new object
            index := node.index(obj)
            if index != -1 {
                node.Nodes[index].Insert(obj)
            } else {
                tmp[tmp.Count()] = obj
            }

            node.Objects = tmp
            return nil
        }
    } else {
        // Our object list has space
        node.Objects[node.Objects.Count()] = obj
    }

    return nil
}

func (node *QuadTree) childObjects() (ret Rects) {
    // Get and return all child object of a node
    ret = make(Rects, 0)
    if node.Nodes.Count() != 0 {
        for _, cn := range node.Nodes {
            ret = append(ret, cn.Objects...)
            ret = append(ret, cn.childObjects()...)
        }
    }

    return ret
}

func (node *QuadTree) Retrieve(obj *Rect) (ret Rects) {
    // Build a list of all objects we could be in collsion with
    ret = make(Rects, 0)

    index := node.index(obj)
    if index != -1 && node.Nodes.Count() != 0 {
        ret = append(ret, node.Nodes[index].Retrieve(obj)...)
    } else {
        // If the object is in this node it could collide
        // with any of this nodes children
        ret = append(ret, node.childObjects()...)
    }

    // Any object in this node could collide with a child
    ret = append(ret, node.Objects...)

    return ret
}
