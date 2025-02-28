/**
 * A generated module for TictactoeClient functions
 */
import {
  dag,
  Container,
  Directory,
  object,
  func,
  Service,
} from "@dagger.io/dagger";

@object()
export class TictactoeClient {
  game: string;
  gameId: string;

  constructor(game: string, gameId: string) {
    this.game = game;
    this.gameId = gameId;
  }
  /**
   * Make a move and get back the game state and turn information
   * @param position The game board position to play
   */
  @func()
  async move(position: string): Promise<string> {
    let turnInfo = `You played at position ${position}\n`;
    let response = JSON.parse(await this.post(position));
    if (response.error) {
      // LOL THIS IS GOLANG NOW
      return response.error;
    }

    // Check gameover
    if (response.status == "draw") {
      return turnInfo + `The game is over. It ended in a draw.`;
    }

    if (response.status == "won") {
      return turnInfo + `The game is over. Player ${response.winner} won.`;
    }
    await new Promise((resolve) => setTimeout(resolve, 1000));
    let newState = JSON.parse(await this.get());
    while (((newState = JSON.parse(await this.get())), newState.turn == "X")) {
      // wait for turn to be O
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }

    // Check for opponent move
    for (let i = 0; i < newState.board.length; i++) {
      if (newState.board[i] != response.board[i]) {
        turnInfo += `Your opponent played at position ${i}\n`;
        break;
      }
    }

    // Check gameover again
    if (response.status == "draw") {
      return turnInfo + `The game is over. It ended in a draw.`;
    }

    if (response.status == "won") {
      return turnInfo + `The game is over. Player ${response.winner} won.`;
    }

    return (
      turnInfo +
      `It is now your turn.\nThe current board is:\n${this.stringifyBoard(newState.board)}\n`
    );
  }

  /**
   * Read the current state of the game
   */
  @func()
  async read(): Promise<string> {
    let state = JSON.parse(await this.get());
    return `The current board is:\n${this.stringifyBoard(state.board)}\n`;
  }

  async get(): Promise<string> {
    return await dag
      .container()
      .from("alpine/curl")
      // .withServiceBinding("game", this.game)
      .withExec([
        "curl",
        "-X",
        "GET",
        `http://${this.game}:3000/games/${this.gameId}`,
      ])
      .stdout();
  }

  async post(position: string): Promise<string> {
    return await dag
      .container()
      .from("alpine/curl")
      // .withServiceBinding("game", this.game)
      .withExec([
        "curl",
        "-X",
        "POST",
        "-H",
        "Content-Type: application/json",
        "-d",
        `{"position": ${position}}`,
        `http://${this.game}:3000/games/${this.gameId}/moves`,
      ])
      .stdout();
  }

  stringifyBoard(board: string[]): string {
    let boardStr = "";
    for (let i = 0; i < board.length; i++) {
      boardStr += this.stringifyBoardSpace(board[i]);
      if (i % 3 == 2) {
        boardStr += "\n";
      } else {
        boardStr += " ";
      }
    }
    return boardStr;
  }

  stringifyBoardSpace(space: string): string {
    if (space == "O") {
      return "O";
    }
    if (space == "X") {
      return "X";
    }
    return "_";
  }
}
