package component

type Button struct {
	Id       string                                `json:"id"`
	Caption  string                                `json:"caption"`
	Icon     int                                   `json:"icon"`
	Url      string                                `json:"url"`
	Function func(Page, map[string][]string) error `json:"-"`
}

func NewButton(id string, caption string, icon int, url string, f func(Page, map[string][]string) error) *Button {
	return &Button{
		Id:       id,
		Caption:  caption,
		Url:      url,
		Function: f,
	}
}
