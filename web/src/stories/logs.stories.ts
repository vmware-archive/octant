import { storiesOf } from '@storybook/angular';
import {LogsComponent} from "../app/modules/shared/components/smart/logs/logs.component";


const view = {
  metadata: {
    type: "logs",
    title: [{metadata: {type: "text"}, config: {value: "Logs"}}],
    accessor: "logs"
  },
  config: {namespace: "default", name: "default-name-57466fd965-xprw9", containers: [""]}
};

const containerLogs=  [
  {
    timestamp: '2020-06-02T11:42:36.554540433Z',
    message: 'Here is a sample message',
    container: 'test-container',
  },
    {
      timestamp: '2020-06-02T11:44:19.554540433Z',
      message: 'that somehow',
      container: 'test-container',
    },
    {
      timestamp: '2020-06-02T12:59:06.554540433Z',
      message: 'showed up in this log',
      container: 'test-container',
    },
  ];
storiesOf('Logs', module).add('with logs', () => ({
  props: {
    view: view,
    containerLogs: containerLogs
  },
  component: LogsComponent,
}));
