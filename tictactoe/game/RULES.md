# Tic Tac Toe API - README

This README provides instructions on how to play Tic Tac Toe using the provided HTTP API.

## Game Rules

Tic Tac Toe is a two-player game played on a 3x3 grid.  The players take turns marking spaces on the grid with their symbol, either 'X' or 'O'.

* **Objective:**  The goal is to be the first player to get three of your symbols in a row, either horizontally, vertically, or diagonally.
* **Turns:** Players alternate turns. Player 'X' always goes first.
* **Winning:** A player wins if they have three of their symbols in a row.
* **Draw:** If all the spaces on the grid are filled and no player has three in a row, the game is a draw.

## API Endpoints

The API provides the following endpoints:

*   `POST /games`: Creates a new game.
*   `GET /games/:gameId`: Gets the current state of a game.
*   `POST /games/:gameId/moves`: Makes a move in a game.

## How to Play

1.  **Create a New Game:**

    *   Send a `POST` request to `/games`.
    *   **Request:**
        ```
        POST /games
        ```
    *   **Response (Success - 201 Created):**

        ```json
        {
            "gameId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
            "gameState": {
                "board": [null, null, null, null, null, null, null, null, null],
                "currentPlayer": "X",
                "status": "ongoing",
                "winner": null
            }
        }
        ```
        *   The `gameId` is a unique identifier for the game.  **Save this ID!** You'll need it for subsequent requests.
        *   The `gameState` provides the initial state of the game:
            *   `board`: An array representing the 3x3 board. `null` indicates an empty space. Indices 0-8 represent the board as follows:

                ```
                0 | 1 | 2
                ---------
                3 | 4 | 5
                ---------
                6 | 7 | 8
                ```
            *   `currentPlayer`: The player whose turn it is ('X' or 'O').
            *   `status`: The current status of the game ('ongoing', 'won', or 'draw').
            *   `winner`:  `null` if the game is ongoing, or the winning player ('X' or 'O') if the game is won.
    *   **Response (Error):** N/A - This endpoint should always succeed unless the server encounters an unexpected error.

2.  **Make a Move:**

    *   Send a `POST` request to `/games/:gameId/moves`, replacing `:gameId` with the actual `gameId` you received when creating the game.
    *   The request body should be a JSON object containing the `position` where you want to place your symbol.  The `position` must be an integer between 0 and 8, representing the board index as described above.
    *   **Request:**

        ```
        POST /games/a1b2c3d4-e5f6-7890-1234-567890abcdef/moves
        Content-Type: application/json

        {
            "position": 4
        }
        ```

        *   In this example, player 'X' is placing their symbol in the center of the board (position 4).
    *   **Response (Success - 200 OK):**

        ```json
        {
            "board": [null, null, null, null, "X", null, null, null, null],
            "currentPlayer": "O",
            "status": "ongoing",
            "winner": null
        }
        ```

        *   The response shows the updated `gameState`. The board now has "X" in the center. The `currentPlayer` has switched to "O", indicating it's now O's turn.
    *   **Response (Error - 400 Bad Request):**

        ```json
        {
            "error": "Invalid position. Must be a number between 0 and 8."
        }
        ```

        *   Returned if the `position` is not a valid integer between 0 and 8.
        ```json
        {
            "error": "Position already taken."
        }
        ```

        *   Returned if the selected `position` is already occupied.

        ```json
        {
            "error": "Game is already over."
        }
        ```

        *   Returned if the game has already been won or drawn.
    *   **Response (Error - 404 Not Found):**
        ```json
        {
            "error": "Game not found"
        }
        ```
        *   Returned if the `gameId` does not exist.

3.  **Get Game State:**

    *   Send a `GET` request to `/games/:gameId`, replacing `:gameId` with the `gameId`.
    *   **Request:**

        ```
        GET /games/a1b2c3d4-e5f6-7890-1234-567890abcdef
        ```
    *   **Response (Success - 200 OK):**

        ```json
        {
            "board": [null, null, null, null, "X", null, null, null, null],
            "currentPlayer": "O",
            "status": "ongoing",
            "winner": null
        }
        ```

        *   The response provides the current `gameState`.
    *   **Response (Error - 404 Not Found):**

        ```json
        {
            "error": "Game not found"
        }
        ```

        *   Returned if the `gameId` does not exist.

4.  **Continue Playing:**

    *   Players alternate making moves by sending `POST` requests to `/games/:gameId/moves` with their desired `position`.
    *   After each move, use `GET /games/:gameId` to check the updated game state and see if there's a winner or a draw.

5.  **Game Over:**

    *   The game ends when either:
        *   The `status` is "won".  The `winner` field will indicate the winning player.
        *   The `status` is "draw".

## Example Game Flow

1.  Player 1 creates a new game and gets `gameId`: `abc-123`
2.  Player 1 makes a move at position 0: `POST /games/abc-123/moves` with body `{"position": 0}`
3.  Player 2 gets the game state: `GET /games/abc-123` and sees the board and that it is their turn.
4.  Player 2 makes a move at position 4: `POST /games/abc-123/moves` with body `{"position": 4}`
5.  Players continue alternating moves until a player wins or the game is a draw.
