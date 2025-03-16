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

func main() {
	rl.InitWindow(128, 128, "ping pong")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	MX := rl.GetMonitorWidth(rl.GetCurrentMonitor())
	MY := rl.GetMonitorHeight(rl.GetCurrentMonitor())
	k := int(min(0.8*float64(MX), 0.8*float64(MY)))
	rl.SetWindowSize(k, k)
	rl.SetWindowPosition(MX/2-k/2, MY/2-k/2)
	rl.SetTargetFPS(60)

	// DIMENTIONS
	DIM_I := rl.GetScreenWidth()
	DIM := float32(DIM_I)

	// GAME VALS
	maxVel := DIM / 72
	ballStartSpeed := DIM / 128
	menuFontSize := DIM / 18
	msgFontSize := DIM / 24
	scoreFontSize := DIM / 12
	pIncSpeed := float32(2.0)
	pSize := rl.Vector2{X: DIM / 20, Y: DIM / 6.0}
	pGrith := DIM / 108
	bSize := DIM / 56
	startBtn := rl.NewRectangle(DIM/2-menuFontSize*2, DIM/2-menuFontSize-menuFontSize-10, menuFontSize*4, menuFontSize*2)
	exitBtn := rl.NewRectangle(DIM/2-menuFontSize*2, DIM/2-menuFontSize+menuFontSize+10, menuFontSize*4, menuFontSize*2)

	// CONST
	const START_TIMER = 36
	const BLINK_TIMER = 10
	const WIN_COUNT = 5
	const EPSILON = 0.5

	// COLORS
	lineCol := rl.NewColor(40, 90, 70, 255)
	bgCol := rl.NewColor(31, 31, 31, 255)
	textCol := rl.NewColor(216, 216, 216, 255)
	playerCol := rl.NewColor(64, 172, 172, 255)
	ballCol := rl.NewColor(240, 170, 200, 255)
	clearCol := rl.NewColor(216, 216, 216, 255)

	// MUSIC AND SFX
	clickSound := rl.LoadSound("./resources/click.mp3")
	hitSound := rl.LoadSound("./resources/hit.mp3")
	menuSound := rl.LoadSound("./resources/menu.mp3")
	bgmSounds := []rl.Sound{
		rl.LoadSound("./resources/bgm1.mp3"),
		rl.LoadSound("./resources/bgm2.mp3"),
		rl.LoadSound("./resources/bgm3.mp3"),
	}
	current_bgm := 1

	rl.SetSoundVolume(clickSound, 0.3)
	rl.SetSoundVolume(hitSound, 0.2)
	rl.SetSoundVolume(menuSound, 0.3)

	p1 := Entity{
		pos: rl.Vector2{X: -pSize.X + pGrith, Y: 0}, vel: 0, pp: 0, blink: 0,
	}
	p2 := Entity{
		pos: rl.Vector2{X: DIM - pGrith, Y: 0}, vel: 0, pp: 0, blink: 0,
	}
	ball := Ball{
		pos: rl.Vector2{X: DIM / 2, Y: DIM / 2}, vel: rl.Vector2{X: 0, Y: 0},
	}
	state := GAMESTATE_MENU
	winner := GAMEWINNER_NONE

	start_btn_hover := false
	exit_btn_hover := false

	started := false
	frame_count := 0

	for !rl.WindowShouldClose() && state != GAMESTATE_EXIT {

		// EVENTS BEGIN
		if state == GAMESTATE_GAMEPLAY || state == GAMESTATE_GAMEOVER {
			if rl.IsKeyDown(rl.KeyA) {
				if p1.vel < 0 {
					p1.vel = -p1.vel
				}
				p1.vel += pIncSpeed * frame_time()
			}
			if rl.IsKeyDown(rl.KeyD) {
				if p1.vel > 0 {
					p1.vel = -p1.vel
				}
				p1.vel -= pIncSpeed * frame_time()
			}
			if rl.IsKeyDown(rl.KeyLeft) {
				if p2.vel < 0 {
					p2.vel = -p2.vel
				}
				p2.vel += pIncSpeed * frame_time()
			}
			if rl.IsKeyDown(rl.KeyRight) {
				if p2.vel > 0 {
					p2.vel = -p2.vel
				}
				p2.vel -= pIncSpeed * frame_time()
			}
			if rl.IsKeyPressed(rl.KeyEnter) && !started {
				rl.PlaySound(clickSound)
				started = true
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

				randVect = rl.Vector2Scale(rl.Vector2Normalize(randVect), ballStartSpeed)
				ball.vel = randVect
			}
		} else if state == GAMESTATE_MENU {
			if rl.IsKeyPressed(rl.KeyEnter) {
				rl.PlaySound(clickSound)
				state = GAMESTATE_GAMEPLAY
			}

			if rl.IsKeyPressed(rl.KeyQ) {
				state = GAMESTATE_EXIT
			}
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				rl.PlaySound(clickSound)
				if start_btn_hover {
					state = GAMESTATE_GAMEPLAY
				}
				if exit_btn_hover {
					state = GAMESTATE_EXIT
				}
			}
		}
		// EVENTS END
		// UPDATE BEGIN
		if state == GAMESTATE_GAMEPLAY || state == GAMESTATE_GAMEOVER {
			rl.HideCursor()
			prev_bgm := current_bgm - 1
			if prev_bgm < 0 {
				prev_bgm = len(bgmSounds) - 1
			}
			if !rl.IsSoundPlaying(bgmSounds[prev_bgm]) {
				rl.StopSound(menuSound)
				for _, bgm := range bgmSounds {
					rl.StopSound(bgm)
				}
				rl.SetSoundVolume(bgmSounds[current_bgm], 0.2)
				rl.PlaySound(bgmSounds[current_bgm])
				current_bgm = (current_bgm + 1) % len(bgmSounds)
			}
			// p1
			p1.pos.Y += p1.vel * frame_time()
			if p1.pos.Y < 0 {
				p1.pos.Y = 0
			} else if p1.pos.Y+pSize.Y > DIM {
				p1.pos.Y = DIM - pSize.Y
			}
			p1.vel *= 0.9
			if p1.vel < -maxVel {
				p1.vel = -maxVel
			} else if p1.vel > maxVel {
				p1.vel = maxVel
			}

			// p2
			p2.pos.Y += p2.vel * frame_time()
			if p2.pos.Y < 0 {
				p2.pos.Y = 0
			} else if p2.pos.Y+pSize.Y > DIM {
				p2.pos.Y = DIM - pSize.Y
			}
			p2.vel *= 0.9
			if p2.vel < -maxVel {
				p2.vel = -maxVel
			} else if p2.vel > maxVel {
				p2.vel = maxVel
			}

			// ball
			ball.pos.X += ball.vel.X * frame_time()
			ball.pos.Y += ball.vel.Y * frame_time()
			if ball.pos.X-bSize < -bSize*4 {
				started = false
				ball.pos = rl.Vector2{X: DIM / 2, Y: DIM / 2}
				ball.vel = rl.NewVector2(0, 0)
				p2.pp++
				p2.blink = BLINK_TIMER * 5
				if p2.pp == WIN_COUNT {
					state = GAMESTATE_GAMEOVER
					winner = GAMEWINNER_TWO
				}
			} else if ball.pos.X+bSize > DIM+bSize*4 {
				ball.pos = rl.Vector2{X: DIM / 2, Y: DIM / 2}
				started = false
				ball.vel = rl.NewVector2(0, 0)
				p1.pp++
				p1.blink = BLINK_TIMER * 5
				if p1.pp == WIN_COUNT {
					state = GAMESTATE_GAMEOVER
					winner = GAMEWINNER_ONE
				}
			}
			if ball.pos.Y-bSize < 0 {
				rl.PlaySound(hitSound)
				ball.pos.Y = bSize
				ball.vel.Y = -ball.vel.Y
			} else if ball.pos.Y+bSize > DIM {
				rl.PlaySound(hitSound)
				ball.pos.Y = DIM - bSize
				ball.vel.Y = -ball.vel.Y
			}
			// BALL x PLAYER 1 COLLISION
			if isCollidingCircleRec(ball, bSize, p1, pSize) {
				rl.PlaySound(hitSound)
				// if   ball.pos.X - ballSize < p1.pos.X + playerSize
				if rl.CheckCollisionCircleLine(ball.pos, bSize, p1.pos, rl.NewVector2(p1.pos.X+pSize.X, p1.pos.Y)) && ball.vel.Y > 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p1.pos.Y - bSize - EPSILON

				}
				if rl.CheckCollisionCircleLine(ball.pos, bSize, rl.NewVector2(p1.pos.X, p1.pos.Y+pSize.Y), rl.Vector2Add(p1.pos, pSize)) && ball.vel.Y < 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p1.pos.Y + pSize.Y + bSize + EPSILON
				}
				if ball.pos.X-bSize > p1.pos.X {
					ball.pos.X = p1.pos.X + pSize.X + bSize + EPSILON
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
			if isCollidingCircleRec(ball, bSize, p2, pSize) {
				rl.PlaySound(hitSound)
				if rl.CheckCollisionCircleLine(ball.pos, bSize, p2.pos, rl.NewVector2(p2.pos.X+pSize.X, p2.pos.Y)) && ball.vel.Y > 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p2.pos.Y - bSize - EPSILON
				}
				if rl.CheckCollisionCircleLine(ball.pos, bSize, rl.NewVector2(p2.pos.X, p2.pos.Y+pSize.Y), rl.Vector2Add(p2.pos, pSize)) && ball.vel.Y < 0 {
					ball.vel.Y *= -1
					ball.pos.Y = p2.pos.Y + pSize.Y + bSize + EPSILON
				}
				if ball.pos.X-bSize < p2.pos.X+pSize.X {
					ball.pos.X = p2.pos.X - bSize - EPSILON
					ball.vel.X *= -1
					// incrementing some ball velocity based on player vel
					v := p1.vel * 0.1
					if v < 0 {
						v *= -1
					}
					ball.vel = rl.Vector2Add(ball.vel, rl.Vector2Scale(rl.Vector2Normalize(ball.vel), v))
				}
			}
			frame_count++
		} else if state == GAMESTATE_MENU {
			if !rl.IsSoundPlaying(menuSound) {
				rl.PlaySound(menuSound)
			}
			rl.ShowCursor()
			start_btn_hover = rl.CheckCollisionPointRec(rl.GetMousePosition(), startBtn)
			exit_btn_hover = rl.CheckCollisionPointRec(rl.GetMousePosition(), exitBtn)
		}
		// UPDATE END

		// DRAW BEGIN
		rl.BeginDrawing()
		if state == GAMESTATE_GAMEPLAY || state == GAMESTATE_GAMEOVER {

			rl.ClearBackground(clearCol)
			rl.DrawRectangle(0, 0, int32(DIM_I), int32(DIM_I), bgCol)

			if p1.blink <= 0 || (frame_count/BLINK_TIMER)%2 == 0 {
				DrawTextCenter(fmt.Sprint(p1.pp), DIM_I/4, int(scoreFontSize)*2, int(scoreFontSize), textCol)
				p1.blink = max(p1.blink-1, 0)
			}
			if p2.blink <= 0 || (frame_count/BLINK_TIMER)%2 == 0 {
				DrawTextCenter(fmt.Sprint(p2.pp), 3*DIM_I/4, int(scoreFontSize)*2, int(scoreFontSize), textCol)
				p2.blink = max(p2.blink-1, 0)
			}

			rl.DrawRectangleV(p1.pos, pSize, playerCol)
			rl.DrawRectangleV(p2.pos, pSize, playerCol)

			if !started {
				if state == GAMESTATE_GAMEPLAY && (frame_count/START_TIMER)%2 == 0 {
					DrawRectCenter(DIM_I/2, DIM_I/2, int(msgFontSize*23*0.9), int(msgFontSize*3), rl.Fade(rl.DarkGray, 0.3))
					DrawTextCenter("Press Enter to start!", DIM_I/2, DIM_I/2, int(msgFontSize), textCol)
				} else if state == GAMESTATE_GAMEOVER {
					DrawRectCenter(DIM_I/2, DIM_I/2, int(msgFontSize*23*0.9), int(msgFontSize*6), rl.Fade(rl.DarkGray, 0.3))
					DrawTextCenter(fmt.Sprint("Player ", winner, " won the match"), DIM_I/2, DIM_I/2-int(msgFontSize)-10, int(msgFontSize), textCol)
					DrawTextCenter("Press Enter to restart!", DIM_I/2, DIM_I/2+int(msgFontSize)+10, int(msgFontSize), textCol)
				}
			} else {
				DrawRectCenter(DIM_I/2, DIM_I/2, 10, DIM_I, lineCol)
				rl.DrawCircleV(ball.pos, bSize, ballCol)
			}
		} else if state == GAMESTATE_MENU {
			rl.ClearBackground(rl.Black)
			rl.DrawRectangle(0, 0, int32(DIM_I), int32(DIM_I), bgCol)

			if start_btn_hover {
				rl.DrawRectangleRec(startBtn, rl.DarkGray)
				DrawTextCenter("Start", DIM_I/2, DIM_I/2-int(menuFontSize)-10, int(menuFontSize), rl.White)
			} else {
				rl.DrawRectangleRec(startBtn, rl.RayWhite)
				DrawTextCenter("Start", DIM_I/2, DIM_I/2-int(menuFontSize)-10, int(menuFontSize), rl.Black)
			}
			if exit_btn_hover {
				rl.DrawRectangleRec(exitBtn, rl.DarkGray)
				DrawTextCenter("Exit", DIM_I/2, DIM_I/2+int(menuFontSize)+10, int(menuFontSize), rl.White)
			} else {
				rl.DrawRectangleRec(exitBtn, rl.RayWhite)
				DrawTextCenter("Exit", DIM_I/2, DIM_I/2+int(menuFontSize)+10, int(menuFontSize), rl.Black)
			}
		}
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
		// DRAW END
	}
	fmt.Println("EXITED SUCESSFULLY")
}
