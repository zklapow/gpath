package quadtree

import "github.com/ajstarks/svgo"

type Drawer interface {
    Draw(*svg.SVG, string)
}

func (rect *Rect) Draw(out *svg.SVG, opt string) {
    out.Rect(rect.X, rect.Y, rect.Width, rect.Height, opt)
}

func (pt *Point) Draw(out *svg.SVG, opt string) {
    // Draw points as circles with r=1
    out.Circle(pt.X, pt.Y, 1, opt)
}

func (l *Line) Draw(out *svg.SVG, opt string) {
    out.Line(l.Start.X, l.Start.Y, l.End.X, l.End.Y, opt)
}
