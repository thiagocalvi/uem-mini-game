package main

import (
	"cart/w4"
	"strconv"
)

// Position armazena as coordenadas X e Y da nossa entidade.
type Position struct {
	X, Y int
}

// Drawable armazena as informações necessárias para desenhar a entidade.
type Drawable struct {
	Size  uint
	Color uint16
}

// PlayerInput marca a entidade como sendo controlável pelo jogador e define sua velocidade.
type PlayerInput struct {
	Speed int
}

// Physics armazena informações de física como velocidade e gravidade.
type Physics struct {
	VelocityX float32 // Velocidade horizontal
	VelocityY float32 // Velocidade vertical
	OnGround  bool    // Se está no chão
	Gravity   float32 // Força da gravidade
	JumpPower float32 // Força do pulo
	Friction  float32 // Atrito para desacelerar movimento horizontal
}

// Depth armazena informações de profundidade para efeito 3D
type Depth struct {
	Z        float32 // Profundidade (0.0 = fundo da tela, 1.0 = frente da tela)
	Speed    float32 // Velocidade de aproximação
	BaseSize uint    // Tamanho base do obstáculo
	MaxSize  uint    // Tamanho máximo do obstáculo
	Active   bool    // Se o obstáculo está ativo
}

// Armazenamento de Componentes para o jogador e obstáculos
var (
	playerPosition PlayerInput
	squarePosition Position
	squareDrawable Drawable
	squarePhysics  Physics

	// Arrays para múltiplos obstáculos
	obstaclePositions [10]Position // Máximo de 10 obstáculos
	obstacleDrawables [10]Drawable
	obstacleDepths    [10]Depth

	// Variáveis para geração de obstáculos
	obstacleSpawnTimer int
	nextObstacleIndex  int

	// Variáveis para pontuação
	score      int
	scoreTimer int

	// Variáveis para controle do menu e estado do jogo
	gameState     int
	menuOption    int
	previousInput uint8
)

const (
	SQUARE_SIZE    = 20
	PLAYER_SPEED   = 2
	GRAVITY        = 0.5
	JUMP_POWER     = 8.0
	FRICTION       = 0.8
	GROUND_LEVEL   = 140 // Nível do chão (tela tem 160 de altura)
	MAX_VELOCITY_Y = 15.0

	// Constantes para obstáculos
	OBSTACLE_BASE_SIZE  = 2     // Tamanho inicial do obstáculo (2x1)
	OBSTACLE_MAX_SIZE   = 20    // Tamanho máximo do obstáculo (limitação)
	OBSTACLE_SPEED      = 0.004 // Velocidade de aproximação (mais lenta para melhor efeito)
	OBSTACLE_SPAWN_RATE = 120   // Frames entre spawn de obstáculos (2 segundos a 60fps)

	// Constantes para pontuação
	SCORE_INCREMENT_RATE = 50 // Frames entre incrementos de pontuação (1 segundo a 60fps)

	// Estados do jogo
	STATE_MENU    = 0
	STATE_PLAYING = 1
	STATE_PAUSED  = 2

	// Opções do menu
	MENU_START   = 0
	MENU_RESUME  = 1
	MENU_RESTART = 2
)

//go:export start
func start() {
	// Inicializa os componentes da nossa entidade "quadrado".

	// Componente PlayerInput
	playerPosition.Speed = PLAYER_SPEED

	// Componente Position: coloca o quadrado no chão, centralizado horizontalmente
	squarePosition.X = (w4.SCREEN_SIZE / 2) - (SQUARE_SIZE / 2)
	squarePosition.Y = GROUND_LEVEL

	// Componente Drawable
	squareDrawable.Size = SQUARE_SIZE
	squareDrawable.Color = 0x43 // Cor primária da paleta

	// Componente Physics
	squarePhysics.VelocityX = 0
	squarePhysics.VelocityY = 0
	squarePhysics.OnGround = true
	squarePhysics.Gravity = GRAVITY
	squarePhysics.JumpPower = JUMP_POWER
	squarePhysics.Friction = FRICTION

	// Inicializa todos os obstáculos como inativos
	for i := 0; i < 10; i++ {
		obstacleDepths[i].Active = false
	}

	// Inicializa pontuação
	score = 0
	scoreTimer = 0

	// Inicializa o estado do jogo no menu
	gameState = STATE_MENU
	menuOption = MENU_START
	previousInput = 0
}

//go:export update
func update() {
	gamepad := *w4.GAMEPAD1

	switch gameState {
	case STATE_MENU:
		menuSystem()
		renderMenu()
	case STATE_PLAYING:
		// Verifica se o botão 2 foi pressionado para pausar
		if gamepad&w4.BUTTON_2 != 0 && previousInput&w4.BUTTON_2 == 0 {
			gameState = STATE_PAUSED
			menuOption = MENU_RESUME
		}

		// Executamos nossos sistemas em ordem a cada frame.
		inputSystem()
		physicsSystem()
		movementSystem()
		collisionSystem() // Sistema de detecção de colisão
		obstacleSystem()  // Sistema para gerenciar obstáculos
		scoreSystem()     // Sistema para gerenciar pontuação
		renderSystem()
	case STATE_PAUSED:
		menuSystem()
		renderPauseMenu()
	}

	previousInput = gamepad
}

// inputSystem lê a entrada do controle e atualiza a velocidade horizontal e pulo.
func inputSystem() {
	// Este sistema opera em entidades que possuem o componente PlayerInput.

	gamepad := *w4.GAMEPAD1 // Pega o estado do controle do jogador 1.

	// Movimento horizontal
	if gamepad&w4.BUTTON_LEFT != 0 {
		squarePhysics.VelocityX = -float32(playerPosition.Speed)
	} else if gamepad&w4.BUTTON_RIGHT != 0 {
		squarePhysics.VelocityX = float32(playerPosition.Speed)
	} else {
		// Aplica atrito quando não há entrada
		squarePhysics.VelocityX *= squarePhysics.Friction
	}

	// Pulo (só pode pular se estiver no chão)
	if gamepad&w4.BUTTON_UP != 0 && squarePhysics.OnGround {
		squarePhysics.VelocityY = -squarePhysics.JumpPower
		squarePhysics.OnGround = false
	}
}

// physicsSystem aplica gravidade e limita velocidades
func physicsSystem() {
	// Aplica gravidade se não estiver no chão
	if !squarePhysics.OnGround {
		squarePhysics.VelocityY += squarePhysics.Gravity

		// Limita a velocidade de queda
		if squarePhysics.VelocityY > MAX_VELOCITY_Y {
			squarePhysics.VelocityY = MAX_VELOCITY_Y
		}
	}
}

// movementSystem atualiza a posição baseada na velocidade e lida com colisões nas laterais e no chão
func movementSystem() {
	// Atualiza posição horizontal
	squarePosition.X += int(squarePhysics.VelocityX)

	// Colisão com as laterais da tela
	if squarePosition.X < 0 {
		squarePosition.X = 1
		squarePhysics.VelocityX = 0
	} else if squarePosition.X > w4.SCREEN_SIZE-int(squareDrawable.Size) {
		squarePosition.X = w4.SCREEN_SIZE - int(squareDrawable.Size) - 1
		squarePhysics.VelocityX = 0
	}

	// Atualiza posição vertical
	squarePosition.Y += int(squarePhysics.VelocityY)

	// Colisão com o chão
	if squarePosition.Y >= GROUND_LEVEL {
		squarePosition.Y = GROUND_LEVEL
		squarePhysics.VelocityY = 0
		squarePhysics.OnGround = true
	}
}

// Desenha as bordas da estrada, criando um efeito de profundidade
func drawBorder() {
	//Lateral direita
	xd, yd, wd := 159, 139, uint(1)
	for i := 0; i <= 79; i++ {
		w4.Rect(xd, yd, wd, 1)
		yd--
		xd--
		wd++
	}

	//Lateral esquerda
	xe, ye, we := 0, 139, uint(1)
	for i := 0; i <= 79; i++ {
		w4.Rect(xe, ye, we, 1)
		ye--
		we++
	}
}

// renderSystem desenha as entidades na tela.
func renderSystem() {
	// Desenha a pontuação no canto superior direito
	drawScore()

	drawBorder()

	// Desenha todos os obstáculos ativos
	for i := 9; i >= 0; i-- {
		if obstacleDepths[i].Active {
			*w4.DRAW_COLORS = obstacleDrawables[i].Color
			w4.Rect(obstaclePositions[i].X, obstaclePositions[i].Y, obstacleDrawables[i].Size, obstacleDrawables[i].Size)
		}
	}

	// Desenha o jogador
	*w4.DRAW_COLORS = squareDrawable.Color
	w4.Rect(squarePosition.X, squarePosition.Y, squareDrawable.Size, squareDrawable.Size)
}

// obstacleSystem gerencia a criação, movimento e renderização dos obstáculos
func obstacleSystem() {
	// Incrementa o timer de spawn
	obstacleSpawnTimer++

	// Cria um novo obstáculo se chegou a hora
	if obstacleSpawnTimer >= OBSTACLE_SPAWN_RATE {
		spawnObstacle()
		obstacleSpawnTimer = 0
	}

	// Atualiza todos os obstáculos ativos
	for i := 0; i < 10; i++ {
		if obstacleDepths[i].Active {
			updateObstacle(i)
		}
	}
}

// spawnObstacle cria um novo obstáculo na posição específica (79, 61)
func spawnObstacle() {
	// Encontra um slot livre para o obstáculo
	for i := 0; i < 10; i++ {
		if !obstacleDepths[i].Active {
			// Posição fixa onde o obstáculo aparece
			obstaclePositions[i].X = 79 // Posição X fixa
			obstaclePositions[i].Y = 61 // Posição Y fixa

			obstacleDrawables[i].Size = 2     // Tamanho inicial 2x1 (será tratado como width)
			obstacleDrawables[i].Color = 0x32 // Cor vermelha

			obstacleDepths[i].Z = 0.0 // Começa no fundo (pequeno)
			obstacleDepths[i].Speed = OBSTACLE_SPEED
			obstacleDepths[i].BaseSize = 2 // Tamanho inicial 2x1
			obstacleDepths[i].MaxSize = 30 // Tamanho máximo (limitação)
			obstacleDepths[i].Active = true

			break // Sai do loop após criar um obstáculo
		}
	}
}

// updateObstacle atualiza um obstáculo específico
func updateObstacle(index int) {
	// Aproxima do jogador
	obstacleDepths[index].Z += obstacleDepths[index].Speed

	// Se chegou muito perto desativa o obstáculo
	if obstaclePositions[index].Y+int(obstacleDrawables[index].Size) == 160 {
		obstacleDepths[index] = Depth{}       // Reseta o obstáculo
		obstaclePositions[index] = Position{} // Reseta a posição
		obstacleDrawables[index] = Drawable{} // Reseta o drawable
		obstacleDepths[index].Active = false  // Marca como inativo
		obstacleSpawnTimer = 0                // Reseta o timer de spawn
		return
	}

	// Calcula o tamanho baseado na profundidade (tanto largura quanto altura)
	sizeRange := float32(obstacleDepths[index].MaxSize - obstacleDepths[index].BaseSize)
	currentSize := obstacleDepths[index].BaseSize + uint(obstacleDepths[index].Z*sizeRange)
	obstacleDrawables[index].Size = currentSize

	// Move o obstáculo em direção ao jogador (efeito de aproximação)
	// Começando em (79, 61) e se expandindo/movendo conforme cresce

	// Calcula o deslocamento baseado no crescimento
	sizeOffset := int(currentSize) / 2

	// Ajusta posição X (centraliza o crescimento)
	obstaclePositions[index].X = 79 - sizeOffset

	// Ajusta posição Y (o obstáculo "desce" conforme cresce)
	obstaclePositions[index].Y = 61 + int(obstacleDepths[index].Z*40) // Move para baixo

	// Garante que não saia da tela
	if obstaclePositions[index].X < 0 {
		obstaclePositions[index].X = 0
	}
	if obstaclePositions[index].X+int(currentSize) > w4.SCREEN_SIZE {
		obstaclePositions[index].X = w4.SCREEN_SIZE - int(currentSize)
	}
	if obstaclePositions[index].Y+int(currentSize) > w4.SCREEN_SIZE {
		obstaclePositions[index].Y = w4.SCREEN_SIZE - int(currentSize)
	}
}

// collisionSystem detecta colisões entre o jogador e os obstáculos
func collisionSystem() {
	// Verifica colisão com todos os obstáculos ativos
	for i := 9; i >= 0; i-- {
		if obstacleDepths[i].Active {
			// Calcula as bordas do jogador
			playerLeft := squarePosition.X
			playerRight := squarePosition.X + int(squareDrawable.Size)
			playerTop := squarePosition.Y
			playerBottom := squarePosition.Y + int(squareDrawable.Size)

			// Calcula as bordas do obstáculo
			obstacleLeft := obstaclePositions[i].X
			obstacleRight := obstaclePositions[i].X + int(obstacleDrawables[i].Size)
			obstacleTop := obstaclePositions[i].Y
			obstacleBottom := obstaclePositions[i].Y + int(obstacleDrawables[i].Size)

			// Verifica se há sobreposição (colisão)
			if playerRight >= obstacleLeft && playerLeft <= obstacleRight &&
				playerBottom >= obstacleTop && playerTop <= obstacleBottom && obstacleBottom >= 158 {
				// Colisão detectada. Reinicia o jogo
				resetGame()
				return
			}
		}
	}
}

// resetGame reinicia o jogo quando há colisão
func resetGame() {
	// Reposiciona o jogador
	squarePosition.X = (w4.SCREEN_SIZE / 2) - (SQUARE_SIZE / 2)
	squarePosition.Y = GROUND_LEVEL

	// Reseta a física
	squarePhysics.VelocityX = 0
	squarePhysics.VelocityY = 0
	squarePhysics.OnGround = true

	// Remove todos os obstáculos
	for i := 0; i < 10; i++ {
		obstacleDepths[i].Active = false
	}

	// Reseta o timer de spawn
	obstacleSpawnTimer = 0

	// Reseta a pontuação
	score = 0
	scoreTimer = 0
}

// renderGame renderiza o jogo em execução
func renderGame() {
	drawBorder()

	// Desenha todos os obstáculos ativos
	for i := 9; i >= 0; i-- {
		if obstacleDepths[i].Active {
			*w4.DRAW_COLORS = obstacleDrawables[i].Color
			w4.Rect(obstaclePositions[i].X, obstaclePositions[i].Y, obstacleDrawables[i].Size, obstacleDrawables[i].Size)
		}
	}

	// Desenha o jogador
	*w4.DRAW_COLORS = squareDrawable.Color
	w4.Rect(squarePosition.X, squarePosition.Y, squareDrawable.Size, squareDrawable.Size)
}

// scoreSystem incrementa a pontuação a cada intervalo de tempo
func scoreSystem() {
	scoreTimer++
	if scoreTimer >= SCORE_INCREMENT_RATE {
		score++
		scoreTimer = 0
	}
}

// drawScore desenha a pontuação no canto superior esquerdo
func drawScore() {
	valueScore := strconv.Itoa(score) // Converte a pontuação para string
	w4.Text("Score:", 0, 0)           // Desenha o texto "Score:" no canto superior esquerdo
	w4.Text(valueScore, 48, 0)
}

// menuSystem controla a navegação e seleção do menu
func menuSystem() {
	gamepad := *w4.GAMEPAD1

	// Navegação no menu
	if gamepad&w4.BUTTON_DOWN != 0 && previousInput&w4.BUTTON_DOWN == 0 {
		if gameState == STATE_MENU {
			menuOption = MENU_START
		} else if gameState == STATE_PAUSED {
			menuOption++
			if menuOption > MENU_RESTART {
				menuOption = MENU_RESUME
			}
		}
	} else if gamepad&w4.BUTTON_UP != 0 && previousInput&w4.BUTTON_UP == 0 {
		if gameState == STATE_PAUSED {
			menuOption--
			if menuOption < MENU_RESUME {
				menuOption = MENU_RESTART
			}
		}
	}

	// Seleção da opção
	if gamepad&w4.BUTTON_1 != 0 && previousInput&w4.BUTTON_1 == 0 {
		switch menuOption {
		case MENU_START:
			// Inicia o jogo
			gameState = STATE_PLAYING
			score = 0
			scoreTimer = 0
		case MENU_RESUME:
			// Retoma o jogo
			gameState = STATE_PLAYING
		case MENU_RESTART:
			// Reinicia o jogo
			resetGame()
			gameState = STATE_PLAYING
			score = 0
			scoreTimer = 0
		}
	}
}

// renderMenu renderiza o menu principal
func renderMenu() {
	*w4.DRAW_COLORS = 0x4
	w4.Text("MINI GAME", 45, 40)

	w4.Text("X para iniciar", 32, 70)

	w4.Text("Controles:", 40, 90)
	w4.Text("Mover \x84 e \x85", 32, 100)
	w4.Text("\x86  para pular", 32, 110)
	w4.Text("Z para pausar", 32, 120)
}

// renderPauseMenu renderiza o menu de pausa
func renderPauseMenu() {
	// Desenha pontuação
	drawScore()

	// Desenha menu de pausa por cima
	*w4.DRAW_COLORS = 0x3
	w4.Text("PAUSADO", 55, 50)
	w4.Text("Use \x86\x87 para navegar", 5, 100)
	w4.Text("X para selecionar", 10, 110)

	// Desenha opções do menu de pausa
	for i := MENU_RESUME; i <= MENU_RESTART; i++ {
		if i == menuOption {
			*w4.DRAW_COLORS = 0x3
			w4.Text(">", 40, 70+(i-1)*10) // Seta indicadora
			*w4.DRAW_COLORS = 0x4         // Cor escura para o texto selecionado
		} else {
			*w4.DRAW_COLORS = 0x3
		}

		switch i {
		case MENU_RESUME:
			w4.Text("Continuar", 50, 70)
		case MENU_RESTART:
			w4.Text("Reiniciar", 50, 80)
		}
	}
}
