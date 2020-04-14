import { Component } from '@angular/core';
import { ButtonGroupView } from '../../../../../../../src/app/modules/shared/models/content';

const buttonGroupView: ButtonGroupView = {
  config: {
    buttons: [
      {
        name: 'Delete',
        payload: {},
        confirmation: {
          title: 'Delete Pod',
          body: 'Are you sure you want to delete *Pod* **pod**?',
        },
      },
    ],
  },
  metadata: {
    type: 'buttonGroup',
  },
};

const multipleButtonsView: ButtonGroupView = {
  config: {
    buttons: [
      {
        name: 'Button 1',
        payload: {},
      },
      {
        name: 'Button 2',
        payload: {},
      },
    ],
  },
  metadata: {
    type: 'buttonGroup',
  },
};

const buttonGroupCode = `buttonGroup := component.NewButtonGroup()

buttonGroup.AddButton(
    component.NewButton("Delete",
      action.CreatePayload(octant.ActionDeleteObject, key.ToActionPayload()),
      component.WithButtonConfirmation(
        "Delete Pod",
        "Are you sure you want to delete *Pod* **pod**?",
      )))
`;

const multipleButtonsCode = `buttonGroup := component.NewButtonGroup()
buttonGroup.AddButton(component.NewButton("Button 1", nil, nil))
buttonGroup.AddButton(component.NewButton("Button 2", nil, nil))
`;

const buttonGroupJSON = JSON.stringify(buttonGroupView, null, 4);

const multipleButtonsJSON = JSON.stringify(multipleButtonsView, null, 4);

@Component({
  selector: 'app-angular-button-group-demo',
  templateUrl: './angular-button-group.demo.html',
})
export class AngularButtonGroupDemoComponent {
  buttonGroupView = buttonGroupView;
  multipleButtonsView = multipleButtonsView;
  buttonGroupCode = buttonGroupCode;
  multipleButtonsCode = multipleButtonsCode;
  buttonGroupJSON = buttonGroupJSON;
  multipleButtonsJSON = multipleButtonsJSON;
}
