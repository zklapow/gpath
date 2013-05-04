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
    // Any given point must be FULLY contained by a child
    // This function CANNOT return -1
    xmid := node.Bounds.X + (node.Bounds.Width / 2)
    ymid := node.Bounds.Y + (node.Bounds.Height / 2)

    // Top half
    if pt.Y < ymid {
        // Left
        if pt.X < xmid {
            ret = 1
        // Right
        } else {
            ret = 0
        }
    // Bottom
    } else {
        // Left
        if pt.X < xmid {
            ret = 2
        // Right
        } else {
            ret = 3
        }
    }

    return ret
}
