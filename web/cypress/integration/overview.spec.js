describe('Overview page', () => {
  it('loads page', () => {
    cy.visit('/');
    cy.title().should('include', 'Octant');
  });

  it('navigate to title', () => {
    cy.get('[class="title"]').click();
  });
});
