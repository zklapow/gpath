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

    Nodes   []*QuadTree
    Objects []*Rect
}

func New(level int, bounds Rect) (ret *QuadTree) {
    ret = &QuadTree{Bounds: bounds, Level: level}

    ret.Nodes = make([]*QuadTree, 4)
    ret.Objects = make([]*Rect, MAXOBJECTS)

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
    nw := int(node.Bounds.Width/2)
    nh := int(node.Bounds.Height/2)

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
    if len(node.Nodes) != 0 {
        index := node.index(obj)

        // This object belongs to a child
        if index != -1 {
            node.Nodes[index].Insert(obj)
        }
    }
        
    // Check if its time to split this node
    if len(node.Objects) == MAXOBJECTS {
        // Create a new list for any object which need to stay here
        tmp := make([]*Rect, MAXOBJECTS)
        node.Split()

        for _, n := range node.Objects {
            index := node.index(n)
            if index != -1 {
                node.Nodes[index].Insert(n)
            } else {
                tmp[len(tmp)+1] = n
            }
        }

        if len(tmp) == MAXOBJECTS {
            // If the new list is the same size we have a problem
            node.Objects = tmp
            return errors.New("Quadtree appears to have been overloaded!")
        } else {
            // We can now figure out where to put our new object
            index := node.index(obj)
            if index != -1 {
                node.Nodes[index].Insert(obj)
            } else {
                tmp[len(tmp)+1] = obj
                node.Objects = tmp
            }
        }
    } else {
        // Our object list has space
        node.Objects[len(node.Objects)+1] = obj
    }

    return nil
}

func (node *QuadTree) childObjects() (ret []*Rect) {
    // Get and return all child object of a node
    ret = make([]*Rect, 0)
    if len(node.Nodes) != 0 {
        for _, cn := range node.Nodes {
            ret = append(ret, cn.Objects...)
            ret = append(ret, cn.childObjects()...)
        }
    }

    return ret
}

func (node *QuadTree) Retrieve(obj *Rect) (ret []*Rect) {
    // Build a list of all objects we could be in collsion with
    ret = make([]*Rect, 0)

    index := node.index(obj)
    if index != -1 && len(node.Nodes) != 0 {
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