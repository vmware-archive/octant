# Contributing to Octant's UI

This document describes how to setup your development environment to contribute to Octant's UI.

**Note**: Most the code related to this development is organized under the `web/` directory.

## Getting started

### Dependencies

Our web UI is built on Node.js 10+ and npm 6+. It was generated with [Angular CLI](https://github.com/angular/angular-cli) version 7.3.3 so you'll want to install & get familiar with that tool to understand how we some of our npm scripts work. Here are some of the major libraries we use:

- [Angular v7.2+](http://angular.io)
- [TypeScript v3+](https://www.typescriptlang.org/)
- [Clarity v1+](https://clarity.design/)
- [Lodash v4+](https://lodash.com/)

There are different ways to installing these dependencies:

The most thorough way that verifies that you have the needed dependencies to build/develop both the Octant's UI and go server component is run `make ci` from the root directory. This is the command used to build a final distributable.

To install just the UI dependencies from the root directory, we have the Makefile command `web-deps` for just npm installation.

You can also run `npm install` yourself, but this will only work if you are within the `web/` directory.

### Running development mode

Once you have the necessary dependencies installed, you can start the backend and the frontend server with the following Make command:

    make -j ui-server ui-client

This will run both server processes in parallel. You can also run each command in different terminals (which would help with debugging through logs) with `make ui-server` and `make ui-client`. This should open up your browser pointing to `[http://localhost:4200/](http://localhost:4200/)`.

## Directory structure

The entry file into our application is `src/main.ts` but most of our UI logic is written under `src/app/modules` with our base module being `OverviewModule`.

Here is a summary of the `app/` directory structure:

    ➜ tree ./web/src/app -L 1
    ├── components     // components living outside any modules
    ├── models         // typescript definitions of backend payloads
    ├── modules        // NgModules (https://angular.io/guide/architecture-modules)
    ├── services       // globally available service classes
    ├── testing        // testing mocks and stubs
    └── util           // app-wide reused functions

## Testing & Linting

For testing we use [Karma's Test Runner](https://karma-runner.github.io/latest/index.html) with [Jasmine](https://jasmine.github.io). To lint the codebase we rely on [TSLint](https://palantir.github.io/tslint/) & [Prettier](https://prettier.io/) to keep our codebase formatted.

There are 3 commands that can help keep PRs tested and linted properly:

- `npm test` uses Karma's Chrome launcher to open up an instance of Chrome and run our test suite against that environment. This is primarily the tool our team uses.
- `npm run test:headless` helps run our test suite against a headless version of Chrome
- `npm run lint` runs our static analysis tools against our TypeScript code

## Building production

To build a production version of the UI, you can run `npm run build` yourself which will build the productinized assets to `dist/`.

To build a full production binary (including the backend server), you can run `make ci` from the root directory.
