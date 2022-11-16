package convert_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gjerry134679/bank/pkg/convert"
)

func TestString(t *testing.T) {
	var s string = "abc"
	var st convert.StrTo = "abc"

	require.Equal(t, s, st.String())
	require.Equal(t, st, convert.StrTo(s))
}

func TestInt(t *testing.T) {
	var i int = 123

	sti := convert.StrTo(fmt.Sprintf("%d", i))
	v, err := sti.Int()
	require.NoError(t, err)
	require.Equal(t, i, v)

	var s string = "12$"
	sti = convert.StrTo(s)
	_, err = sti.Int()
	require.Error(t, err)

	s = "5123"
	sti = convert.StrTo(s)
	v = sti.MustInt()
	require.Equal(t, 5123, v)
}

func TestUint32(t *testing.T) {
	s := "1234"
	u32, err := convert.StrTo(s).UInt32()
	require.NoError(t, err)
	require.Equal(t, uint32(1234), u32)

	u32 = convert.StrTo(s).MustUInt32()
	require.Equal(t, uint32(1234), u32)

	s = "1df123"
	_, err = convert.StrTo(s).UInt32()
	require.Error(t, err)

	s = "4294967295"
	u32, err = convert.StrTo(s).UInt32()
	require.NoError(t, err)
	require.Equal(t, uint32(4294967295), u32)

	s = "4294967296"
	_, err = convert.StrTo(s).UInt32()
	require.EqualError(t, err, convert.ErrOverFlow.Error())

	s = "-10"
	_, err = convert.StrTo(s).UInt32()
	require.EqualError(t, err, convert.ErrOverFlow.Error())
}
