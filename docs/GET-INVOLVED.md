# Get Involved

## 🎯 Objectives (Objectives and non-objectives)

### Objectives

1. Make interactions with LLMs easier.
2. Keep it lightweight and focus on the core.
3. Make it accessible to all people.
4. Build with security and privacy in mind: focus on local LLMs.

### Non-Objectives

1. Overbloat the libary/package with external dependencies and features.
2. Mimic other big libraries that enable interaction with LLMs

## 🏛️ System Design

There are two main distinctions. Firstly, the Core and secondly the Agents.
In short, Core consists of these three components: `Input -> LLM -> Output`.

Agents utilises the ReAct pattern and can call tools:
In short, Agent consists of three components: `Tools -> Agent -> Executor`
The executor will run n times till the the LLM returns a final result.

## 🚴🏽‍♂️ Roadmap

- ✅ Build core components to interact with large language models.
- ✅ Create an Agent that can interact with the outside world using tools.
- 🔜 Develop a Director Agent that manages complex tasks by coordinating multiple Agents.

## 🧑🏽‍💻 Developer Guide

1. Select an issue from the issues tab.
2. Fork the repository
3. Checkout to a new branch `git checkout -b <descriptive-name>`
4. Do the changes according to the issue (if not clear feel free to ask)
5. Add, commit and push the changes (please use [conventional commits](https://www.conventionalcommits.org))
6. Create a pull request. It will be tested and merged by the maintainers

## Contributing (i.e. Cheatsheet)

## Code Of Conduct

## Maintainers

1. Tobias Gleiter

## Release Process

## Documentation
