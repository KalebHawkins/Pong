```mermaid
---
title: Game State Document
---
stateDiagram-v2
    state "Main Menu" as main
    state "Game Loop" as gameLoop
    state "Paused" as pause
    state "Game Over" as gameOver
    state numPlayers <<choice>>
    state restartQuit <<choice>>

    [*] --> main

    main --> numPlayers
    numPlayers --> OnePlayer
    numPlayers --> TwoPlayer

    OnePlayer --> gameLoop
    TwoPlayer --> gameLoop

    gameLoop --> pause
    pause --> gameLoop
    gameLoop --> gameOver
    gameOver --> restartQuit

    restartQuit --> Quit
    restartQuit --> Restart

    Restart --> main
    Quit --> [*]