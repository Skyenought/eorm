// Copyright 2021 ecodeclub
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

package eorm

import (
	"context"
	"database/sql"

	"github.com/ecodeclub/eorm/internal/datasource"
)

// Executor sql 语句执行器
type Executor interface {
	Exec(ctx context.Context) Result
}

// QueryBuilder 普通 sql 构造抽象
type QueryBuilder interface {
	Build() (Query, error)
}

// Session 代表一个抽象的概念，即会话
type Session interface {
	getCore() core
	queryContext(ctx context.Context, query datasource.Query) (*sql.Rows, error)
	execContext(ctx context.Context, query datasource.Query) (sql.Result, error)
}
