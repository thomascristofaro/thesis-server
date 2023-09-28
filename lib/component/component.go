package component

type PageType int

const (
	PageList PageType = 0
	PageCard          = 1
	Home              = 2
)

type SectionType int

const (
	Repeater  SectionType = 0
	Group                 = 1
	PieChart              = 2
	LineChart             = 3
	SubPage               = 4
)

type FieldType int

const (
	TextType     FieldType = 0
	IntType                = 1
	DecimalType            = 2
	BooleanType            = 3
	DateType               = 4
	TimeType               = 5
	DatetimeType           = 6
)

type Page interface {
	GetSchema() ([]byte, error)
	Get(filters map[string][]string) ([]byte, error)
	Post(body []byte) ([]byte, error)
	Patch(body []byte) ([]byte, error)
	Delete(filters map[string][]string) ([]byte, error)
	GetId() string
}
