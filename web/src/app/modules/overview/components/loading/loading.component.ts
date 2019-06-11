import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LoadingView } from 'src/app/models/content';

@Component({
  selector: 'app-view-loading',
  templateUrl: './loading.component.html',
  styleUrls: ['./loading.component.scss'],
})
export class LoadingComponent implements OnChanges {
  @Input() view: LoadingView;

  value: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LoadingView;
      this.value = view.config.value;
    }
  }
}
