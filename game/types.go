package game

const (
	FieldEmpty = iota
	FieldWall
	FieldPlayer
	FieldBox
	FieldTarget
)

const (
	EmojiLeft    = `⬅️`
	EmojiRight   = `➡️`
	EmojiUp      = `⬆️`
	EmojiDown    = `⬇️`
	EmojiRestart = `🔄`
)

const (
	EmojiEmpty  = `⬛`
	EmojiWall   = `⬜`
	EmojiPlayer = `👷`
	EmojiBox    = `🎁`
	EmojiTarget = `🟥`
)

type LevelLine []int
type LevelArea []LevelLine

type Vector struct {
	X int
	Y int
}

type Game interface {
	LoadLevel(level int) (err error)
	ProcessAction(action string) (err error)
}

type Playground interface {
	LoadData(levelIndex int) error
	AnyTargetLeft() bool
	IsValidPosition(position Vector) bool
	SetPosition(position Vector, ID int)
}
