package quadtree

import(
    "testing"
    "math/rand"
    "time"
    "fmt"
)

const (
    TESTRECTS = 100
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

func TestQuadTree(t *testing.T) {
    cur := time.Now()
    rand.Seed(cur.Unix())

    qt := New(0, Rect{X:0, Y:0, Width: 800, Height:800})
    items := make(Rects, TESTRECTS)
    for i := 0; i < TESTRECTS; i++ {
        tmp := &Rect{X: rand.Intn(720), Y: rand.Intn(720),
                     Width: rand.Intn(80), Height: rand.Intn(80)}
        err := qt.Insert(tmp)
        if err != nil {
            t.Fatal(err)
        }

        items[i] = tmp
    }

    tree := make(Rects, 0)
    tree = append(tree, qt.Objects...)
    tree = append(tree, qt.childObjects()...)

    if tree.Count() != items.Count() {
        qt.Draw("TestQuadTreeFail.svg")
        t.Fatalf("%v items were created but %v were inserted!", items.Count(), tree.Count())
    }

    qt.Draw("TestQuadTree.svg")
}
