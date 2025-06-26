# Hermes

Hermes is a static site generator (SSG) written in Go.

It aims to provide a clean, fast, and flexible way to build static websites. While still early in development, Hermes is designed for simplicity and will evolve with new features as the project grows.

![New Content Screenshot](docs/img/new-content.png)

## Notes

Hermes implements an authentication and authorization system with support for multiple users and teams. However, the initial implementation is designed for single-user operation on a personal machine (localhost). Multi-user support will be optimized in future iterations.

All authentication features should work, but since the SSG feature shares some logic with authorization and is being iteratively improved to create a simple API, some changes may impact authentication-related functionality. Please keep this in mind as the system evolves.
