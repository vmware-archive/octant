## Build with Drone

Drone is allows running tests and builds inside containers. This document shows how to run Drone locally using Docker.

1. Install the CLI for your operating system per instructions in [https://docs.drone.io/cli/install/](https://docs.drone.io/cli/install/)

2. Make sure to have an up-to-date version of Docker

3. Run `drone exec` from the base of the octant repository to start frontend and backend tests in parallel. To start only a single step of the pipeline, add `--include` with the step name (e.g. `drone exec --include backend`) 
