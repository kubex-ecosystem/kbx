package events

// Handler is a function that processes an event of type T.
type Handler[T any] func(T) error

// Collection is an ordered, generic hook registry.
type Collection[T any] struct {
	handlers []Handler[T]
}

// Add appends a handler to the collection.
func (c *Collection[T]) Add(h Handler[T]) {
	c.handlers = append(c.handlers, h)
}

// Fire calls every registered handler in order. Returns the first error encountered.
func (c *Collection[T]) Fire(event T) error {
	for _, h := range c.handlers {
		if err := h(event); err != nil {
			return err
		}
	}
	return nil
}
