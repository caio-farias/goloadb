package components

type Queue interface {
	add(val string)
	remove(val string)
	next()
}

type CircularQueue struct {
	queue []string
	len   int
	head  string
	tail  string
}

func NewCircularQueue(size int) *CircularQueue {
	return &CircularQueue{
		queue: make([]string, size),
		len:   0,
	}
}

func (c *CircularQueue) add(val string) {
	c.queue = append(c.queue, val)
	c.len = +1
}

func (c *CircularQueue) remove(index int) string {
	removed := c.queue[index]
	for ii := index; ii < c.len-1; ii++ {
		c.queue[ii] = c.queue[ii+1]
	}
	return removed
}

func (c *CircularQueue) next() string {
	next := c.remove(0)
	c.add(next)

	return next
}
