describe('Namespace', () => {
  before(() => {
    cy.visit('/');
  });

  it('namespace dropdown', () => {
    cy.get('input[role="combobox"]').click();

    cy.contains('octant-cypress').click();

    cy.location('hash').should('include', '/' + 'octant-cypress');
  });
});
