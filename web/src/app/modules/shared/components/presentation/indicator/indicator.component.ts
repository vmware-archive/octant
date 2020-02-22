import { Component, Input, OnInit } from '@angular/core';

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
export class IndicatorComponent implements OnInit {
  @Input()
  status: number;

  constructor() {}

  ngOnInit(): void {}

  name() {
    return statusLookup[this.status];
  }
}
