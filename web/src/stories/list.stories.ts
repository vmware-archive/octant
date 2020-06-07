import {ListView, TableView, TextView} from "../app/modules/shared/models/content";
import {storiesOf} from "@storybook/angular";
import {ListComponent} from "../app/modules/shared/components/presentation/list/list.component";

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

const listView: ListView = {
  config: {
    iconName: 'test',
    items: [tableView],
  },
  metadata: {
    type: 'list',
  },
};

storiesOf('List', module).add('List component', () => ({
  props: {
    view: listView
  },
  component: ListComponent,
}));
