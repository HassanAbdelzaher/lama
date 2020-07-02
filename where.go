package lama

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
)

type Where struct {
	Expr  string
	Op    string
	Value interface{}
	Or    []Where
	Args  []sql.NamedArg
	Raw   string
}

func (w *Where) Build() (string, []sql.NamedArg) {
	if w.Args == nil {
		w.Args = make([]sql.NamedArg, 0)
	}
	bnN := "arg"
	if len(w.Expr) > 0 {
		n := strings.TrimLeft(w.Expr, "(")
		n = strings.TrimLeft(n, "/")
		n = strings.TrimLeft(n, ":")
		n = strings.TrimLeft(n, ".")
		n = strings.TrimLeft(n, ",")
		n = strings.TrimLeft(n, "'")
		n = strings.TrimLeft(n, `"`)
		ln := len(n)
		if ln >= 3 {
			chr := []rune(n)
			bnN = string(chr[0:3])
		} else {
			chr := []rune(n)
			bnN = string(chr[0:ln])
		}
	}
	bnN = strings.TrimSpace(bnN)
	if bnN == "" {
		bnN = "arg"
	}
	name := bnN + strconv.Itoa(rand.Int())
	or := ""
	if w.Or != nil {
		for idx, o := range w.Or {
			stm, args := o.Build()
			if idx == 0 {
				or = stm
			} else {
				or = or + " OR " + stm
			}
			w.Args = append(w.Args, args...)
		}
	}
	wh := ""
	if len(w.Expr) > 0 {
		wh = "(" + w.Expr + w.Op + ":" + name + ")"
		w.Args = append(w.Args, sql.NamedArg{Name: name, Value: w.Value})
	}
	if len(w.Raw) > 0 {
		wh = "(" + w.Raw + ")"
	}
	if len(or) > 0 {
		if len(wh) > 0 {
			wh = "(" + wh + " and (" + or + "))"
		} else {
			wh = "(" + or + ")"
		}
	}
	return wh, w.Args
}
