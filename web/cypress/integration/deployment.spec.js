describe('Deployment', () => {
  before(() => {
    cy.visit('/');
  });

  it('port forward nginx deployment', () => {
    // TODO (GuessWhoSamFoo) Start octant with temp context to avoid passing namespaces (GH#362)
    cy.exec(
      `kubectl config view --minify --output 'jsonpath={..namespace}'`
    ).then(result => {
      cy.exec(
        'kubectl apply -f ../examples/resources/deployment.yaml --namespace ' +
          result.stdout
      )
        .its('stdout')
        .should('contains', 'nginx-deployment');

      cy.visit(
        '/#/overview/namespace/' + result.stdout + '/workloads/deployments'
      );

      cy.get('span[class="ng-value-label ng-star-inserted"]').should('contain', result.stdout);

      cy.get('[class="ng-star-inserted"]')
        .contains('nginx-deployment')
        .click();
      cy.contains(/^nginx-deployment-[a-z0-9]+-[a-z0-9]+/).click();

      cy.location('hash').should('include', '/' + 'pods');

      cy.contains('Running').should('be.visible')
      cy.contains(/Start port forward/)
        .click();

      cy.get('a[class="open-pf ng-star-inserted"')
        .contains('localhost');

      cy.get('[class=port-actions]')
        .should('have.length', 1)
        .first()
        .contains('Stop port forward');

      cy.contains(/Port Forwards/).click();

      cy.contains(/Stop port forward/).should('be.visible');
      cy.exec(
        'kubectl delete pod -l app.kubernetes.io/name=nginx,app.kubernetes.io/instance=sample,app.kubernetes.io/version=v1 --namespace ' +
          result.stdout
      )
        .its('stdout')
        .should('contain', 'deleted');

      cy.contains(/Stop port forward/).should('not.exist');
    });
  });

  it('resource viewer to pod', () => {
    cy.exec(
      `kubectl config view --minify --output 'jsonpath={..namespace}'`
    ).then(result => {
      cy.visit(
        '/#/overview/namespace/' + result.stdout + '/workloads/deployments'
      );

      cy.get('span[class="ng-value-label ng-star-inserted"]').should('contain', result.stdout);

      cy.get('[class="ng-star-inserted"]')
        .contains('nginx-deployment')
        .click();
      cy.contains(/^nginx-deployment-[a-z0-9]+-[a-z0-9]+/).click();

      cy.contains(/Resource Viewer/).click();
      // Check canvas is drawn
      cy.get('[data-id="layer0-selectbox"]')
        .invoke('width')
        .should('be.greaterThan', 0);
      cy.get('[data-id="layer2-node"]').click(370, 360, { force: true });

      cy.get('app-heptagon-grid svg g:first')
        .children()
        .should('have.length', 3);
      cy.get('app-heptagon-grid svg g:first')
        .children()
        .last()
        .click();

      cy.contains(/Container nginx/);
    });
  });
});
