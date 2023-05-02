package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// JSON: { x, y }
type Move struct {
	X int `json:"x" validate:"required"`
	Y int `json:"y" validate:"required"`
}

// JSON: { size, moves: [{ x, y }] }
type MoveSet struct {
	Size  int     `json:"size" validate:"required,min=2"`
	Moves []*Move `json:"moves" validate:"required,min=2"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validateMoveSetStruct(m MoveSet) []*ErrorResponse {
	var errors []*ErrorResponse

	validate := validator.New()

	err := validate.Struct(m)

	if err != nil {
		// _ = index
		for _, currentErr := range err.(validator.ValidationErrors) {
			var element ErrorResponse

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

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func validateMoveSet(m MoveSet) (bool, error) {
	moveSetSize := m.Size

	for i := 0; i < len(m.Moves); i++ {
		currentMove := m.Moves[i]

		isWithinBounds := currentMove.X >= 0 && currentMove.X <= moveSetSize-1 && currentMove.Y >= 0 && currentMove.Y <= moveSetSize-1

		if !isWithinBounds {
			return false, fmt.Errorf("snake hit the wall: size=%dx%d moves[%d] { %d, %d }", moveSetSize, moveSetSize, i, currentMove.X, currentMove.Y)
		}

		isSecondMoveAndAbove := i > 0

		if isSecondMoveAndAbove {
			prevMove := m.Moves[i-1]

			isPrevAndCurrentMoveSame := prevMove.X == currentMove.X && prevMove.Y == currentMove.Y

			if isPrevAndCurrentMoveSame {
				return false, fmt.Errorf("snake cannot stop: moves[%d] and moves[%d] { %d, %d }", i, i-1, currentMove.X, currentMove.Y)
			}

			isPrevAndCurrentMoveAdjacent := (Abs(currentMove.X-prevMove.X) + Abs(currentMove.Y-prevMove.Y)) == 1

			if !isPrevAndCurrentMoveAdjacent {
				return false, fmt.Errorf("snake cannot move diagonally: moves[%d] { %d, %d } vs moves[%d] { %d, %d }", i, currentMove.X, currentMove.Y, i-1, prevMove.X, prevMove.Y)
			}

			isThirdMoveAndAbove := i > 1

			if isThirdMoveAndAbove {
				prevPrevMove := m.Moves[i-2]

				isMoveReversed := prevPrevMove.X == currentMove.X && prevPrevMove.Y == currentMove.Y

				if isMoveReversed {
					return false, fmt.Errorf("snake cannot reverse: moves[%d] { %d, %d } vs moves[%d] { %d, %d }", i, currentMove.X, currentMove.Y, i-2, prevPrevMove.X, prevPrevMove.Y)
				}
			}
		}
	}

	return true, nil
}

func ValidateMoveSet(c *fiber.Ctx) error {
	var moveSet MoveSet

	// { size: 2, moves: [{0, 0}, {1, 0}, ..., {x, y}] }
	if err := c.BodyParser(&moveSet); err != nil {
		return err
	}

	errors := validateMoveSetStruct(moveSet)

	if errors != nil {
		return c.JSON(errors)
	}

	isValid, err := validateMoveSet(moveSet)

	if err != nil {
		return c.JSON(fiber.Map{"isValid": isValid, "message": err.Error()})

	}

	return c.JSON(fiber.Map{"isValid": isValid})
}

// func GetAllProducts(c *fiber.Ctx) error {
// 	var products []*Product

// 	// await cur.forEach((p) => {
// 	// 	products.push(p);
// 	// })
// 	for cur.Next(context.TODO()) {
// 		var p Product

// 		err := cur.Decode(&p)

// 		if err != nil {
// 			return err
// 		}

// 		products = append(products, &p)
// 	}

// 	return c.JSON(products)
// }
