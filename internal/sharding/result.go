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

package sharding

import "database/sql"

type Result struct {
	err error
	res []sql.Result
}

func (r Result) Err() error {
	return r.err
}

func (r Result) SetErr(err error) Result {
	return Result{res: r.res, err: err}
}

func (r Result) LastInsertId() (int64, error) {
	return r.res[len(r.res)-1].LastInsertId()
}
func (r Result) RowsAffected() (int64, error) {
	var sum int64
	for _, i := range r.res {
		n, err := i.RowsAffected()
		if err != nil {
			return 0, err
		}
		sum += n
	}
	return sum, nil
}

func NewResult(res []sql.Result) Result {
	return Result{res: res}
}
