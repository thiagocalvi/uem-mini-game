# UEM Mini Game

A 3D road runner game developed as an academic project for the Computer Science program at UEM (State University of Maringá). This game was created for the "Hardware and Software Interface Programming" course in collaboration with the "Programming Languages" course, demonstrating practical application of systems programming concepts and modern software architecture patterns.

## 🎮 Game Overview

Navigate a square character down a perspective road while dodging randomly moving obstacles that approach with a realistic 3D depth effect. The game features physics-based movement, collision detection, and an interactive menu system.

## ✨ Features

- **3D Depth Effect**: Obstacles spawn in the distance and grow as they approach the player
- **Physics System**: Realistic gravity, jumping mechanics, and collision detection
- **Dynamic Obstacles**: Random movement patterns within road boundaries
- **Interactive Menus**: Start screen and pause system with navigation
- **Real-time Scoring**: Score tracking that increases over time
- **Smooth Controls**: Responsive player movement and jumping

## 🎯 How to Play

1. **Start**: Press X to begin the game
2. **Move**: Use left/right arrow keys to move your character
3. **Jump**: Press up arrow to jump over obstacles
4. **Avoid**: Dodge the approaching obstacles to survive
5. **Pause**: Press Z to pause and access the menu
6. **Score**: Survive as long as possible to increase your score

## 🎮 Controls

- **Arrow Keys**: Move left/right, jump up
- **X Button**: Start game / Menu selection
- **Z Button**: Pause game
- **Arrow Keys (Menu)**: Navigate menu options

## 🛠️ Technical Implementation

### Architecture
- **Language**: Go with TinyGo compiler
- **Architecture**: Entity Component System (ECS)
- **Platform**: WASM-4 Fantasy Console
- **Rendering**: Custom 3D perspective system

### ECS Components
- **Position**: Entity coordinates (X, Y)
- **Drawable**: Visual properties (size, color)
- **Physics**: Movement and gravity system
- **Depth**: 3D effect and perspective calculations

### Systems
- **Input System**: Handles player controls
- **Physics System**: Applies gravity and movement
- **Collision System**: Detects player-obstacle interactions
- **Render System**: Draws all game elements with proper depth sorting
- **Obstacle System**: Manages obstacle spawning and movement

## 🚀 Building and Running

### Prerequisites
- [TinyGo](https://tinygo.org/getting-started/install/)
- [WASM-4](https://wasm4.org/docs/getting-started/setup)

### Build Commands
```bash
# Build the cartridge
make

# Run with WASM-4
w4 run build/cart.wasm

# Or run in browser
w4 run-native build/cart.wasm
```

### Development
```bash
# Clean build files
make clean

# Force rebuild
make all
```

## 🎓 Educational Context

This project was developed as part of the Computer Science curriculum at UEM, specifically for:

- **Hardware and Software Interface Programming**: Demonstrating low-level programming concepts, memory management, and hardware constraints
- **Programming Languages**: Showcasing Go language features, compilation to WebAssembly, and cross-platform development

The game serves as a practical example of:
- Systems programming in resource-constrained environments
- Modern software architecture patterns (ECS)
- Real-time graphics and game development
- Cross-compilation and WebAssembly deployment

## 📁 Project Structure

```
uem-mini-game/
├── main.go           # Main game logic and ECS implementation
├── w4/
│   └── wasm4.go      # WASM-4 bindings and API
├── Makefile          # Build configuration
├── README.md         # Project documentation
└── build/
    └── cart.wasm     # Compiled game cartridge
```

## 🔗 References

- [WASM-4 Documentation](https://wasm4.org/docs/)
- [TinyGo Documentation](https://tinygo.org/docs/)
- [Entity Component System Pattern](https://en.wikipedia.org/wiki/Entity_component_system)
- [WebAssembly Specification](https://webassembly.org/)