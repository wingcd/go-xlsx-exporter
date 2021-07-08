module go-xlsx-protobuf

go 1.16

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.4.0
	github.com/wingcd/go-xlsx-protobuf v0.0.0-20210605005928-d24eefb22e15
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/wingcd/go-xlsx-protobuf => ./
)
