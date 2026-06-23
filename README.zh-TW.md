# Go 時區型別

[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/tz.svg)](https://pkg.go.dev/github.com/min0625/tz)
[![codecov](https://codecov.io/gh/min0625/tz/branch/main/graph/badge.svg)](https://codecov.io/gh/min0625/tz)

[English](./README.md) | **繁體中文**

針對 Go 所提供的 `TimeZone` 型別，以 IANA 時區資料庫為基礎的小型函式庫。

## 特色

- **零值即 UTC** — `TimeZone{}` 代表 UTC；載入 `"UTC"` 或 `""` 必定回傳零值。
- **IANA 名稱** — 接受任何 [IANA 時區資料庫](https://www.iana.org/time-zones) 名稱（例如 `"America/New_York"`）。
- **不支援 `"Local"`** — `"Local"` 時區一律以錯誤拒絕。
- **可比較** — `TimeZone` 為一般結構體；可直接使用 `==` 進行相等判斷。
- **豐富的介面支援** — 實作 `fmt.Stringer`、`sql.Scanner`、`driver.Valuer`、`encoding.TextMarshaler/Unmarshaler` 及 `json.Marshaler/Unmarshaler`。
- **零相依** — 完全建構於標準函式庫之上。

## 安裝

```sh
go get github.com/min0625/tz
```

需要 Go 1.24 或更新版本。

## 快速開始

> **注意**：匯入 `_ "time/tzdata"` 可將 IANA 時區資料庫直接嵌入二進位檔，
> 確保在系統時區資料可能不存在的環境（例如 scratch 或 Alpine 容器）中也能正確運作。

```go
package main

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/min0625/tz"
)

func main() {
	z, err := tz.LoadTimeZone("America/New_York")
	if err != nil {
		panic(err)
	}

	fmt.Println(z)                                      // America/New_York
	fmt.Println(time.Now().In(z.Location()).Location()) // America/New_York

	// UTC 為零值
	utc, _ := tz.LoadTimeZone("UTC")
	fmt.Println(utc == tz.TimeZone{}) // true
}
```

## 文件

- API 參考文件：[pkg.go.dev/github.com/min0625/tz](https://pkg.go.dev/github.com/min0625/tz)
- 可執行範例：[time_zone_example_test.go](./time_zone_example_test.go)

## 授權

請參閱 [LICENSE](./LICENSE)。
