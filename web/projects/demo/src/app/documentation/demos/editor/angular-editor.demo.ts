import { Component } from '@angular/core';
import { EditorView } from '../../../../../../../src/app/modules/shared/models/content';

const view: EditorView = {
  config: {
    value:
      '---\nmetadata:\n  creationTimestamp: null\nspec:\n  containers: null\nstatus: {}\n',
    language: 'yaml',
    readOnly: false,
    metadata: {
      apiVersion: 'v1',
      kind: 'Pod',
      name: 'nginx-deployment-75c5cb5f44-7wh55',
      namespace: 'default',
    },
  },
  metadata: {
    type: 'editor',
  },
};
const code = `data := "---\nmetadata:\n  creationTimestamp: null\nspec:\n  containers: null\nstatus: {}\n"
component.NewEditor(component.TitleFromString("YAML"), data, false)
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-editor-demo',
  templateUrl: './angular-editor.demo.html',
  styles: [':host ::ng-deep .editor-container .editor { height: 500px; }'],
})
export class AngularEditorDemoComponent {
  view = view;
  code = code;
  json = json;
}
