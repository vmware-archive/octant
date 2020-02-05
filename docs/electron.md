# Octant as an Electron application

Octant is migrating to an Electron application. This document provides
information on design and how to run the application in development
mode.

## Running

Octant is still primarily a browser-based application. To run the
Electron application, the Angular frontend and Go backend must be
running first.

Three terminal sessions are required.

* Session 1: `$ npm run start`
* Session 2: `$ go run build.go build run-dev`
* Session 3: `$ npm run  electron:serve`
