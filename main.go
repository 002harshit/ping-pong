package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity struct {
	pos rl.Vector2
	vel float32
	pp  int
}

type Ball struct {
	pos rl.Vector2
	vel rl.Vector2
}

const dim = 1280

var LINE_COLOR = rl.Red
var BG_COLOR = rl.SkyBlue
var TEXT_COLOR = rl.Black
var PLAYER_COLOR = rl.Gray
var CLEAR_COLOR = rl.RayWhite

const MAX_VEL = 16

func main() {
	rl.InitWindow(dim, dim, "basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	playerSize := rl.Vector2{X: dim / 32, Y: dim / 6}
	var ballSize float32 = dim / 32
	isStarted := false

	p1 := Entity{
		pos: rl.Vector2{X: 0, Y: 0},
	}
	p2 := Entity{
		pos: rl.Vector2{X: dim - playerSize.X, Y: 0},
	}
	ball := Ball{
		pos: rl.Vector2{X: dim / 2, Y: dim / 2}, vel: rl.Vector2{X: 0, Y: 0},
	}

	for !rl.WindowShouldClose() {
		// EVENTS BEGIN
		{
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

		}
		// EVENTS END
		// UPDATE BEGIN
		{
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
				// panic("PLAYER 2 GAIN ONE POINT")
			} else if ball.pos.X+ballSize > dim+ballSize*4 {
				ball.pos = rl.Vector2{X: dim / 2, Y: dim / 2}
				isStarted = false
				ball.vel = rl.NewVector2(0, 0)
				p1.pp++
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
				if ball.pos.X-ballSize > p1.pos.X+playerSize.X/1.5 {
					ball.pos.X = p1.pos.X + playerSize.X + ballSize + 0.5
					ball.vel.X *= -1
					ball.vel.Y *= -1
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
				if ball.pos.X-ballSize < p2.pos.X+playerSize.X-playerSize.X/1.5 {
					ball.pos.X = p2.pos.X - ballSize - 0.5
					ball.vel.X *= -1
				}
			}
		}
		// UPDATE END

		// DRAW BEGIN
		rl.BeginDrawing()
		{

			rl.ClearBackground(CLEAR_COLOR)
			rl.DrawRectangle(0, 0, dim, dim, BG_COLOR)
			DrawRectCenter(dim/2, dim/2, 10, dim, LINE_COLOR)
			// DrawTextCenter(message, dim/2, dim/2, font_size, TEXT_COLOR)

			rl.DrawRectangleV(p1.pos, playerSize, PLAYER_COLOR)
			rl.DrawRectangleV(p2.pos, playerSize, PLAYER_COLOR)
			rl.DrawCircleV(ball.pos, ballSize, PLAYER_COLOR)
		}
		rl.EndDrawing()
		// DRAW END
	}
	fmt.Println("EXITED SUCESSFULLY")
}
