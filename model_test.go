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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTagMetaRegistry(t *testing.T) {
	tm := &TestModel{}
	meta, err := defaultMetaRegistry.Register(tm)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, err)
	assert.Equal(t, 4, len(meta.columns))
	assert.Equal(t, 4, len(meta.fieldMap))
	assert.Equal(t, reflect.TypeOf(tm), meta.typ)
	assert.Equal(t, "test_model", meta.tableName)

	idMeta := meta.fieldMap["Id"]
	assert.Equal(t, "id", idMeta.columnName)
	assert.Equal(t, "Id", idMeta.fieldName)
	assert.Equal(t, reflect.TypeOf(int64(0)), idMeta.typ)
	assert.True(t, idMeta.isAutoIncrement)
	assert.True(t, idMeta.isPrimaryKey)

	// 把剩下的三个字段都断言一遍，注意它们的类型，而且它们也不是主键，也不是自增
}
