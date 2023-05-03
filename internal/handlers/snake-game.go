package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nubrid/go-api-snake-demo/internal/utils"

	"math/rand"
	"time"
)

type gameDimension struct {
	Width  int `query:"w" validate:"required,min=2"`
	Height int `query:"h" validate:"required,min=2"`
}

type fruit struct {
	X int `json:"x" validate:"min=0"`
	Y int `json:"y" validate:"min=0"`
}

type snake struct {
	X    int `json:"x" validate:"min=0"`
	Y    int `json:"y" validate:"min=0"`
	VelX int `json:"velX" validate:"oneof=-1 0 1"` // X velocity of the snake (one of -1, 0, 1)
	VelY int `json:"velY" validate:"oneof=-1 0 1"` // Y velocity of the snake (one of -1, 0, 1)
}

type state struct {
	GameID string `json:"gameId" validate:"required,uuid4"`
	Width  int    `json:"width" validate:"required,min=2"`
	Height int    `json:"height" validate:"required,min=2"`
	Score  int    `json:"score" validate:"min=0"`
	Fruit  fruit  `json:"fruit" validate:"required"`
	Snake  snake  `json:"snake" validate:"required"`
}

type velocity struct {
	VelX int `json:"velX" validate:"oneof=-1 0 1"` // X velocity of the snake (one of -1, 0, 1)
	VelY int `json:"velY" validate:"oneof=-1 0 1"` // Y velocity of the snake (one of -1, 0, 1)
}

type moveSet struct {
	State state       `json:"state" validate:"required"`
	Ticks []*velocity `json:"ticks" validate:"required,min=1"`
}

func generateRandomFruit(dimension gameDimension, excludedFruit fruit) fruit {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	newFruit := fruit{X: newRand.Intn(dimension.Height),
		Y: newRand.Intn(dimension.Width)}

	if newFruit.X == excludedFruit.X && newFruit.Y == excludedFruit.Y {
		return generateRandomFruit(dimension, excludedFruit)
	}

	return newFruit
}

func CreateNewGame(c *fiber.Ctx) error {
	var newGameDimension gameDimension

	if err := c.QueryParser(&newGameDimension); err != nil {
		return err
	}

	errors := utils.ValidateStruct(newGameDimension)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	newState := state{
		GameID: uuid.NewString(),
		Width:  newGameDimension.Width,
		Height: newGameDimension.Height,
		// Initially, snake will always start at {0, 0}. Likewise, new fruit must not spawn here
		Fruit: generateRandomFruit(newGameDimension, fruit{X: 0, Y: 0}),
		Snake: snake{VelX: 1},
	}

	return c.JSON(newState)
}

func validateTicks(m moveSet) (*snake, int, error) {
	currentSnake := m.State.Snake

	hasSnakeFoundFruit := false

	for i := 0; i < len(m.Ticks); i++ {
		currentTick := m.Ticks[i]

		hasSnakeStoppedMoving := currentTick.VelX == 0 && currentTick.VelY == 0

		if hasSnakeStoppedMoving {
			return nil, fiber.StatusTeapot, fmt.Errorf("game is over, snake stopped moving")
		}

		hasSnakeMovedDiagonally := utils.Abs(currentTick.VelX) == 1 && utils.Abs(currentTick.VelY) == 1

		if hasSnakeMovedDiagonally {
			return nil, fiber.StatusTeapot, fmt.Errorf("game is over, snake made an invalid move diagonally")
		}

		currentSnake = snake{
			X:    currentSnake.X + currentTick.VelX,
			Y:    currentSnake.Y + currentTick.VelY,
			VelX: currentTick.VelX,
			VelY: currentTick.VelY,
		}

		currentState := m.State
		isWithinBounds := currentSnake.X >= 0 && currentSnake.X <= currentState.Width-1 && currentSnake.Y >= 0 && currentSnake.Y <= currentState.Height-1

		if !isWithinBounds {
			return nil, fiber.StatusTeapot, fmt.Errorf("game is over, snake went out of bounds")
		}

		isSecondTickAndAbove := i > 0

		if isSecondTickAndAbove {
			prevTick := m.Ticks[i-1]

			isMoveReversed := utils.Abs(prevTick.VelX)-utils.Abs(currentTick.VelX) == 0 && utils.Abs(prevTick.VelY)-utils.Abs(currentTick.VelY) == 0

			if isMoveReversed {
				return nil, fiber.StatusTeapot, fmt.Errorf("game is over, snake made an invalid move by reversing")
			}
		}

		// possible that snake found fruit, but still moving, i.e. more ticks
		if !hasSnakeFoundFruit {
			hasSnakeFoundFruit = currentSnake.X == m.State.Fruit.X && currentSnake.Y == m.State.Fruit.Y
		}
	}

	if !hasSnakeFoundFruit {
		return nil, fiber.StatusNotFound, fmt.Errorf("fruit not found, the ticks do not lead the snake to the fruit position")
	}

	return &currentSnake, fiber.StatusOK, nil
}

func ValidateMoveSet(c *fiber.Ctx) error {
	var moveSet moveSet

	if err := c.BodyParser(&moveSet); err != nil {
		return err
	}

	errors := utils.ValidateStruct(moveSet)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	newSnakeLocation, errorStatus, err := validateTicks(moveSet)

	if err != nil {
		return c.Status(errorStatus).JSON(fiber.Map{"error": err.Error()})
	}

	currentState := moveSet.State

	newState := state{
		GameID: currentState.GameID,
		Width:  currentState.Width,
		Height: currentState.Height,
		Score:  currentState.Score + 1,
		Fruit:  generateRandomFruit(gameDimension{Width: currentState.Width, Height: currentState.Height}, fruit{X: newSnakeLocation.X, Y: newSnakeLocation.Y}),
		Snake:  *newSnakeLocation,
	}

	return c.JSON(newState)
}
