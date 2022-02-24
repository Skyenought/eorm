// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eql

import (
	"errors"
	"github.com/gotomicro/eql/internal"
	"github.com/valyala/bytebufferpool"
)

// QueryBuilder is used to build a query
type QueryBuilder interface {
	Build() (*Query, error)
}

// Query represents a query
type Query struct {
	SQL  string
	Args []interface{}
}

type builder struct {
	registry MetaRegistry
	dialect  Dialect
	// Use bytebufferpool to reduce memory allocation.
	// After using buffer, it must be put back in bytebufferpool.
	// Call bytebufferpool.Get() to get a buffer, call bytebufferpool.Put() to put buffer back to bytebufferpool.
	buffer  *bytebufferpool.ByteBuffer
	meta    *TableMeta
	args    []interface{}
	aliases map[string]struct{}
}

func (b builder) quote(val string) {
	_ = b.buffer.WriteByte(b.dialect.quote)
	_, _ = b.buffer.WriteString(val)
	_ = b.buffer.WriteByte(b.dialect.quote)
}

func (b builder) space() {
	_ = b.buffer.WriteByte(' ')
}

func (b builder) end() {
	_ = b.buffer.WriteByte(';')
}

func (b builder) comma() {
	_ = b.buffer.WriteByte(',')
}

func (b *builder) parameter(arg interface{}) {
	if b.args == nil {
		// TODO 4 may be not a good number
		b.args = make([]interface{}, 0, 4)
	}
	_ = b.buffer.WriteByte('?')
	b.args = append(b.args, arg)
}

func (b *builder) buildExpr(expr Expr) error {
	switch e := expr.(type) {
	case RawExpr:
		b.buildRawExpr(e)
	case Column:
		if e.name != "" {
			_, ok := b.aliases[e.name]
			if ok {
				b.quote(e.name)
				return nil
			}
			cm, ok := b.meta.fieldMap[e.name]
			if !ok {
				return internal.NewInvalidColumnError(e.name)
			}
			b.quote(cm.columnName)
		}
	case Aggregate:
		if err := b.buildHavingAggregate(e); err != nil {
			return err
		}
	case valueExpr:
		b.parameter(e.val)
	case MathExpr:
		if err := b.buildBinaryExpr(binaryExpr(e)); err != nil {
			return err
		}
	case binaryExpr:
		if err := b.buildBinaryExpr(e); err != nil {
			return err
		}
	case Predicate:
		if err := b.buildBinaryExpr(binaryExpr(e)); err != nil {
			return err
		}
	case nil:
	default:
		return errors.New("unsupported expr")
	}
	return nil
}

func (b *builder) buildPredicates(predicates []Predicate) error {
	p := predicates[0]
	for i := 1; i < len(predicates); i++ {
		p = p.And(predicates[i])
	}
	return b.buildExpr(p)
}

func (b *builder) buildHavingAggregate(aggregate Aggregate) error {
	_, _ = b.buffer.WriteString(aggregate.fn)
	_ = b.buffer.WriteByte('(')
	cMeta, ok := b.meta.fieldMap[aggregate.arg]
	if !ok {
		return internal.NewInvalidColumnError(aggregate.arg)
	}
	b.quote(cMeta.columnName)
	_ = b.buffer.WriteByte(')')
	return nil
}

func (b *builder) buildBinaryExpr(e binaryExpr) error {
	err := b.buildSubExpr(e.left)
	if err != nil {
		return err
	}
	_, _ = b.buffer.WriteString(e.op.text)
	return b.buildSubExpr(e.right)
}

func (b *builder) buildRawExpr(e RawExpr) {
	_, _ = b.buffer.WriteString(e.raw)
	b.args = append(b.args, e.args...)
}

func (b *builder) buildSubExpr(subExpr Expr) error {
	switch r := subExpr.(type) {
	case MathExpr:
		_ = b.buffer.WriteByte('(')
		if err := b.buildBinaryExpr(binaryExpr(r)); err != nil {
			return err
		}
		_ = b.buffer.WriteByte(')')
	case binaryExpr:
		_ = b.buffer.WriteByte('(')
		if err := b.buildBinaryExpr(r); err != nil {
			return err
		}
		_ = b.buffer.WriteByte(')')
	case Predicate:
		_ = b.buffer.WriteByte('(')
		if err := b.buildBinaryExpr(binaryExpr(r)); err != nil {
			return err
		}
		_ = b.buffer.WriteByte(')')
	default:
		if err := b.buildExpr(r); err != nil {
			return err
		}
	}
	return nil
}
