/**
 * A generated module for Tictactoe functions
 */
import { dag, Container, object, func, Service } from "@dagger.io/dagger";

@object()
export class Tictactoe {
  constructor() {}

  /**
   * Run the Tic Tac Toe game server
   */
  @func()
  server(): Service {
    const port = 3000;
    const game = dag
      .currentModule()
      .source()
      .directory("game")
      .withoutDirectory("node_modules");

    return dag
      .container()
      .from("node")
      .withWorkdir("/app")
      .withDirectory("/app", game)
      .withMountedCache("/root/.npm", dag.cacheVolume("tictactoe_node_modules"))
      .withExec(["npm", "install"])
      .withExec(["/usr/local/bin/npx", "tsc"])
      .withEnvVariable("PORT", `${port}`)
      .withExposedPort(port)
      .asService({ args: ["node", "dist/index.js"] });
  }
  /**
   * Play the Tic Tac Toe game
   */
  @func()
  play(): Service {
    return dag
      .container()
      .from("nginx")
      .withFile(
        "/usr/share/nginx/html/index.html",
        dag.currentModule().source().file("game/index.html"),
      )
      .withExposedPort(80)
      .asService({ useEntrypoint: true });
  }

  /**
   * Have the LLM play the Tic Tac Toe game
   * FIXME: This needs the IP instead of the Service itself until a Dagger bug is fixed.
   */
  @func()
  async ai(gameEndpoint: string, gameId: string): Promise<string[]> {
    // get an instance of the game client to give to the llm
    const client = dag.tictactoeClient(gameEndpoint, gameId);
    // the game client gives the LLM the tools to play
    // it will make moves and wait for the opponent to make moves
    return dag
      .llm()
      .withTictactoeClient(client)
      .withPromptFile(dag.currentModule().source().file("prompt.md"))
      .history();
  }
}
