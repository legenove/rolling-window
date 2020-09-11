package windows

// CirCircularQueue
//   capacity is the numbers of this queue's ability (capacity - 1)
type CircularQueue struct {
	capacity int
	elements []float64
	first    int
	end      int
}

// NewCircularQueue create new CircularQueue passing a integer as its size
// and return its pointer
func NewCircularQueue(size int) *CircularQueue {
	cq := CircularQueue{capacity: size + 1, first: 0, end: 0}
	cq.elements = make([]float64, cq.capacity)
	return &cq
}

// IsEmpty return if this queue is empty
func (c CircularQueue) IsEmpty() bool {
	return c.first == c.end
}

func (c CircularQueue) Len() int {
	return (c.end + c.capacity - c.first) % c.capacity
}

// IsFull return if this queue is full
func (c CircularQueue) IsFull() bool {
	return c.first == (c.end+1)%c.capacity
}

// Push pushing a element to this queue
// note: if pushing into a full queue, it will panic
func (c *CircularQueue) Push(e float64) {
	if c.IsFull() {
		panic("Queue is full")
	}
	c.elements[c.end] = e
	c.end = (c.end + 1) % c.capacity
}

func (c *CircularQueue) PushEmpty(i int) {
	if c.Len()+i+1 > c.capacity {
		panic("Queue is full")
	}
	c.end = (c.end + i) % c.capacity
}

// Shift shift a element witch pushed earlist
// note: if will return nil if this queue is empty
func (c *CircularQueue) Shift() (e float64) {
	if c.IsEmpty() {
		return 0
	}
	e = c.elements[c.first]
	c.first = (c.first + 1) % c.capacity
	return
}

func (c *CircularQueue) Pop() (e float64) {
	if c.IsEmpty() {
		return 0
	}
	e = c.elements[(c.end + c.capacity - 1) % c.capacity]
	c.end = (c.end + c.capacity - 1) % c.capacity
	return
}

func (c *CircularQueue) First() float64 {
	if c.IsEmpty() {
		return 0
	}
	return c.elements[c.first]
}

func (c *CircularQueue) Last() float64 {
	if c.IsEmpty() {
		return 0
	}
	return c.elements[(c.end + c.capacity - 1) % c.capacity]
}
