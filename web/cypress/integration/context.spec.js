describe('Context', () => {
  before(() => {
    cy.visit('/');
  });

  it('has kubeconfig context', () => {
    cy.contains(' octant-temporary ').click();

    cy.get('.active').contains(' octant-temporary ');
  });
});
