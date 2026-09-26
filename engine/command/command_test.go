package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCommand struct {
	Active bool
}

func (c *TestCommand) IsActive() bool {
	return c.Active
}

func TestCommandBasic(t *testing.T) {
	require := require.New(t)
	button := &TestCommand{}

	c := New[int]()
	c.Register(0, button.IsActive)

	c.Update()
	require.False(c.Pressed(0))
	require.False(c.Check(0))
	require.False(c.Released(0))
	require.False(c.CheckAny())

	button.Active = true

	c.Update()
	require.True(c.Pressed(0))
	require.True(c.Check(0))
	require.False(c.Released(0))
	require.True(c.CheckAny())

}

func TestCommandMultiRegister(t *testing.T) {
	require := require.New(t)
	c := New[int]()

	button1 := &TestCommand{}
	button2 := &TestCommand{}

	c.Register(0, button1.IsActive)
	c.Register(0, button2.IsActive)

	c.Update()
	require.False(c.Check(0))

	button1.Active = true
	button2.Active = false
	c.Update()
	require.True(c.Check(0))

	button1.Active = false
	button2.Active = true
	c.Update()
	require.True(c.Check(0))

	button1.Active = false
	button2.Active = false
	c.Update()
	require.False(c.Check(0))

	button3 := &TestCommand{}
	c.Register(0, button3.IsActive)

	button3.Active = true
	c.Update()
	require.True(c.Check(0))
}
