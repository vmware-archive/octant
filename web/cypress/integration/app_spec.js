describe('Octant Integration Tests', () => {
    before(function() {
        cy.exec('kubectl create namespace octant-cypress')
        cy.exec('kubectl config set-context octant-temporary --namespace octant-cypress \
        --cluster $(kubectl config get-contexts $CURRENT_CONTEXT | tail -1 | awk \'{print $3}\') \
        --user $(kubectl config get-contexts $CURRENT_CONTEXT | tail -1 | awk \'{print $4}\')', { env: { CURRENT_CONTEXT: Cypress.env('CURRENT_CONTEXT' )} })
        cy.exec('kubectl config use-context octant-temporary')
    })

    it('loads page', () => {
        cy.visit('/')
        cy.title().should('include', 'Octant')
    })

    it('namespace dropdown', () => {
        cy
            .get('input[role="combobox"]')
            .click()

        cy
            .contains('octant-cypress')
            .click()
    })

    it('has kubeconfig context', () => {
        cy.contains(' octant-temporary ').click()
    
        cy
            .get('.active')
            .contains(' octant-temporary ')
    })
    
    it('create and delete port forward', () => {
        // Find first nginx pod in overview
        cy
            .exec('kubectl apply -f ../examples/resources/deployment.yaml --namespace octant-cypress')
            .its('stdout')
            .should('contains', 'nginx-deployment')
    
        cy.contains(/^nginx-deployment-[a-z0-9]+-[a-z0-9]+/).click()
        cy.contains('Running').should('be.visible')
    
        cy.contains(/Start port forward/).click().should('not.exist')
    
        cy
            .get('[class=port-actions]')
            .should('have.length', 1)
            .first()
            .contains('Stop port forward')
    
        cy.contains(/Port Forwards/).click()
    
        cy.contains(/Stop port forward/).should('be.visible')
        cy.exec('kubectl delete pod -l app=nginx-deployment --namespace octant-cypress')
            .its('stdout')
            .should('contain', 'deleted')
    
        cy.contains(/Stop port forward/).should('not.exist')
    })

    it('navigate to title', () => {
        cy.get('[class="title"]').click()
    })

    it('check plugin tab', () => {
        cy.contains(/^nginx-deployment-[a-z0-9]+-[a-z0-9]+/).click()
        cy
            .get('[class=nav]')
            .find('button')
            .should('have.length', 5)
            .last()
            .contains('Extra Pod Details')
    
        cy.contains(/Extra Pod Details/).click()
        cy.contains('content')
    })

    it('cleanup context and namespace', () => {
        cy.exec('kubectl config use-context $CURRENT_CONTEXT', { env: { CURRENT_CONTEXT: Cypress.env('CURRENT_CONTEXT' )} })
        cy.exec('kubectl delete namespace octant-cypress')
        cy.exec('kubectl config delete-context octant-temporary')
      })
})
