package convert

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := s.Int()
	if err != nil {
		return uint32(0), err
	}
	return uint32(v), nil
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

// func (s StrTo) PgNumeric() (pgtype.Numeric, error) {
// 	num := new(pgtype.Numeric)
// 	err := num.Set(s)
// 	return *num, err
// }

// func (s StrTo) MustPgNumeric() pgtype.Numeric {
// 	num, _ := s.PgNumeric()
// 	return num
// }

func (s StrTo) Decimal() (decimal.Decimal, error) {
	return decimal.NewFromString(string(s))
}

func (s StrTo) MustDecimal() decimal.Decimal {
	d, _ := s.Decimal()
	return d
}
