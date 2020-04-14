import { Component } from '@angular/core';
import {
  ListView,
  TextView,
  TableView,
} from '../../../../../../../src/app/modules/shared/models/content';

const title: TextView = {
  metadata: {
    type: 'text',
  },
  config: {
    value: 'Title',
  },
};

const sampleText: TextView = {
  config: {
    value: 'sample text',
  },
  metadata: {
    type: 'text',
  },
};

const tableView: TableView = {
  metadata: {
    type: 'table',
    title: [title],
  },
  config: {
    columns: [
      {
        name: 'ColumnA',
        accessor: 'ColumnA',
      },
      {
        name: 'ColumnB',
        accessor: 'ColumnB',
      },
    ],
    rows: [
      {
        ColumnA: sampleText,
        ColumnB: sampleText,
      },
      {
        ColumnA: sampleText,
        ColumnB: sampleText,
      },
    ],
    emptyContent: 'There are no items!',
    loading: false,
    filters: {},
  },
};

const view: ListView = {
  config: {
    iconName: 'test',
    items: [tableView, tableView],
  },
  metadata: {
    type: 'list',
  },
};

const code = `cols := component.NewTableCols("ColumnA", "ColumnB")

component.NewList([]component.TitleComponent{}, []component.Component{
	component.NewList([]component.TitleComponent{}, []component.Component{
		component.NewTableWithRows("Title", "There are no items!", cols, []component.TableRow{
			// Table data
		})}),
})
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-list-demo',
  templateUrl: './angular-list.demo.html',
})
export class AngularListDemoComponent {
  view = view;
  code = code;
  json = json;
}
