import { Component } from '@angular/core';
import { TimestampView } from '../../../../../../../src/app/modules/shared/models/content';

const view: TimestampView = {
  config: {
    timestamp: 1586469079,
  },
  metadata: {
    type: 'timestamp',
  },
};

const code = `timestamp := component.NewTimestamp(time.Unix(1586469079, 0))`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-timestamp-demo',
  templateUrl: './angular-timestamp.demo.html',
})
export class AngularTimestampDemoComponent {
  view = view;
  code = code;
  json = json;
}
