import dagger
from dagger import dag, function, object_type


@object_type
class MultiagentDemo:
    @function
    async def demo(
        self,
        chat_model: str = "gpt-4o",
        coder_model: str = "gpt-o1",
    ) -> str:
        """Demo of a multi-agent system"""

        ws = dag.toy_workspace()

        weather_env = (
            dag.env()
            .with_toy_workspace_input("workspace", ws, "tools to build code")
            .with_string_input("assignment",
                """
                write a program called weather
                that retrieves current weather in San Francisco from wttr.in
                and prints a short report about the temperature and precipitation
                to the console
                """, "the task to be completed")
            .with_toy_workspace_output("workspace", "completed assignment")
        )

        # write a program that gets the current weather report
        coder = (
           dag.llm(model = coder_model)
            .with_env(weather_env)
            .with_prompt_file(dag.current_module().source().file("coder.txt"))
        )

        # save the report to a file
        current_weather = (
            coder
            .env()
            .output("workspace")
            .as_toy_workspace()
            .container()
            .with_exec(["go", "run", "."], redirect_stdout="weather.txt")
            .file("weather.txt")
        )

        summary_env = (
            dag.env()
            .with_file_input("weather", current_weather, "the weather report")
        )
        # give the report to another llm
        summarizer = (
            dag.llm(model = chat_model)
            .with_env(summary_env)
            .with_prompt("""
            The file $weather describes the current weather conditions in San Francisco,
            Don't tell me about the structure or content of the file,
            Briefly, using the weather information provided in the file, tell me if I need to wear a jacket today.
            """)
        )

        return await summarizer.last_reply()
