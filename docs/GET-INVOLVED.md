# Get Involved

## 🎯 Objectives (What we aim to do and what we don’t)

### Objectives

1. Make interactions with LLMs easier and observable.
2. Keep it lightweight and focused on the core.
3. Make it accessible to everyone.
4. Build with security and privacy in mind — prioritize local LLMs.

### Non-Objectives

1. Bloat the library/package with external dependencies or excessive features.
2. Mimic large libraries that enable LLM interaction.

---

## 🏛️ System Design

### Core

The Core groups the following modules: `Input`, `LLM`, `Output`, and `Pipe`. The image below illustrates this concept.

![High Level Design](/docs/concept/core.svg)  
**Caption**: The image shows the three core components — `Input`, `LLM`, and `Output`. A session is defined as the flow from `Input` to `Output`, enabling fine-grained observability.

---

### Agent

There is currently an experimental agent implementation that follows the ReAct pattern to call tools. The current focus remains on the Core.

In short, the Agent consists of three components: `Tools -> Agent -> Runner`. The Runner runs multiple cycles until the LLM returns a final result.

---

## 🚴🏽‍♂️ Roadmap

The focus is first **observability**, then **security**, and finally **compliance** within the Gogantic Core.

### Phase 1

**Objective**: Implement observability into Gogantic Core

- [x] Build core components to interact with large language models
- [ ] Conduct research and developer/company interviews to validate the need for lightweight LLM interaction (focused on observability, security, and compliance)
- [ ] Publish observability flow for core components
- [ ] Release Gogantic Core with observability support for beta testing

### Phase 2

_To be determined_

### Phase 3

_To be determined_

### Side Projects

We are also working on an agent capable of autonomously executing tasks. Here's the current status:

- ✅ Create an agent that can interact with external tools
- 🔜 Develop a **Director Agent** to manage complex tasks by coordinating multiple agents

---

## 🧑🏽‍💻 Developer Guide

1. Select an issue from the issues tab
2. Fork the repository
3. Checkout a new branch: `git checkout -b <descriptive-name>`
4. Make changes based on the issue (feel free to ask if anything is unclear)
5. Add, commit, and push your changes using [Conventional Commits](https://www.conventionalcommits.org)
6. Create a pull request — see: [PR Example](/docs/PULL-REQUEST-EXAMPLE.md). A maintainer will review, test, and merge it

---

## 🛠️ Contributing (Cheat Sheet)

_This section will provide a quick reference for contributing. Coming soon._

---

## 🧭 Code of Conduct

_This section will define expected behavior and guidelines for contributors. Coming soon._

---

## 🔧 Maintainers

- Tobias Gleiter

---

## 🚀 Release Process

_This section will explain how releases are versioned and deployed. Coming soon._

---

## 📚 Documentation

_This section will link or reference further documentation. Coming soon._
