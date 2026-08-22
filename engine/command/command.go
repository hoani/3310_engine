package command

type entry struct {
	checks   []func() bool
	active   bool
	pressed  bool
	released bool
}

type Command[T ~int] interface {
	Pressed(cmd T) bool
	Released(cmd T) bool
	Check(cmd T) bool
}

type CommandImpl[T ~int] struct {
	registry map[T]*entry
}

func New[T ~int]() *CommandImpl[T] {
	return &CommandImpl[T]{
		registry: make(map[T]*entry),
	}
}

func (c *CommandImpl[T]) Register(cmd T, check func() bool) {
	if _, ok := c.registry[cmd]; !ok {
		c.registry[cmd] = &entry{
			checks: make([]func() bool, 0),
		}
	}
	entry := c.registry[cmd]
	entry.checks = append(entry.checks, check)
}

func (c *CommandImpl[T]) Update() {
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

func (c *CommandImpl[T]) Pressed(cmd T) bool {
	entry, ok := c.registry[cmd]
	return ok && entry.pressed
}

func (c *CommandImpl[T]) Released(cmd T) bool {
	entry, ok := c.registry[cmd]
	return ok && entry.released
}

func (c *CommandImpl[T]) Check(cmd T) bool {
	entry, ok := c.registry[cmd]
	return ok && entry.active
}
