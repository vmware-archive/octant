describe('Octant Integration Tests', () => {
  beforeEach(function(){
    // baseURL should be localhost:7777
    cy.visit('/')
  })

  it('loads page', () => {
    cy.title().should('include', 'Octant')
  })

  it('has kubeconfig context', () => {
    cy.viewport(1440, 900)
    cy.get('[class=dropdown]').click()

    cy
    .get('.active')
    .should('have.length', 2)
    .first()
    .contains(' kubernetes-admin@kind ')
  })
})