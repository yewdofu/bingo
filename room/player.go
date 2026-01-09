package room

type Player struct {
	DiscordUserID string
	Name          string
	Color         string
	Progress      [25]bool
	Rating        int
}
