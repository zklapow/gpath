package quadtree

const (
    MAXLEVELS = 5
    MAXOBJECTS = 10
)

type QuadTree struct {
    Bounds  Rect
    Level   int

    Nodes   Trees
    Objects Shapes
}

type Trees []*QuadTree
func (l Trees) Count() int {
    c := 0
    for _, val := range l {
        if val != nil {c++}
    }
    return c
}

type Shapes []Shape
func (l Shapes) Count() int {
    c := 0
    for _, val := range l {
        if val != nil {c++}
    }
    return c
}

func New(level int, bounds Rect) (ret *QuadTree) {
    ret = &QuadTree{Bounds: bounds, Level: level}

    ret.Nodes = make(Trees, 4)
    ret.Objects = make(Shapes, MAXOBJECTS)

    return ret
}


func (node *QuadTree) Clear() {
    // Make a new empty list
    node.Objects = make(Shapes, MAXOBJECTS)

    // If there are no object then there
    // can't be any child nodes
    node.Nodes = make(Trees, 4)
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

func (node *QuadTree) Insert(obj Shape) (err error) {
    if node.Nodes.Count() != 0 {
        index := obj.In(node)

        // This object belongs to a child
        if index != -1 {
            node.Nodes[index].Insert(obj)
            return nil
        }
    }
        
    // Check if its time to split this node
    if node.Objects.Count() >= MAXOBJECTS {
        // Create a new list for any object which need to stay here
        tmp := make(Shapes, node.Objects.Count())

        // Its possible we do not need to make new nodes
        if node.Nodes.Count() == 0 {
            node.Split()
        }

        for i := 0; i < node.Objects.Count(); i++ {
            if node.Objects[i] != nil {
                index := node.Objects[i].In(node)
                if index != -1 {
                    node.Nodes[index].Insert(node.Objects[i])
                } else {
                    tmp[tmp.Count()] = node.Objects[i]
                }
            }
        }

        // We can now figure out where to put our new object
        index := obj.In(node)
        if index != -1 {
            node.Nodes[index].Insert(obj)
        } else {
            if tmp.Count() >= node.Objects.Count() {
                // We are going to need to expand the array
                tmp = append(tmp, Shapes{obj}...)
            } else {
                tmp[tmp.Count()] = obj
            }
        }

        node.Objects = tmp
        return nil
    } else {
        // Our object list has space
        node.Objects[node.Objects.Count()] = obj
    }

    return nil
}

func (node *QuadTree) childObjects() (ret Shapes) {
    // Get and return all child object of a node
    // TODO: This should be concurrent (nodes of larger depths take longer)
    ret = make(Shapes, 0)
    if node.Nodes.Count() != 0 {
        for _, cn := range node.Nodes {
            ret = append(ret, cn.Objects...)
            ret = append(ret, cn.childObjects()...)
        }
    }

    return ret
}

func (node *QuadTree) Retrieve(obj Shape) (ret Shapes) {
    // Build a list of all objects we could be in collsion with
    ret = make(Shapes, 0)

    index := obj.In(node)
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
