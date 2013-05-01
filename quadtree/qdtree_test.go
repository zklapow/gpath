package quadtree

import(
    "testing"
    "math/rand"
    "time"
)

func BenchmarkInsert(b *testing.B) {
    cur := time.Now()
    rand.Seed(cur.Unix())

    qt := New(0, Rect{X:0, Y:0, Width: 80, Height: 80})
    items := make([]*Rect, b.N)
    for i := 0; i < b.N; i++ {
        // Run the benchmark
        tmp := &Rect{X: rand.Intn(80), Y: rand.Intn(80),
                     Width: rand.Intn(40), Height: rand.Intn(40)}

        qt.Insert(tmp)
        items[i] = tmp
    }
}
