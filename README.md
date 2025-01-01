## AgentSmithU

A forge for making agents - become an agentsmith!

AgentSmith is a command-line tool, chat assistant (command-line and GUI), Python module, and IDE for using and creating AI Agents. It features everything you need to craft AI Agents and the ability to call these agents to complete tasks. Add AI functionality to any program with less than four lines of code. 

Users can freely edit/save/load prompts, chats, and functions to experiment with responses. Agents have the ability to automatically request functions and can interact with literally any other program*. 

*You may have to build an interface but there's tools/examples to cover that

This is a fork of [AgentSmith](https://github.com/mattyleecifer/AgentSmith/) = I updated it so it can work with OpenAI, Mistral, Anthropic, and Ollama (plus anything that uses the OpenAI messages API format). It has a very experimental function that uses a local LLM to convert any unfamiliar API responses to OpenAI's messages API format so that it can read it - this is still a little janky, but it kind of works.

I had to remove "Functions" functionality as that seems to be a more OpenAI specific thing, but I have ideas on how to bring it back. 

For now, basic chat functionality all works and you can build basic agents, just without function calling like [AgentSmith](https://github.com/mattyleecifer/AgentSmith/).

# 2025 note

I have updated the package to be more organized/usable (lol). You will now be able to import agentsmithu/agent to create agents in golang (without using core.go like before). I'll probably be doing updates like this as I learn more about how to do things properly. 


### Features

- **Chat** - You can chat with it like any OpenAI chatbot
- **CLI and GUI** - You can interact with AgentSmith via the CLI or GUI. The CLI is mostly just for chat, but there are a few handy functions in there - type 'help' to see an overview
- **Edit/delete/save/load chats** - Allows you to easily modify chats (even change the AI's response) and store/retrieve them for later use
- **Cost Estimator** - The GUI shows estimated call costs
- **Prompt editor** - An interface for easily editing/saving/loading/deleting prompts. Makes it really easy to prompt engineer
- **Fully Customizable Agents** - You can control which models the main chat uses, max tokencount, autofunction, and more in 'Settings'

The default directory for the AgentSmith is `~/AgentSmith`, but you can easily set this using the `-home` flag.

### Efficient

Being able to remove/edit responses means you can remove redundant information to keep token counts low while retaining important information. Use the call cost estimator to keep costs down.

### How to build agents

# This will need a rewrite with the new version

Look at `/examples` for an examples on how to build simple agents in Golang and Python. The plugins are also good examples for how to build interfaces for the agent to call to access other programs.

`core.go` contains everything you need to build an agent in Go. The `AgentSmith.py` contains everything you need to build an agent in Python. It takes less than four lines of code.

The agents are able to make/receive/process data in whichever way they're programmed - you can even chain agents within agents and get them to talk to each other.

This allows anyone to easily create complex AI apps with multiple agents all with different prompts/functions that can work together to do anything.

### How to run

To run as just a command-line chat, run `agentsmith --console`

To start the GUI, you just have to run: `agentsmith --gui`

(Or `agentsmith.exe --gui` on Windows, etc.)

This will start a server at http://127.0.0.1:49327 - the server is secured so only localhost can connect to it. To allow external connections, launch the app with `-ip <ipaddress>` or `-allowallips`. Use `-port` to specify port.

The default folder is `~/AgentSmith` but this can be set with the `-home` flag

You can create an agent using flags.

Flags:
- `-key` api key" (this must be first)
- `-home` set the home directory for agent
- `-save` save the chat + response to `homedir/Saves/filename.json`
- `-load` load chat from `homedir/Saves` eg `-load example.json` will load `homedir/Saves/example.json`
- `-prompt` set model prompt - otherwise there is a default assistant prompt
- `-model` model name - default is gpt-3.5-turbo
- `-maxtokens` default max tokens is 2048
- `-message` add message from user to chat
- `-messageassistant` add message from assistant to chat

This can be used to build a full agent. The Python module basically follows the same idea - you set the flags/messages and then make a call.

If you want to build from source, you might want to change the PGP keys in `core.go` (this protects your API key*) or set your API key manually in the code. Then it's just `go build` and run.

*The app stores an encrypted API key in `homedir` by default. It will not do this if you specify a key with the `-key` flag.

### How to contribute

I'm just a hobby developer trying to hone my skills. If you want to help, out feel free to open an issue, make a fork, or [email me](mailto:mattyleedev@gmail.com).
