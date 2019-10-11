describe('Octant Integration Tests', () => {
    var namespace = 'octant-cypress-' +  random_string()
    var startingContext = ''

    before(function() {
        cy.exec('kubectl config current-context').then((result) => {
            startingContext = result.stdout
        })

        cy.exec('kubectl create namespace ' + namespace)
        cy.exec('kubectl config set-context $CYPRESS_CONTEXT --namespace ' + namespace +
        ' --cluster $(kubectl config get-contexts ' + startingContext + ' | tail -1 | awk \'{print $3}\') \
        --user $(kubectl config get-contexts ' + startingContext + ' | tail -1 | awk \'{print $4}\')',
         { env: { CYPRESS_CONTEXT: Cypress.env('CYPRESS_CONTEXT')} })
        cy.exec('kubectl config use-context $CYPRESS_CONTEXT', { env: { CYPRESS_CONTEXT: Cypress.env('CYPRESS_CONTEXT')} })
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
        cy
            .location('hash')
            .should('include', '/' + namespace)
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
            .exec('kubectl apply -f ../examples/resources/deployment.yaml --namespace ' + namespace)
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
        cy.exec('kubectl delete pod -l app.kubernetes.io/name=nginx,app.kubernetes.io/instance=sample,app.kubernetes.io/version=v1 --namespace ' + namespace)
            .its('stdout')
            .should('contain', 'deleted')
    
        cy.contains(/Stop port forward/).should('not.exist')
    })

    it('navigate to title', () => {
        cy.get('[class="title"]').click()
    })

    it('check plugin tab', () => {
        cy.get('[href="#/overview/namespace/' + namespace + '/workloads/pods"]').click()
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

    it('check resource viewer', () => {
        cy.contains(/Resource Viewer/).click()
        // Check canvas is drawn
        cy.get('[data-id="layer0-selectbox"]').invoke('width').should('be.greaterThan', 0)
        cy.get('[data-id="layer2-node"]').click(400, 320, {force: true})

        cy.get('app-heptagon-grid svg g:first').children().should('have.length', 3)
        cy.get('app-heptagon-grid svg g:first').children().last().click()

        cy.contains(/Container nginx/)
    })

    it('cleanup context and namespace', () => {
        cy.exec('kubectl config use-context ' + startingContext)
        cy.exec('kubectl delete namespace '  + namespace)
        cy.exec('kubectl config delete-context $CYPRESS_CONTEXT', { env: { CYPRESS_CONTEXT: Cypress.env('CYPRESS_CONTEXT')} })
    })

    function random_string() {
        var text = "";
        var chars = "abcdefghijklmnopqrstuvwxyz123456789";
        for (var i = 0; i < 6; i++)
          text += chars.charAt(Math.floor(Math.random() * chars.length));

        return text;
      }
})
