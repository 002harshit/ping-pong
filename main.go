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

type GameWinner int

const (
	GAMEWINNER_NONE GameWinner = iota
	GAMEWINNER_ONE
	GAMEWINNER_TWO
)

var LINE_COLOR = rl.NewColor(40, 90, 70, 255)
var BG_COLOR = rl.NewColor(31, 31, 31, 255)
var TEXT_COLOR = rl.NewColor(216, 216, 216, 255)
var PLAYER_COLOR = rl.NewColor(64, 172, 172, 255)
var BALL_COLOR = rl.NewColor(240, 170, 200, 255)
var CLEAR_COLOR = rl.NewColor(216, 216, 216, 255)

func main() {
	rl.InitWindow(0, 0, "ping pong")
	defer rl.CloseWindow()
	var k = rl.GetScreenWidth()
	if k > rl.GetScreenHeight() {
		k = rl.GetScreenHeight()
	}
	k = int(float32(k) * 0.8)
	dim_i := int(k)
	dim := float32(k)
	rl.SetWindowSize(dim_i, dim_i)
	rl.SetTargetFPS(60)

	MAX_VEL := dim / 72
	BALL_START_SPEED := dim / 128
	MENU_FONT_SIZE := dim / 18
	MSG_FONT_SIZE := dim / 24
	SCORE_FONT_SIZE := dim / 12
	PLAYER_INC_SPEED := float32(2.0)
	P_SIZE := rl.Vector2{X: dim / 20, Y: dim / 6.0}
	P_GRITH := dim / 108
	B_SIZE := dim / 56
	START_TIMER := 36
	BLINK_TIMER := 10
	WIN_COUNT := 1
	EPSILON := float32(0.5)

	p1 := Entity{
		pos: rl.Vector2{X: -P_SIZE.X + P_GRITH, Y: 0}, vel: 0, pp: 0, blink: 0,
	}
	p2 := Entity{
		pos: rl.Vector2{X: dim - P_GRITH, Y: 0}, vel: 0, pp: 0, blink: 0,
	}
	ball := Ball{
		pos: rl.Vector2{X: dim / 2, Y: dim / 2}, vel: rl.Vector2{X: 0, Y: 0},
	}
	state := GAMESTATE_MENU
	winner := GAMEWINNER_NONE

	start_rec := rl.NewRectangle(dim/2-MENU_FONT_SIZE*2, dim/2-MENU_FONT_SIZE-MENU_FONT_SIZE-10, MENU_FONT_SIZE*4, MENU_FONT_SIZE*2)
	start_is_hover := false
	exit_rec := rl.NewRectangle(dim/2-MENU_FONT_SIZE*2, dim/2-MENU_FONT_SIZE+MENU_FONT_SIZE+10, MENU_FONT_SIZE*4, MENU_FONT_SIZE*2)
	exit_is_hover := false

	isStarted := false
	frameCount := 0
	for !rl.WindowShouldClose() && state != GAMESTATE_EXIT {

		// EVENTS BEGIN
		if state == GAMESTATE_GAMEPLAY || state == GAMESTATE_GAMEOVER {
			if rl.IsKeyDown(rl.KeyA) {
				if p1.vel < 0 {
					p1.vel = -p1.vel
				}
				p1.vel += PLAYER_INC_SPEED * frame_time()
			}
			if rl.IsKeyDown(rl.KeyD) {
				if p1.vel > 0 {
					p1.vel = -p1.vel
				}
				p1.vel -= PLAYER_INC_SPEED * frame_time()
			}
			if rl.IsKeyDown(rl.KeyLeft) {
				if p2.vel < 0 {
					p2.vel = -p2.vel
				}
				p2.vel += PLAYER_INC_SPEED * frame_time()
			}
			if rl.IsKeyDown(rl.KeyRight) {
				if p2.vel > 0 {
					p2.vel = -p2.vel
				}
				p2.vel -= PLAYER_INC_SPEED * frame_time()
			}
			if rl.IsKeyPressed(rl.KeyEnter) && !isStarted {

				isStarted = true
				if state == GAMESTATE_GAMEOVER {
					state = GAMESTATE_GAMEPLAY
					winner = GAMEWINNER_NONE
					p1.pp = 0
					p2.pp = 0
				}
				randVect := rl.Vector2{
					X: 100,
					Y: float32(rl.GetRandomValue(-50, 50)),
				}
				if rl.GetRandomValue(-2, 2) < 0 {
					randVect.X *= -1
				}

				randVect = rl.Vector2Scale(rl.Vector2Normalize(randVect), BALL_START_SPEED)
				ball.vel = randVect
			}
		} else if state == GAMESTATE_MENU {
			if rl.IsKeyPressed(rl.KeyEnter) {
				state = GAMESTATE_GAMEPLAY
			}
			if rl.IsKeyPressed(rl.KeyQ) {
				state = GAMESTATE_EXIT
			}
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				if start_is_hover {
					state = GAMESTATE_GAMEPLAY
				}
				if exit_is_hover {
					state = GAMESTATE_EXIT
				}
			}
		}
		// EVENTS END
		// UPDATE BEGIN
		if state == GAMESTATE_GAMEPLAY || state == GAMESTATE_GAMEOVER {
			rl.HideCursor()
			// p1
			p1.pos.Y += p1.vel * frame_time()
			if p1.pos.Y < 0 {
				p1.pos.Y = 0
			} else if p1.pos.Y+P_SIZE.Y > dim {
				p1.pos.Y = dim - P_SIZE.Y
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
			} else if p2.pos.Y+P_SIZE.Y > dim {
				p2.pos.Y = dim - P_SIZE.Y
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
			if ball.pos.X-B_SIZE < -B_SIZE*4 {
				isStarted = false
				ball.pos = rl.Vector2{X: dim / 2, Y: dim / 2}
				ball.vel = rl.NewVector2(0, 0)
				p2.pp++
				p2.blink = BLINK_TIMER * 5
				if p2.pp == WIN_COUNT {
					state = GAMESTATE_GAMEOVER
					winner = GAMEWINNER_TWO
				}
			} else if ball.pos.X+B_SIZE > dim+B_SIZE*4 {
				ball.pos = rl.Vector2{X: dim / 2, Y: dim / 2}
				isStarted = false
				ball.vel = rl.NewVector2(0, 0)
				p1.pp++
				p1.blink = BLINK_TIMER * 5
				if p1.pp == WIN_COUNT {
					state = GAMESTATE_GAMEOVER
					winner = GAMEWINNER_ONE
				}
			}
			if ball.pos.Y-B_SIZE < 0 {
				ball.pos.Y = B_SIZE
				ball.vel.Y = -ball.vel.Y
			} else if ball.pos.Y+B_SIZE > dim {
				ball.pos.Y = dim - B_SIZE
				ball.vel.Y = -ball.vel.Y
			}
			// BALL x PLAYER 1 COLLISION
			if isCollidingCircleRec(ball, B_SIZE, p1, P_SIZE) {
				// if   ball.pos.X - ballSize < p1.pos.X + playerSize
				if rl.CheckCollisionCircleLine(ball.pos, B_SIZE, p1.pos, rl.NewVector2(p1.pos.X+P_SIZE.X, p1.pos.Y)) && ball.vel.Y > 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p1.pos.Y - B_SIZE - EPSILON

				}
				if rl.CheckCollisionCircleLine(ball.pos, B_SIZE, rl.NewVector2(p1.pos.X, p1.pos.Y+P_SIZE.Y), rl.Vector2Add(p1.pos, P_SIZE)) && ball.vel.Y < 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p1.pos.Y + P_SIZE.Y + B_SIZE + EPSILON
				}
				if ball.pos.X-B_SIZE > p1.pos.X {
					ball.pos.X = p1.pos.X + P_SIZE.X + B_SIZE + EPSILON
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
			if isCollidingCircleRec(ball, B_SIZE, p2, P_SIZE) {
				if rl.CheckCollisionCircleLine(ball.pos, B_SIZE, p2.pos, rl.NewVector2(p2.pos.X+P_SIZE.X, p2.pos.Y)) && ball.vel.Y > 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p2.pos.Y - B_SIZE - EPSILON
				}
				if rl.CheckCollisionCircleLine(ball.pos, B_SIZE, rl.NewVector2(p2.pos.X, p2.pos.Y+P_SIZE.Y), rl.Vector2Add(p2.pos, P_SIZE)) && ball.vel.Y < 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p2.pos.Y + P_SIZE.Y + B_SIZE + EPSILON
				}
				if ball.pos.X-B_SIZE < p2.pos.X+P_SIZE.X {
					ball.pos.X = p2.pos.X - B_SIZE - EPSILON
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
		} else if state == GAMESTATE_MENU {
			rl.ShowCursor()
			start_is_hover = rl.CheckCollisionPointRec(rl.GetMousePosition(), start_rec)
			exit_is_hover = rl.CheckCollisionPointRec(rl.GetMousePosition(), exit_rec)
		}
		// UPDATE END

		// DRAW BEGIN
		rl.BeginDrawing()
		if state == GAMESTATE_GAMEPLAY || state == GAMESTATE_GAMEOVER {

			rl.ClearBackground(CLEAR_COLOR)
			rl.DrawRectangle(0, 0, int32(dim_i), int32(dim_i), BG_COLOR)

			if p1.blink <= 0 || (frameCount/BLINK_TIMER)%2 == 0 {
				DrawTextCenter(fmt.Sprint(p1.pp), dim_i/4, int(SCORE_FONT_SIZE)*2, int(SCORE_FONT_SIZE), TEXT_COLOR)
				p1.blink = max(p1.blink-1, 0)
			}
			if p2.blink <= 0 || (frameCount/BLINK_TIMER)%2 == 0 {
				DrawTextCenter(fmt.Sprint(p2.pp), 3*dim_i/4, int(SCORE_FONT_SIZE)*2, int(SCORE_FONT_SIZE), TEXT_COLOR)
				p2.blink = max(p2.blink-1, 0)
			}

			rl.DrawRectangleV(p1.pos, P_SIZE, PLAYER_COLOR)
			rl.DrawRectangleV(p2.pos, P_SIZE, PLAYER_COLOR)

			if !isStarted {
				if state == GAMESTATE_GAMEPLAY && (frameCount/START_TIMER)%2 == 0 {
					DrawRectCenter(dim_i/2, dim_i/2, int(MSG_FONT_SIZE*23*0.9), int(MSG_FONT_SIZE*3), rl.Fade(rl.DarkGray, 0.3))
					DrawTextCenter("Press Enter to start!", dim_i/2, dim_i/2, int(MSG_FONT_SIZE), TEXT_COLOR)
				} else if state == GAMESTATE_GAMEOVER {
					DrawRectCenter(dim_i/2, dim_i/2, int(MSG_FONT_SIZE*23*0.9), int(MSG_FONT_SIZE*6), rl.Fade(rl.DarkGray, 0.3))
					DrawTextCenter(fmt.Sprint("Player ", winner, " won the match"), dim_i/2, dim_i/2-int(MSG_FONT_SIZE)-10, int(MSG_FONT_SIZE), TEXT_COLOR)
					DrawTextCenter("Press Enter to restart!", dim_i/2, dim_i/2+int(MSG_FONT_SIZE)+10, int(MSG_FONT_SIZE), TEXT_COLOR)
				}
			} else {
				DrawRectCenter(dim_i/2, dim_i/2, 10, dim_i, LINE_COLOR)
				rl.DrawCircleV(ball.pos, B_SIZE, BALL_COLOR)
			}
		} else if state == GAMESTATE_MENU {
			rl.ClearBackground(rl.Black)
			rl.DrawRectangle(0, 0, int32(dim_i), int32(dim_i), BG_COLOR)

			if start_is_hover {
				rl.DrawRectangleRec(start_rec, rl.DarkGray)
				DrawTextCenter("Start", dim_i/2, dim_i/2-int(MENU_FONT_SIZE)-10, int(MENU_FONT_SIZE), rl.White)
			} else {
				rl.DrawRectangleRec(start_rec, rl.RayWhite)
				DrawTextCenter("Start", dim_i/2, dim_i/2-int(MENU_FONT_SIZE)-10, int(MENU_FONT_SIZE), rl.Black)
			}
			if exit_is_hover {
				rl.DrawRectangleRec(exit_rec, rl.DarkGray)
				DrawTextCenter("Exit", dim_i/2, dim_i/2+int(MENU_FONT_SIZE)+10, int(MENU_FONT_SIZE), rl.White)
			} else {
				rl.DrawRectangleRec(exit_rec, rl.RayWhite)
				DrawTextCenter("Exit", dim_i/2, dim_i/2+int(MENU_FONT_SIZE)+10, int(MENU_FONT_SIZE), rl.Black)
			}

		}
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
		// DRAW END
	}
	fmt.Println("EXITED SUCESSFULLY")
}
