package hypr

const (
	USER_TERMINAL = "kitty"
	NEOVIM        = "nvim"
	YOUTUBE       = "YouTube"
	CLASS_SPOTIFY = "Spotify"
	CLASS_FIREFOX = "firefox"
	CLASS_VESKTOP = "vesktop"
)

type Window struct {
	Class string `json:"class"`
	Title string `json:"title"`
}

// func (w Window) IsVim() bool {
// 	words := strings.Split(w.Title, " ")
//
// 	if len(words) < 1 {
// 		return false
// 	}
//
// 	if w.Class == USER_TERMINAL && words[0] == NEOVIM {
// 		return true
// 	}
//
// 	return false
// }
//
// func (w Window) IsSpotify() bool {
// 	return w.Class == CLASS_SPOTIFY
// }
//
// func (w Window) IsYouTube() bool {
// 	if w.Class != CLASS_FIREFOX {
// 		return false
// 	}
//
// 	for _, word := range strings.Split(w.Title, " ") {
// 		if word == YOUTUBE {
// 			return true
// 		}
// 	}
//
// 	return false
// }
//
// func (w Window) IsDisord() bool {
// 	return w.Class == CLASS_VESKTOP
// }
