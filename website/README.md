# SOMAS2020 Visualisation Website

https://somas2020.github.io/SOMAS2020

## Getting started

### Requirements

- [Node.js 14.x.x](https://nodejs.org/en/)
  - `node -v` should produce `v14.x.x`
- [Yarn 1.22.x](https://yarnpkg.com/getting-started/install)
  - `yarn --version` should produce `1.22.x`

### Make sure you have yarn installed on your local machine

Install yarn if you do not already have it installed. Yarn manages the dependencies for the website. The easiest way to do this is

```bash
npm i -g yarn
# On a mac or if you get permission errors
sudo npm i -g yarn
```

### Install dependencies

`yarn install`

## Scripts

### `yarn start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

<!-- ### `yarn test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information. -->

### `yarn build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.

<!-- ### `yarn deploy`

Deploy the app into [GitHub Pages](https://somas2020.github.com/SOMAS2020).\
This should be run automatically by CI. -->

## Information

- This website uses [React](https://reactjs.org/) and [TypeScript](https://www.typescriptlang.org/).
- The library used for UI/UX is [React Bootstrap](https://react-bootstrap.github.io/).
## Deployment

Deployment is done automatically via GitHub Actions whenever a push occurs in the `main` branch (which includes a PR merged into `main`).

You can enable GitHub pages in your own fork to have your own fork's website. Just go to "Settings" on your fork, set the GitHub Pages Source to the `gh-pages` branch and save it. The next merge into your `main` branch should deploy a new version of the page.
