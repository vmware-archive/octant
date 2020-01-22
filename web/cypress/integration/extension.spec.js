describe('Extension', () => {
  before(() => {
    cy.visit('/');
  });

  it('container exec creates tab', () => {
    cy.exec(
      `kubectl config view --minify --output 'jsonpath={..namespace}'`
    ).then(result => {
      cy.exec(
        'kubectl apply -f ../examples/resources/log-pod.yaml --namespace ' +
          result.stdout
      )
        .its('stdout')
        .should('contains', 'logme');

      cy.visit(
        '/#/overview/namespace/' + result.stdout + '/workloads/pods/logme'
      );

      cy.contains(/^logme/).click();

      cy.get('[class="btn btn-sm btn-link ng-star-inserted"]')
        .contains('Execute Command')
        .click();

      cy.get('[class="clr-input ng-untouched ng-pristine ng-valid"]').type('bash{enter}');

      cy.get('[class="tab-button btn btn-link nav-link active"]').should('contain', ' bash ');

      cy.exec(
        'kubectl delete pod logme --namespace ' + result.stdout
      )
        .its('stdout')
        .should('contain', 'deleted');
    });
  });
})