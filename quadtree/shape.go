package quadtree

type Shape interface {
    In(*QuadTree) int
}

type Rect struct {
    X       int
    Y       int
    Width   int
    Height  int
}

func (rect *Rect) In(node *QuadTree) (ret int) {
    // Calculate what quadrant the object should be placed in
    // -1 indicates the object does not fit entirely within
    // a child node, thus it belong to this node
    ret = -1

    xmid := node.Bounds.X + (node.Bounds.Width / 2)
    ymid := node.Bounds.Y + (node.Bounds.Height / 2)

    top := (rect.Y < ymid && rect.Y + rect.Height < ymid)
    bot := (rect.Y > ymid)

    // Left half
    if rect.X < xmid && (rect.X + rect.Width) < xmid {
        if top {
            ret = 1
        } else if bot {
            ret = 2
        }
    // Right half
    }else if rect.X > xmid {
        if top {
            ret = 0
        } else if bot {
            ret = 3
        }
    }

    return ret
}

type Point struct {
    X int
    Y int
}

func (pt *Point) In(node *QuadTree) (ret int) {
    // If no child contains this point return -1
    // This is not the same as other methods
    // On insert this method will never return -1 because
    // Assuming a node contains a point, one of its children
    // MUST also contain that point
    // -1 then indicated out of bounds of the PARENT node
    ret = -1

    xmid := node.Bounds.X + (node.Bounds.Width / 2)
    ymid := node.Bounds.Y + (node.Bounds.Height / 2)

    // Top half
    if pt.Y < ymid && pt.Y > node.Bounds.Y {
        // Left
        if pt.X < xmid && pt.X > node.Bounds.X {
            ret = 1
        // Right
        } else if pt.X > xmid && pt.X < (node.Bounds.X + node.Bounds.Width) {
            ret = 0
        }
    // Bottom
    } else if pt.Y > ymid && pt.Y < (node.Bounds.Y + node.Bounds.Height) {
        // Left
        if pt.X < xmid && pt.X > node.Bounds.X {
            ret = 2
        // Right
        } else if pt.X > xmid && pt.X < (node.Bounds.X + node.Bounds.Width) {
            ret = 3
        }
    }

    return ret
}

type Line struct {
    Start Point
    End Point
}

func (l *Line) In(node *QuadTree) (ret int) {
    // Leverage the points in method to find the a node containing both points
    // If no node contains both ret=-1
    ret = -1

    for i, n := range node.Nodes {
        // If both points are in the same node the line is that node
        if l.Start.In(n) != -1 && l.End.In(n) != -1 {
            ret = i
            return ret
        }
    }

    return ret
}
