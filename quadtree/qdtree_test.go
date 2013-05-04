package quadtree

import(
    "testing"
    "math/rand"
    "time"
)

const (
    TESTRECTS = 100
    TESTPOINTS = 100
)

func BenchmarkInsert(b *testing.B) {
    // Don't time setup
    b.StopTimer()

    cur := time.Now()
    rand.Seed(cur.Unix())

    qt := New(0, Rect{X:0, Y:0, Width: 800, Height: 800})
    items := make([]*Rect, b.N)

    b.StartTimer()
    for i := 0; i < b.N; i++ {
        // Run the benchmark
        tmp := &Rect{X: rand.Intn(700), Y: rand.Intn(700),
                     Width: rand.Intn(80), Height: rand.Intn(80)}

        qt.Insert(tmp)
        items[i] = tmp
    }
}

func ConfirmCount(qt *QuadTree, items Shapes) bool {
    // Helper function for making sure insert worked properly
    tree := make(Shapes, 0)
    tree = append(tree, qt.Objects...)
    tree = append(tree, qt.childObjects()...)

    return tree.Count() == items.Count()
}

func TestInsertRects(t *testing.T) {
    cur := time.Now()
    rand.Seed(cur.Unix())

    qt := New(0, Rect{X:0, Y:0, Width: 800, Height:800})
    items := make(Shapes, TESTRECTS)
    for i := 0; i < TESTRECTS; i++ {
        tmp := &Rect{X: rand.Intn(720), Y: rand.Intn(720),
                     Width: rand.Intn(80), Height: rand.Intn(80)}
        err := qt.Insert(tmp)
        if err != nil {
            t.Fatal(err)
        }

        items[i] = tmp
    }

    if !ConfirmCount(qt, items) {
        qt.Draw("TestRects.svg")
        t.Fatal("Number of items inserted does not match number created!")
    }

    qt.Draw("TestRects.svg")
}

func TestInsertPoints(t *testing.T) {
    cur := time.Now()
    rand.Seed(cur.Unix())

    qt := New(0, Rect{X:0, Y:0, Width: 800, Height:800})
    items := make(Shapes, TESTPOINTS)
    for i := 0; i < TESTPOINTS; i++ {
        tmp := &Point{X: rand.Intn(800), Y: rand.Intn(800)}

        err := qt.Insert(tmp)
        if err != nil {
            t.Fatal(err)
        }

        items[i] = tmp
    }

    if !ConfirmCount(qt, items) {
        qt.Draw("TestPoints.svg")
        t.Fatal("Number of items inserted does not match number created!")
    }

    qt.Draw("TestPoints.svg")
}
