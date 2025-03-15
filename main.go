package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity struct {
	pos   rl.Vector2
	vel   float32
	pp    int
	blink int
}

type Ball struct {
	pos rl.Vector2
	vel rl.Vector2
}

type GameState int

const (
	GAMESTATE_MENU GameState = iota
	GAMESTATE_GAMEPLAY
	GAMESTATE_GAMEOVER
	GAMESTATE_EXIT
)

const dim = 1280

var LINE_COLOR = rl.Red
var BG_COLOR = rl.SkyBlue
var TEXT_COLOR = rl.Black
var PLAYER_COLOR = rl.Gray
var CLEAR_COLOR = rl.RayWhite

const MAX_VEL = 16

func main() {
	rl.InitWindow(dim, dim, "ping pong")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	playerSize := rl.Vector2{X: dim / 32, Y: dim / 6}
	var ballSize float32 = dim / 32
	isStarted := false
	frameCount := 0

	p1 := Entity{
		pos: rl.Vector2{X: 0, Y: 0},
	}
	p2 := Entity{
		pos: rl.Vector2{X: dim - playerSize.X, Y: 0},
	}
	ball := Ball{
		pos: rl.Vector2{X: dim / 2, Y: dim / 2}, vel: rl.Vector2{X: 0, Y: 0},
	}
	gamestate := GAMESTATE_MENU

	start_rec := rl.NewRectangle(dim/2-200, dim/2-60-50, 400, 100)
	start_is_hover := false
	exit_rec := rl.NewRectangle(dim/2-200, dim/2+60-50, 400, 100)
	exit_is_hover := false

	for !rl.WindowShouldClose() && gamestate != GAMESTATE_EXIT {
		// EVENTS BEGIN
		if gamestate == GAMESTATE_GAMEPLAY {
			if rl.IsKeyDown(rl.KeyA) {
				if p1.vel < 0 {
					p1.vel = -p1.vel
				}
				p1.vel += 2 * frame_time()
			}
			if rl.IsKeyDown(rl.KeyD) {
				if p1.vel > 0 {
					p1.vel = -p1.vel
				}
				p1.vel -= 2 * frame_time()
			}
			if rl.IsKeyDown(rl.KeyLeft) {
				if p2.vel < 0 {
					p2.vel = -p2.vel
				}
				p2.vel += 2 * frame_time()
			}
			if rl.IsKeyDown(rl.KeyRight) {
				if p2.vel > 0 {
					p2.vel = -p2.vel
				}
				p2.vel -= 2 * frame_time()
			}
			if rl.IsKeyPressed(rl.KeyEnter) && !isStarted {
				startSpeed := 10
				isStarted = true
				randVect := rl.Vector2{
					X: 100,
					Y: float32(rl.GetRandomValue(-50, 50)),
				}
				if rl.GetRandomValue(-2, 2) < 0 {
					randVect.X *= -1
				}

				randVect = rl.Vector2Scale(rl.Vector2Normalize(randVect), float32(startSpeed))
				ball.vel = randVect
			}
		} else if gamestate == GAMESTATE_MENU {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				if start_is_hover {
					// panic("MAA PHIRSE CHUD GYI")
					gamestate = GAMESTATE_GAMEPLAY
				}
				if exit_is_hover {
					gamestate = GAMESTATE_EXIT
				}
			}
		}
		// EVENTS END
		// UPDATE BEGIN
		if gamestate == GAMESTATE_GAMEPLAY {
			// p1
			p1.pos.Y += p1.vel * frame_time()
			if p1.pos.Y < 0 {
				p1.pos.Y = 0
			} else if p1.pos.Y+playerSize.Y > dim {
				p1.pos.Y = dim - playerSize.Y
			}
			p1.vel *= 0.9
			if p1.vel < -MAX_VEL {
				p1.vel = -MAX_VEL
			} else if p1.vel > MAX_VEL {
				p1.vel = MAX_VEL
			}

			// p2
			p2.pos.Y += p2.vel * frame_time()
			if p2.pos.Y < 0 {
				p2.pos.Y = 0
			} else if p2.pos.Y+playerSize.Y > dim {
				p2.pos.Y = dim - playerSize.Y
			}
			p2.vel *= 0.9
			if p2.vel < -MAX_VEL {
				p2.vel = -MAX_VEL
			} else if p2.vel > MAX_VEL {
				p2.vel = MAX_VEL
			}

			// ball
			ball.pos.X += ball.vel.X * frame_time()
			ball.pos.Y += ball.vel.Y * frame_time()
			if ball.pos.X-ballSize < -ballSize*4 {
				isStarted = false
				ball.pos = rl.Vector2{X: dim / 2, Y: dim / 2}
				ball.vel = rl.NewVector2(0, 0)
				p2.pp++
				p2.blink = 60
				// panic("PLAYER 2 GAIN ONE POINT")
			} else if ball.pos.X+ballSize > dim+ballSize*4 {
				ball.pos = rl.Vector2{X: dim / 2, Y: dim / 2}
				isStarted = false
				ball.vel = rl.NewVector2(0, 0)
				p1.pp++
				p1.blink = 20
				// panic("PLAYER 1 GAIN ONE POINT")
			}
			if ball.pos.Y-ballSize < 0 {
				ball.pos.Y = ballSize
				ball.vel.Y = -ball.vel.Y
			} else if ball.pos.Y+ballSize > dim {
				ball.pos.Y = dim - ballSize
				ball.vel.Y = -ball.vel.Y
			}
			// BALL x PLAYER 1 COLLISION
			if isCollidingCircleRec(ball, ballSize, p1, playerSize) {
				// if   ball.pos.X - ballSize < p1.pos.X + playerSize
				if rl.CheckCollisionCircleLine(ball.pos, ballSize, p1.pos, rl.NewVector2(p1.pos.X+playerSize.X, p1.pos.Y)) && ball.vel.Y > 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p1.pos.Y - ballSize - 0.5

				}
				if rl.CheckCollisionCircleLine(ball.pos, ballSize, rl.NewVector2(p1.pos.X, p1.pos.Y+playerSize.Y), rl.Vector2Add(p1.pos, playerSize)) && ball.vel.Y < 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p1.pos.Y + playerSize.Y + ballSize + 0.5
				}
				if ball.pos.X-ballSize > p1.pos.X {
					ball.pos.X = p1.pos.X + playerSize.X + ballSize + 0.5
					ball.vel.X *= -1
					ball.vel.Y *= -1
					// incrementing some ball velocity based on player vel
					v := p1.vel * 0.1
					if v < 0 {
						v *= -1
					}
					ball.vel = rl.Vector2Add(ball.vel, rl.Vector2Scale(rl.Vector2Normalize(ball.vel), v))
				}
			}

			// BALL x PLAYER 2 COLLISION
			if isCollidingCircleRec(ball, ballSize, p2, playerSize) {
				if rl.CheckCollisionCircleLine(ball.pos, ballSize, p2.pos, rl.NewVector2(p2.pos.X+playerSize.X, p2.pos.Y)) && ball.vel.Y > 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p2.pos.Y - ballSize - 0.5
				}
				if rl.CheckCollisionCircleLine(ball.pos, ballSize, rl.NewVector2(p2.pos.X, p2.pos.Y+playerSize.Y), rl.Vector2Add(p2.pos, playerSize)) && ball.vel.Y < 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p2.pos.Y + playerSize.Y + ballSize + 0.5
				}
				if ball.pos.X-ballSize < p2.pos.X+playerSize.X {
					ball.pos.X = p2.pos.X - ballSize - 0.5
					ball.vel.X *= -1
					// incrementing some ball velocity based on player vel
					v := p1.vel * 0.1
					if v < 0 {
						v *= -1
					}
					ball.vel = rl.Vector2Add(ball.vel, rl.Vector2Scale(rl.Vector2Normalize(ball.vel), v))
				}
			}
			frameCount++
		} else if gamestate == GAMESTATE_MENU {
			start_is_hover = rl.CheckCollisionPointRec(rl.GetMousePosition(), start_rec)
			exit_is_hover = rl.CheckCollisionPointRec(rl.GetMousePosition(), exit_rec)
		}
		// UPDATE END

		// DRAW BEGIN
		rl.BeginDrawing()
		if gamestate == GAMESTATE_GAMEPLAY {

			rl.ClearBackground(CLEAR_COLOR)
			rl.DrawRectangle(0, 0, dim, dim, BG_COLOR)

			if p1.blink <= 0 || (frameCount/10)%2 == 0 {
				DrawTextCenter(fmt.Sprint(p1.pp), dim/4, 64, 128, TEXT_COLOR)
				p1.blink = max(p1.blink-1, 0)
			}
			if p2.blink <= 0 || (frameCount/10)%2 == 0 {
				DrawTextCenter(fmt.Sprint(p2.pp), 3*dim/4, 64, 128, TEXT_COLOR)
				p2.blink = max(p2.blink-1, 0)
			}

			DrawRectCenter(dim/2, dim/2, 5, dim, LINE_COLOR)

			rl.DrawRectangleV(p1.pos, playerSize, PLAYER_COLOR)
			rl.DrawRectangleV(p2.pos, playerSize, PLAYER_COLOR)

			if !isStarted {
				if (frameCount/60)%2 == 0 {
					DrawTextCenter("Press Enter to start!", dim/2, dim/2, 72, TEXT_COLOR)
				}
			} else {

				rl.DrawCircleV(ball.pos, ballSize, PLAYER_COLOR)
			}
		} else if gamestate == GAMESTATE_MENU {
			rl.ClearBackground(rl.Black)
			rl.DrawRectangle(0, 0, dim, dim, BG_COLOR)

			if start_is_hover {
				rl.DrawRectangleRec(start_rec, rl.DarkGray)
				DrawTextCenter("Start", dim/2, dim/2-60, 72, rl.White)
			} else {
				rl.DrawRectangleRec(start_rec, rl.RayWhite)
				DrawTextCenter("Start", dim/2, dim/2-60, 72, TEXT_COLOR)
			}
			if exit_is_hover {
				rl.DrawRectangleRec(exit_rec, rl.DarkGray)
				DrawTextCenter("Exit", dim/2, dim/2+60, 72, rl.White)
			} else {
				rl.DrawRectangleRec(exit_rec, rl.RayWhite)
				DrawTextCenter("Exit", dim/2, dim/2+60, 72, TEXT_COLOR)
			}

		}
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
		// DRAW END
	}
	fmt.Println("EXITED SUCESSFULLY")
}
