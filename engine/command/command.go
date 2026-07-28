package command

type entry struct {
	checks   []func() bool
	active   bool
	pressed  bool
	released bool
}

type Command[T ~int] struct {
	registry map[T]*entry
}

func New[T ~int]() *Command[T] {
	return &Command[T]{
		registry: make(map[T]*entry),
	}
}

func (c *Command[T]) Register(cmd T, check func() bool) {
	if _, ok := c.registry[cmd]; !ok {
		c.registry[cmd] = &entry{
			checks: make([]func() bool, 0),
		}
	}
	entry := c.registry[cmd]
	entry.checks = append(entry.checks, check)
}

func (c *Command[T]) Update() {
	for _, entry := range c.registry {
		active := false
		for _, check := range entry.checks {
			active = active || check()
		}
		if active {
			if !entry.active {
				entry.pressed = true
			} else {
				entry.pressed = false
			}
			entry.active = true
			entry.released = false
		} else {
			if entry.active {
				entry.released = true
			} else {
				entry.released = false
			}
			entry.active = false
			entry.pressed = false
		}
	}
}

func (c *Command[T]) Pressed(cmd T) bool {
	entry, ok := c.registry[cmd]
	return ok && entry.pressed
}

func (c *Command[T]) Released(cmd T) bool {
	entry, ok := c.registry[cmd]
	return ok && entry.released
}

func (c *Command[T]) Check(cmd T) bool {
	entry, ok := c.registry[cmd]
	return ok && entry.active
}
