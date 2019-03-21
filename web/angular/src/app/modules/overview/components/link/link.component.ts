import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LinkView } from 'src/app/models/content';

@Component({
  selector: 'app-view-link',
  templateUrl: './link.component.html',
  styleUrls: ['./link.component.scss'],
})
export class LinkComponent implements OnChanges {
  @Input() view: LinkView;

  ref: string;
  value: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LinkView;
      this.ref = view.config.ref;
      this.value = view.config.value;
    }
  }
}
