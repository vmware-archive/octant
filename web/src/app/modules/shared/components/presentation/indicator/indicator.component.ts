import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View } from '../../../models/content';
import '@cds/core/icon/register.js';
import {
  ClarityIcons,
  checkCircleIcon,
  exclamationCircleIcon,
  infoCircleIcon,
} from '@cds/core/icon';

/**
 * Status are statuses known to indicator.
 */
export const Status = {
  Ok: 1,
  Warning: 2,
  Error: 3,
};

export const statusLookup = {
  [Status.Ok]: 'success',
  [Status.Warning]: 'warning',
  [Status.Error]: 'danger',
};

export const iconLookup = {
  [Status.Ok]: 'check-circle',
  [Status.Warning]: 'info-circle',
  [Status.Error]: 'exclamation-circle',
};

@Component({
  selector: 'app-indicator',
  templateUrl: './indicator.component.html',
})
export class IndicatorComponent implements OnChanges {
  @Input()
  status: number;

  @Input()
  detail: View;

  currentStatus: string;
  iconShape: string;

  constructor() {
    ClarityIcons.addIcons(
      checkCircleIcon,
      exclamationCircleIcon,
      infoCircleIcon
    );
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.status) {
      this.currentStatus = statusLookup[changes.status.currentValue];
      this.iconShape = iconLookup[changes.status.currentValue];
    }
  }

  name(): string {
    return statusLookup[this.status];
  }
}
