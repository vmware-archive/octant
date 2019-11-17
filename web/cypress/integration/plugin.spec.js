describe('Plugin', () => {
  before(() => {
    cy.visit('/');
  });

  it('click plugin tab', () => {
    // TODO (GuessWhoSamFoo) Start octant with the temp context to avoid passing namespaces
    cy.exec(
      `kubectl config view --minify --output 'jsonpath={..namespace}'`
    ).then(result => {
      cy.exec(
        'kubectl apply -f ../examples/resources/log-pod.yaml --namespace ' +
          result.stdout
      );
      cy.visit('/#/overview/namespace/' + result.stdout + '/workloads/pods');

      cy.location('hash').should('include', '/' + 'pods');

      cy.contains(/logme/).click();

      cy.get('[class=nav]')
        .find('button')
        .should('have.length', 6)
        .last()
        .contains('Extra Pod Details')
        .click();

      cy.contains('content');
    });
  });
});
