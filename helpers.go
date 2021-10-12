package lama

import (
	"database/sql"
	"fmt"
	"strings"
)

///create wheres
// WhereEqual accept nil and zero values
func Eq(key string, value interface{}) Where {
	if value == nil {
		return Where{Raw: fmt.Sprintf(`%s is null`, key)}
	}
	return Where{Expr: key, Op: "=", Value: value}
}

// WhereEqual accept nil and zero values
func NotEq(key string, value interface{}) Where {
	if value == nil {
		return Where{Raw: fmt.Sprintf(`%s is not null`, key)}
	}
	return Where{Expr: key, Op: "<>", Value: value}
}

func Gt(key string, value interface{}) Where {
	if value == nil {
		return Where{Fake: true}
	}
	return Where{Expr: key, Op: ">", Value: value}
}

func Gte(key string, value interface{}) Where {
	if value == nil {
		return Where{Fake: true}
	}
	return Where{Expr: key, Op: ">=", Value: value}
}

func Lt(key string, value interface{}) Where {
	if value == nil {
		return Where{Fake: true}
	}
	return Where{Expr: key, Op: "<", Value: value}
}

func Lte(key string, value interface{}) Where {
	if value == nil {
		return Where{Fake: true}
	}
	return Where{Expr: key, Op: "<=", Value: value}
}

func In(di Dialect, key string, values ...interface{}) Where {
	if values == nil || len(values) == 0 {
		return Where{Fake: true}
	}
	args := make([]sql.NamedArg, 0)
	ins := make([]string, 0)
	for idx := range values {
		nam := (getArgName(fmt.Sprintf(`%s%d`, key, idx)))
		args = append(args, sql.NamedArg{Name: nam, Value: values[idx]})
		//args[nam] = v
		ins = append(ins, di.BindVarStr(nam))
	}
	//stm := " " + key + " between (" + strings.Join(ins, ",") + ")"
	stm := fmt.Sprintf(`%s in (%s)`, key, strings.Join(ins, ","))
	return Where{Raw: stm, Args: args}
}

func Between(di Dialect, key string, value1 interface{}, value2 interface{}) Where {
	if value1 == nil && value2 == nil {
		return Where{Fake: true}
	}
	if value2 == nil {
		return Gte(key, value1)
	}
	if value1 == nil {
		return Lte(key, value2)
	}
	args := make([]sql.NamedArg, 0)
	nam1 := (getArgName(key))
	nam2 := (getArgName(key + "_"))
	args = append(args, sql.NamedArg{Name: nam1, Value: value1})
	args = append(args, sql.NamedArg{Name: nam2, Value: value2})
	stm := fmt.Sprintf(`%s between %s and %s`, key, di.BindVarStr(nam1), di.BindVarStr(nam2))
	return Where{Raw: stm, Args: args}
}

func Like(key string, value1 interface{}) Where {
	if value1 == nil {
		return Where{Fake: true}
	}
	args := make([]sql.NamedArg, 0)
	stm := fmt.Sprintf(`%s like '%%%s%%'`, key, value1)
	return Where{Raw: stm, Args: args}
}

func EndsWith(key string, value1 string) Where {
	args := make([]sql.NamedArg, 0)
	stm := fmt.Sprintf(`%s like '%%%s'`, key, value1)
	return Where{Raw: stm, Args: args}
}

func StartsWith(key string, value1 string) Where {
	args := make([]sql.NamedArg, 0)
	stm := fmt.Sprintf(`%s like '%s%%'`, key, value1)
	return Where{Raw: stm, Args: args}
}
