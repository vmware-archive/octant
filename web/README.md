<p align="center">
  <img src="./header.png" alt="heptio ui-starter" align="center" width="400" />
</p>

<br>

<p align="center">
  :school_satchel: A heptio boilerplate w/ common tooling & patterns for building UI
</p>

## Why

As Heptio grows, we wanted to build a light UI starter kit that developers can get start being productive with and also encapsulate some of the styling, practices & tools we've gathered.
## Features

Some of our current stack:

- _React.js_ - :tada: v16
- _Styles_ - :necktie: SCSS Styling w/ Stylelint 
- _Lint_ - :police_car: Airbnb's ESLint w/ StandardJS overrides 
- _Babel_ - :nut_and_bolt: along w/ Webpack & PostCSS

## Quick start

1. Clone this repo using `git clone git@github.com:heptio/ui-starter.git`
2. Move to the appropriate directory: `cd ui-starter`.<br />
3. Run `npm install` to install dependencies.<br />
4. Run `npm start` to see the example app at `http://localhost:3000`.

## Commands

### `npm run dev`
Same as `npm run start`. Starts a server at `localhost:3000` by default

### `npm run build`
Builds production mode of the single page app into the `/build` directory

### :fire_engine: `npm run fix`
This runs both `eslint --fix`  & `stylelint --fix` over the appropriate files so that you don't have to worry about formatting or your css being valid.

<br>

[Feel free to submit an issue if you think something should change!](https://github.com/heptio/ui-starter/issues)
