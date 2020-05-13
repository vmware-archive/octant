import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View } from '../../../models/content';

/**
 * Status are statuses known to indicator.
 */
export const Status = {
  Ok: 1,
  Warning: 2,
  Error: 3,
};

export const statusLookup = {
  [Status.Ok]: 'ok',
  [Status.Warning]: 'warning',
  [Status.Error]: 'error',
};

@Component({
  selector: 'app-indicator',
  templateUrl: './indicator.component.html',
  styleUrls: ['./indicator.component.scss'],
})
export class IndicatorComponent implements OnChanges {
  @Input()
  status: number;

  @Input()
  detail: View;

  currentStatus: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.status) {
      this.currentStatus = statusLookup[changes.status.currentValue];
    }
  }

  name(): string {
    return statusLookup[this.status];
  }
}
