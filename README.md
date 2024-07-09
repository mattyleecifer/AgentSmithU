# AgentSmithU

This is a fork of [AgentSmith](https://github.com/mattyleecifer/AgentSmith/), a project where I used Golang/HTMX to build a program that can make agents/work as a chat interface. This project expands on that to allow users to use any chat model API that returns responses in a similar format to OpenAI's conversational responses. 

It now works with OpenAI, Mistral, Anthropic, and Ollama (plus anything that uses the OpenAI messages API format).

I had to remove "Functions" functionality as that seems to be a more OpenAI specific thing, but I have ideas on how to bring it back. 

For now, basic chat functionality all works and you can build basic agents, just without function calling like [AgentSmith](https://github.com/mattyleecifer/AgentSmith/).

Ideas/Todo:
- Create a "Function calling" feature that adds an extra prompt that makes the model output in JSON - though, this can just be done with a normal prompt so I don't know how valuable it would be
  
