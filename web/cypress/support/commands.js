// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

var namespace = 'octant-cypress-';
var startingContext = '';

before(() => {
  cy.exec('kubectl config current-context').then(result => {
    startingContext = result.stdout;
  });

  var chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
  for (var i = 0; i < 6; i++)
    namespace += chars.charAt(Math.floor(Math.random() * chars.length));

  cy.log('Creating namespace ' + namespace);
  cy.exec('kubectl create namespace ' + namespace);

  cy.exec(
    'kubectl config set-context $CYPRESS_CONTEXT --namespace ' +
      namespace +
      ' --cluster $(kubectl config get-contexts ' +
      startingContext +
      " | tail -1 | awk '{print $3}') \
    --user $(kubectl config get-contexts " +
      startingContext +
      " | tail -1 | awk '{print $4}')",
    { env: { CYPRESS_CONTEXT: Cypress.env('CYPRESS_CONTEXT') } }
  );
  // Octant is expected to start before cypress in ci
  // Setting context here allows spec to get a namespace
  cy.exec('kubectl config use-context $CYPRESS_CONTEXT', {
    env: { CYPRESS_CONTEXT: Cypress.env('CYPRESS_CONTEXT') },
  });
});

after(() => {
  cy.exec('kubectl config use-context ' + startingContext);
  cy.exec('kubectl delete namespace ' + namespace);
  cy.exec('kubectl config delete-context $CYPRESS_CONTEXT', {
    env: { CYPRESS_CONTEXT: Cypress.env('CYPRESS_CONTEXT') },
  });
});
