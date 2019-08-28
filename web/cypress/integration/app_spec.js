describe('Octant Integration Tests', () => {
  beforeEach(function(){
    // baseURL should be localhost:7777
    cy.server()
    cy.route('/api/v1/namespace').as('overview')
    cy.route({
      method: 'POST',
      url: '/api/v1/content/overview/port-forwards',
    }).as('portforward')
    cy.visit('/')
  })

  it('loads page', () => {
    cy.title().should('include', 'Octant')
  })

  it('has kubeconfig context', () => {
    cy.wait('@overview')
    cy.contains(' kubernetes-admin@kind ').click()

    cy
    .get('.active')
    .should('have.length', 2)
    .first()
    .contains(' kubernetes-admin@kind ')
  })

  it('create and delete port forward', () => {
    // Find first nginx pod in overview
    cy
    .exec('kubectl apply -f ../examples/resources/deployment.yaml')
    .its('stdout')
    .should('contains', 'nginx-deployment')

    cy.contains(/^nginx-deployment-[a-z0-9]+-[a-z0-9]+/).click()
    cy.contains('Running').should('be.visible')

    cy.contains(/Start port forward/).click().should('not.exist')
    cy.wait('@portforward')

    cy
    .get('[class=port-actions]')
    .should('have.length', 1)
    .first()
    .contains('Stop port forward')

    cy.contains(/Port Forwards/).click()

    cy.contains(/Stop port forward/).should('be.visible')
    cy.exec('kubectl delete pod -l app=nginx-deployment')
    .its('stdout')
    .should('contain', 'deleted')

    cy.contains(/Stop port forward/).should('not.exist')
  })
})
