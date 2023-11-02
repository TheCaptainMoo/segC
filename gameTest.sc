playerY = 0
playerX = 0
playerWidth = 40
playerHeight = 20
playerSY = 0

playerSpeed = 2

obstacleX = 0
obstacleY = 0
obstacleWidth = 350
obstacleHeight = 350

score = 0

obstacleXS = 3

width = 0
height = 0

floorCounter = 0

fn isColliding(x0 y0 w0 h0 x1 y1 w1 h1)
{
    x = x1 + w1
    c1 = x0 < x

    x = x0 + w0
    c2 = x > x1

    x = y1 + h1
    c3 = y0 < x

    x = y0 + h0
    c4 = x > y1
    
    ret c1 & c2 & c3 & c4
}

fn start(screenWidth screenHeight)
{
    playerX = screenWidth / 8
    playerX = playerX - playerWidth / 2
    playerY = screenHeight/2

    obstacleX = screenWidth+obstacleWidth
    obstacleY = screenHeight / 2
    obstacleY = obstacleY - obstacleHeight / 2
    
    width = screenWidth
    height = screenHeight

    score = 0

    ret 0
}

fn updatePlayer(jumpKeyPressed)
{
    ; Check if it is outside bounds
    x = height-playerHeight
    if x < playerY
    {
        playerY = x
        floorCounter = floorCounter+1
    }
    elsif playerY < 0
    {
        playerY = 0
        floorCounter = floorCounter+1
    }
    else
    {
        floorCounter = 0
    }

    ; Update Input
    x = 1 - jumpKeyPressed
    x = x * 2
    x = x - 1
    playerSY = x * playerSpeed

    playerY = playerY + playerSY

    ; Add to the score
    if playerX == obstacleX
    {
        score = score + 1
    }

    if floorCounter > 200
    {
        start(width height)
    }

    ret 0
}

fn update(jumpKeyPressed)
{
    updatePlayer(jumpKeyPressed)

    obstacleX = obstacleX-obstacleXS

    if obstacleX < -obstacleWidth
    {
        obstacleX = width + obstacleWidth
    }

    if isColliding(playerX playerY playerWidth playerHeight obstacleX obstacleY obstacleWidth obstacleHeight)
    {
        start(width height)
    }
    
    ret 0
}