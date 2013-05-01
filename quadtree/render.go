package quadtree

import (
    "os"
    "github.com/ajstarks/svgo"
    "fmt"
)

func (tree *QuadTree) Draw(fname string) (err error) {
    f, err := os.Create(fname)
    if err != nil {
        return err
    }

    // Close the file later and check for errors
    defer func() {
        if err := f.Close(); err != nil {
            panic(err)
        }
    }()

    // Initialize the svg with width and height of root node
    out := svg.New(f)
    out.Start(tree.Bounds.Width, tree.Bounds.Height)
    defer out.End() // Make sure to finish the document

    DrawNode(out, tree)

    return nil
}

func DrawNode(out *svg.SVG, node *QuadTree) {
    // Draw the nodes bounds
    out.Rect(node.Bounds.X, node.Bounds.Y, node.Bounds.Width, node.Bounds.Height, "fill:none;stroke:red;stroke-width:4;fill-opacity:0")

    // Draw all objects in the node
    for _, obj := range node.Objects {
        if obj != nil {
            out.Rect(obj.X, obj.Y, obj.Width, obj.Height, "fill:none;stroke:black;stroke-width:2;fill-opacity:0")
        }
    }

    // Draw all child nodes if they exist
    if node.Nodes.Count() != 0 {
        for _, child := range node.Nodes {
            DrawNode(out, child)
        }
    }
}

func (node *QuadTree) Print() {
    prefix := make([]byte, node.Level)
    for i := 0; i < node.Level; i++ {
        prefix[i] = byte('\t')
    }
    fmt.Printf("%vLevel %v\n", string(prefix), node.Level)

    for i := 0; i < node.Objects.Count(); i++ {
        fmt.Printf("%vObject[%v]: %v\n", string(prefix), i, node.Objects[i])
    }

    if node.Nodes.Count() != 0 {
        for i := range node.Nodes {
            node.Nodes[i].Print()
        }
    }
}
