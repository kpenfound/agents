import express, { Request, Response, Application, NextFunction } from "express";
import cors from "cors";
import { v4 as uuidv4 } from "uuid";
import { ParamsDictionary } from "express-serve-static-core";
import { ParsedQs } from "qs";

// Define types for the game state
type Player = "X" | "O";
type Board = (Player | null)[];
type GameStatus = "ongoing" | "won" | "draw";

interface GameState {
  board: Board;
  currentPlayer: Player;
  status: GameStatus;
  winner: Player | null;
}

// Game state manager (in-memory storage)
const games: { [gameId: string]: GameState } = {};

// Function to create a new game
function createNewGame(): { gameId: string; initialState: GameState } {
  const gameId = uuidv4();
  const initialState: GameState = {
    board: Array(9).fill(null),
    currentPlayer: "X",
    status: "ongoing",
    winner: null,
  };
  games[gameId] = initialState;
  return { gameId, initialState };
}

// Function to check for a winner
function checkWinner(board: Board): Player | null {
  const winningCombinations = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8], // Rows
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8], // Columns
    [0, 4, 8],
    [2, 4, 6], // Diagonals
  ];

  for (const combination of winningCombinations) {
    const [a, b, c] = combination;
    if (board[a] && board[a] === board[b] && board[a] === board[c]) {
      return board[a];
    }
  }
  return null;
}

// Function to check for a draw
function checkDraw(board: Board): boolean {
  return board.every((cell) => cell !== null);
}

// Function to update the game state after a move
function updateGameState(gameId: string, position: number): GameState {
  const game = games[gameId];
  if (!game) {
    throw new Error("Game not found.");
  }

  if (game.status !== "ongoing") {
    throw new Error("Game is already over.");
  }

  if (game.board[position] !== null) {
    throw new Error("Position already taken.");
  }

  const newBoard = [...game.board];
  newBoard[position] = game.currentPlayer;

  const winner = checkWinner(newBoard);
  let newStatus: GameStatus = "ongoing";
  let newWinner: Player | null = null;

  if (winner) {
    newStatus = "won";
    newWinner = winner;
  } else if (checkDraw(newBoard)) {
    newStatus = "draw";
  }

  const newGameState: GameState = {
    board: newBoard,
    currentPlayer: game.currentPlayer === "X" ? "O" : "X",
    status: newStatus,
    winner: newWinner,
  };

  games[gameId] = newGameState;
  return newGameState;
}

// Express app setup
const app: Application = express(); // Explicitly type 'app'
const port = process.env.PORT || 3000;

app.use(cors()); // Enable CORS for all origins
app.use(express.json()); // for parsing application/json in req.body

// API endpoints

// 1. Create a new game
const createGameHandler = (
  req: Request,
  res: Response,
  next: NextFunction,
): void => {
  const { gameId, initialState } = createNewGame();
  res.status(201).json({ gameId, gameState: initialState });
  next(); // Important: Call next() to complete the middleware cycle
};
app.post("/games", createGameHandler);

// 2. Get game state by ID
const getGameHandler = (
  req: Request,
  res: Response,
  next: NextFunction,
): void => {
  const gameId = req.params.gameId;
  const game = games[gameId];

  if (game) {
    res.json(game);
  } else {
    res.status(404).json({ error: "Game not found" });
  }
  next(); // Important: Call next() to complete the middleware cycle
};
app.get("/games/:gameId", getGameHandler);

// 3. Make a move
const makeMoveHandler = (
  req: Request,
  res: Response,
  next: NextFunction,
): void => {
  const gameId = req.params.gameId;
  const { position } = req.body; // Expects { position: number }

  if (typeof position !== "number" || position < 0 || position > 8) {
    res
      .status(400)
      .json({ error: "Invalid position. Must be a number between 0 and 8." });
    return; // Terminate the function after sending the error
  }

  try {
    const updatedGameState = updateGameState(gameId, position);
    res.json(updatedGameState);
    next(); // Important: Call next() to complete the middleware cycle
  } catch (error: any) {
    // Explicitly type error as any to access .message
    if (error.message === "Game not found.") {
      res.status(404).json({ error: "Game not found" });
      return; // Terminate the function after sending the error
    } else if (error.message === "Game is already over.") {
      res.status(400).json({ error: "Game is already over" });
      return; // Terminate the function after sending the error
    } else if (error.message === "Position already taken.") {
      res.status(400).json({ error: "Position already taken" });
      return; // Terminate the function after sending the error
    } else {
      console.error("Unexpected error:", error); // Log unexpected errors
      res.status(500).json({ error: "Internal server error" });
      return; // Terminate the function after sending the error
    }
  }
};
app.post("/games/:gameId/moves", makeMoveHandler);

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});
