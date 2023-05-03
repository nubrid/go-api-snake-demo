package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"math/rand"
	"time"
)

// JSON: { x, y }
type move struct {
	X int `json:"x" validate:"required"`
	Y int `json:"y" validate:"required"`
}

// JSON: { moves: [{ x, y }], score, size }
type moveSet struct {
	Moves []*move `json:"moves" validate:"required"`
	Score int     `json:"score" validate:"min=0"`
	Size  int     `json:"size" validate:"required,min=2"`
}

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validateMoveSetStruct(m moveSet) []*errorResponse {
	var errors []*errorResponse

	validate := validator.New()

	err := validate.Struct(m)

	if err != nil {
		// _ = index
		for _, currentErr := range err.(validator.ValidationErrors) {
			var element errorResponse

			// e.g.
			// {
			// 	"FailedField": "MoveSet.Size",
			// 	"Tag": "min",
			// 	"Value": "2"
			// }
			element.FailedField = currentErr.StructNamespace()
			element.Tag = currentErr.Tag()
			element.Value = currentErr.Param()

			errors = append(errors, &element)
		}
	}

	return errors
}

// abs returns the absolute value
func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

func validateMoveSet(m moveSet) error {
	moveSetSize := m.Size

	for i := 0; i < len(m.Moves); i++ {
		currentMove := m.Moves[i]

		isWithinBounds := currentMove.X >= 0 && currentMove.X <= moveSetSize-1 && currentMove.Y >= 0 && currentMove.Y <= moveSetSize-1

		if !isWithinBounds {
			return fmt.Errorf("snake hit the wall: size=%dx%d moves[%d] { %d, %d }", moveSetSize, moveSetSize, i, currentMove.X, currentMove.Y)
		}

		isSecondMoveAndAbove := i > 0

		if isSecondMoveAndAbove {
			prevMove := m.Moves[i-1]

			isPrevAndCurrentMoveSame := prevMove.X == currentMove.X && prevMove.Y == currentMove.Y

			if isPrevAndCurrentMoveSame {
				return fmt.Errorf("snake cannot stop: moves[%d] and moves[%d] { %d, %d }", i, i-1, currentMove.X, currentMove.Y)
			}

			isPrevAndCurrentMoveAdjacent := (abs(currentMove.X-prevMove.X) + abs(currentMove.Y-prevMove.Y)) == 1

			if !isPrevAndCurrentMoveAdjacent {
				return fmt.Errorf("snake cannot move diagonally: moves[%d] { %d, %d } vs moves[%d] { %d, %d }", i, currentMove.X, currentMove.Y, i-1, prevMove.X, prevMove.Y)
			}

			isThirdMoveAndAbove := i > 1

			if isThirdMoveAndAbove {
				prevPrevMove := m.Moves[i-2]

				isMoveReversed := prevPrevMove.X == currentMove.X && prevPrevMove.Y == currentMove.Y

				if isMoveReversed {
					return fmt.Errorf("snake cannot reverse: moves[%d] { %d, %d } vs moves[%d] { %d, %d }", i, currentMove.X, currentMove.Y, i-2, prevPrevMove.X, prevPrevMove.Y)
				}
			}
		}
	}

	return nil
}

func generateRandomFruit(size int, excludeFruit *move) move {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomFruit := move{X: newRand.Intn(size), Y: newRand.Intn(size)}

	if excludeFruit != nil && randomFruit.X == excludeFruit.X && randomFruit.Y == excludeFruit.Y {
		return generateRandomFruit(size, excludeFruit)
	}

	return randomFruit
}

func ValidateMoveSet(c *fiber.Ctx) error {
	var moveSet moveSet

	// { size: 2, moves: [{0, 0}, {1, 0}, ..., {x, y}] }
	if err := c.BodyParser(&moveSet); err != nil {
		return err
	}

	errors := validateMoveSetStruct(moveSet)

	if errors != nil {
		return c.JSON(errors)
	}

	err := validateMoveSet(moveSet)

	if err != nil {
		return c.JSON(fiber.Map{"isValid": false, "message": err.Error()})

	}

	// Initially, snake will always start at {0, 0}. Likewise, new fruit must not spawn here
	prevFruit := move{X: 0, Y: 0}

	hasMoves := len(moveSet.Moves) > 0

	if hasMoves {
		prevFruit = *moveSet.Moves[len(moveSet.Moves)-1]
	}

	newFruit := generateRandomFruit(moveSet.Size, &prevFruit)

	if hasMoves {
		return c.JSON(fiber.Map{"isValid": true, "fruitX": newFruit.X, "fruitY": newFruit.Y, "score": moveSet.Score + 1})
	}

	return c.JSON(fiber.Map{"isValid": true, "fruitX": newFruit.X, "fruitY": newFruit.Y})
}
