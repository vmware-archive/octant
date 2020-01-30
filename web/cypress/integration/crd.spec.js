describe('CRD', () => {
  before(() => {
    cy.visit('/');
  });
});

it('loads CRD', () => {
  cy.exec(
    `kubectl config view --minify --output 'jsonpath={..namespace}'`
  ).then(result => {
    cy.exec(
      'kubectl apply -f ../examples/resources/crd-crontab.yaml --namespace ' +
        result.stdout
    )
      .its('stdout')
      .should('contains', 'crontabs.stable.example.com created');

    cy.exec(
      'kubectl apply -f ../examples/resources/crd-crontab-resource.yaml --namespace ' +
        result.stdout
    )
      .its('stdout')
      .should(
        'contains',
        'crontab.stable.example.com/my-new-cron-object created'
      );

    cy.visit(
      '/#/overview/namespace/' +
        result.stdout +
        '/custom-resources/crontabs.stable.example.com'
    );

    cy.get('a[class="ng-star-inserted"]').contains('my-new-cron-object');

    cy.exec(
      'kubectl delete crd crontabs.stable.example.com --namespace ' +
        result.stdout
    )
      .its('stdout')
      .should('contain', 'deleted');

    cy.contains(/my-new-cron-object/).should('not.exist');
  });
});
