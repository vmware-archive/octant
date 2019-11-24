import { Component, OnChanges, Input, SimpleChanges } from '@angular/core';
import { IFrameView } from 'src/app/models/content';

@Component({
  selector: 'app-iframe',
  templateUrl: './iframe.component.html',
  styleUrls: ['./iframe.component.scss'],
})
export class IFrameComponent implements OnChanges {
  @Input() view: IFrameView;

  url: string;
  title: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as IFrameView;
      this.url = view.config.url;
      this.title = view.config.title;
    }
  }
}
