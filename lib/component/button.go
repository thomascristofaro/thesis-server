package component

type Button struct {
	Id       string `json:"id"`
	Caption  string `json:"caption"`
	Icon     int    `json:"icon"`
	Url      string `json:"url"`
	function func() `json:"-"`
}

func NewButton(id string, caption string, icon int, url string, function func()) *Button {
	return &Button{
		Id:       id,
		Caption:  caption,
		Url:      url,
		function: function,
	}
}
