import { Component, Input, ViewEncapsulation } from '@angular/core';
import { Node } from 'src/app/models/content';

@Component({
  selector: 'app-view-object-status',
  templateUrl: './object-status.component.html',
  styleUrls: ['./object-status.component.scss'],
  encapsulation: ViewEncapsulation.Emulated,
})
export class ObjectStatusComponent {
  @Input() node: Node;

  constructor() { }

  indicatorClass() {
    if (!this.node) {
      return ['progress', 'top', 'success'];
    }

    return [
      'progress', 'top',
      this.node.status === 'ok' ? 'success' : 'danger',
    ];
  }

  detailsTrackBy(index, item) {
    return index;
  }
}
