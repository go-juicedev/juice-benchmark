module github.com/go-juicedev/juice-benchmark

go 1.23

toolchain go1.23.2

require (
	github.com/go-juicedev/juice v1.6.7
	github.com/go-sql-driver/mysql v1.7.1
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)

// 使用本地 juice 包
replace github.com/go-juicedev/juice => ../juice
