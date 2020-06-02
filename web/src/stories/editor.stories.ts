import { storiesOf } from '@storybook/angular';
import {EditorView, TextView} from "../app/modules/shared/models/content";
import {MonacoEditorModule} from "ng-monaco-editor";

const text: TextView = {
  config: {
    value: 'YAML',
  },
  metadata: {
    type: 'text',
  },
};

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
    title: [text]
  },
};

storiesOf('Text', module).add('editor', () => ({
  props: {
    view: view,
  },
  moduleMetadata: {
    imports: [
      MonacoEditorModule.forRoot({
        baseUrl: 'lib',
        defaultOptions: {},
      }),
    ],
  },  template: `
    <div class="main-container">
        <div class="content-container">
            <div class="content-area">
                <app-view-editor [view]="view"></app-view-editor>
            </div>
        </div>
    </div>
    `,
}));
