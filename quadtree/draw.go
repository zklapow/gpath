package quadtree

import "github.com/ajstarks/svgo"

type Drawer interface {
    Draw(*svg.SVG, string)
}

func (rect *Rect) Draw(out *svg.SVG, opt string) {
    out.Rect(rect.X, rect.Y, rect.Width, rect.Height, opt)
}
