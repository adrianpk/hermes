# Hermes

Hermes is a static site generator (SSG) written in Go.

It aims to provide a clean, fast, and flexible way to build static websites. While still early in development, Hermes is designed for simplicity and will evolve with new features as the project grows.

![New Content Screenshot](docs/img/new-content.png)

## Architectural Updates

- The project has moved away from DAOs and transformation layers; now exported structs and pointers are used throughout the codebase.
- The web layer (BFF) no longer calls services directly, but interacts only with the internal API handlers. All business logic and validation is centralized in the API handlers.
- The use of the BFF naming and pattern for this logic is only temporary, to avoid interfering with the current webhandler approach, which remains active for now. This name is not intended for a final release.
- For shared behaviors, we are moving away from struct embedding and instead prefer interfaces implemented via delegation. This approach is a bit more verbose, but it is also more explicit and easier to follow for anyone reading the code.
- The decoupling of the web layer from the service layer also makes a future move to SPA/PWA possible, but this is not a current goal, just a technical side effect of the new design.

## Notes

Although Hermes provides a web interface, it is fundamentally a static site generator (SSG). All content is managed locally (or in the cloud if the service runs on a remote server), and the output is pure HTML, suitable for hosting on platforms like GitHub Pages. Hermes will handle publishing both manually and automatically. Additionally, you can continue to write your markdown content using a text editor if that workflow is more natural or practical for you, rather than using the web interface.

Hermes implements an authentication and authorization system with support for multiple users and teams. However, the initial implementation is designed for single-user operation on a personal machine (localhost). Multi-user support will be optimized in future iterations.

All authentication features should work, but since the SSG feature shares some logic with authorization and is being iteratively improved to create a simple API, some changes may impact authentication-related functionality. Please keep this in mind as the system evolves.
