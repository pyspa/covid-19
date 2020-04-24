// Copyright 2020 pyspa developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

type Prefecture string

const (
	HOKKAIDO  Prefecture = "北海道"
	TOKYO                = "東京都"
	OSAKA                = "大阪府"
	YAMANASHI            = "山梨県"
	UNKNOWN              = "不明"
)

func PrefFromString(s string) Prefecture {
	switch s {
	case "北海道":
		return HOKKAIDO
	case "東京都":
		return TOKYO
	case "大阪府":
		return OSAKA
	case "山梨県":
		return YAMANASHI
	default:
		return UNKNOWN
	}
}
