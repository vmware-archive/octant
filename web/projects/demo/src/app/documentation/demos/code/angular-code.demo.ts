import { Component } from '@angular/core';
import {
  CodeView,
  TableView,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const codeView: CodeView = {
  config: {
    value:
      'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("hello world")\n}',
  },
  metadata: {
    type: 'codeBlock',
  },
};

const titleView: TextView = {
  metadata: {
    type: 'text',
  },
  config: {
    value: 'Data',
  },
};

const textView: TextView = {
  metadata: {
    type: 'text',
  },
  config: {
    value: 'example.go',
  },
};

const view: TableView = {
  metadata: {
    type: 'table',
    title: [titleView],
  },
  config: {
    columns: [
      {
        name: 'Key',
        accessor: 'Key',
      },
      {
        name: 'Value',
        accessor: 'Value',
      },
    ],
    rows: [
      {
        Key: textView,
        Value: codeView,
      },
    ],
    emptyContent: 'There are no values!',
    loading: false,
    filters: {},
  },
};

const code = `component.NewCodeBlock("package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("hello world")\n}")
`;

const json = JSON.stringify(codeView, null, 4);

@Component({
  selector: 'app-angular-code-demo',
  templateUrl: './angular-code.demo.html',
})
export class AngularCodeDemoComponent {
  view = view;
  code = code;
  json = json;
}
