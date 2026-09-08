package model

type WindowState struct {
	X      int  `json:"x"`
	Y      int  `json:"y"`
	Width  int  `json:"width"`
	Height int  `json:"height"`
	Pinned bool `json:"pinned"`
	Locked bool `json:"locked"`
}
