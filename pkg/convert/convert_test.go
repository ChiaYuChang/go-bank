package convert_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func TestUUIDv4(t *testing.T) {
	testUUID := uuid.New()

	newUUID, err := convert.StrTo(testUUID.String()).UUIDv4()
	require.NoError(t, err)
	require.Equal(t, testUUID, newUUID)

	newUUID = convert.StrTo(testUUID.String()).MustUUIDv4()
	require.Equal(t, testUUID, newUUID)

	randomstr := "b569b0cb8-57b-4d69-a62a-3ba05f653dce"
	_, err = convert.StrTo(randomstr).UUIDv4()
	require.Error(t, err)
}

func TestDecimal(t *testing.T) {
	s := "100.23"
	d1, _ := decimal.NewFromString(s)

	d2, err := convert.StrTo(s).Decimal()
	require.NoError(t, err)

	if !d1.Equal(d2) {
		t.Fatal("convert error")
	}

	d3 := convert.StrTo(s).MustDecimal()
	require.NoError(t, err)

	if !d1.Equal(d3) {
		t.Fatal("convert error")
	}

	_, err = convert.StrTo("f123").Decimal()
	require.Error(t, err)
}
