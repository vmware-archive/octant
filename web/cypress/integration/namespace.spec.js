describe('Namespace', () => {
  beforeEach(() => {
    cy.visit('/', {
      onBeforeLoad: spyOnAddEventListener
    }).then(waitForAppStart);
  });

  it('namespaces navigation', () => {
    cy.exec(
      `kubectl config view --minify --output 'jsonpath={..namespace}'`
    ).then(result => {
      cy.contains(/Namespaces/).click();
      cy.contains(result.stdout).should('be.visible');
    });
  });

  it('namespace dropdown', () => {
    cy.get('input[role="combobox"]').click();

    cy.get('span[class="ng-option-label ng-star-inserted"]')
      .contains('octant-cypress')
      .parent()
      .click()

    cy.location('hash').should('include', '/' + 'octant-cypress');
  });
});

let appHasStarted
function spyOnAddEventListener (win) {
  const addListener = win.EventTarget.prototype.addEventListener
  win.EventTarget.prototype.addEventListener = function (name) {
    console.log('Event listener added:', name)
    if (name === 'test') {
      // that means the web application has started
      appHasStarted = true
      win.EventTarget.prototype.addEventListener = addListener
    }
    return addListener.apply(this, arguments)
  }
}
function waitForAppStart() {
  return new Cypress.Promise((resolve, reject) => {
    const isReady = () => {
      if (appHasStarted) {
        return resolve();
      }
      setTimeout(isReady, 0);
    }
    isReady();
  });
};
