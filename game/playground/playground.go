package playground

import (
	"Bobox/game"
	"Bobox/game/levels"
	"errors"
)

type FieldData struct {
	PlayerPosition  game.Vector
	TargetPositions []game.Vector
	Size            game.Vector
	Area            game.LevelArea
}

func New() FieldData {
	return FieldData{}
}

func (f *FieldData) LoadData(levelIndex int) error {
	if levelIndex >= len(levels.Data) {
		return errors.New("level index is out of range")
	}

	originalLevel := levels.Data[levelIndex]
	f.Area = make(game.LevelArea, len(originalLevel))
	for i := range originalLevel {
		f.Area[i] = make(game.LevelLine, len(originalLevel[i]))
		copy(f.Area[i], originalLevel[i])
	}

	f.Size = f.size()

	playerPosition := f.findFields(game.FieldPlayer, true)

	if len(playerPosition) == 0 {
		return errors.New("couldn't find player position")
	}

	f.PlayerPosition = playerPosition[0]

	f.TargetPositions = f.findFields(game.FieldTarget, false)

	if len(f.TargetPositions) == 0 {
		return errors.New("no target fields detected")
	}

	for _, position := range f.TargetPositions {
		f.SetPosition(position, game.FieldEmpty)
	}

	return nil
}

// ObjectFromPosition checks if position is valid for field and it's value
func (f *FieldData) ObjectFromPosition(position game.Vector) int {
	return f.Area[position.Y][position.X]
}

func (f *FieldData) IsValidPosition(position game.Vector) bool {
	return position.X >= 0 && position.Y >= 0 && position.X < f.Size.X && position.Y < f.Size.Y
}

// AnyTargetLeft returns if there is any target left
func (f *FieldData) AnyTargetLeft() bool {
	for i := 0; i < len(f.TargetPositions); i++ {
		if f.ObjectFromPosition(f.TargetPositions[i]) != game.FieldBox {
			return true
		}
	}
	return false
}

func (f *FieldData) SetPosition(position game.Vector, ID int) {
	f.Area[position.Y][position.X] = ID
}

func (f *FieldData) size() game.Vector {
	height := len(f.Area)
	if height == 0 {
		return game.Vector{
			X: 0,
			Y: 0,
		}
	}
	width := len(f.Area[0])

	return game.Vector{
		X: width,
		Y: height,
	}
}

func (f *FieldData) findFields(objectID int, isSingle bool) []game.Vector {
	res := make([]game.Vector, 0)
	for y := 0; y < f.Size.Y; y++ {
		for x := 0; x < f.Size.X; x++ {
			if f.Area[y][x] == objectID {
				res = append(res, game.Vector{
					X: x,
					Y: y,
				})
				if isSingle {
					return res
				}
			}
		}
	}
	return res
}
