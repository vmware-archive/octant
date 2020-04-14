import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { IFrameView, View } from '../../../../shared/models/content';

@Component({
  selector: 'app-iframe',
  templateUrl: './iframe.component.html',
  styleUrls: ['./iframe.component.scss'],
})
export class IFrameComponent implements OnChanges {
  private v: IFrameView;

  @Input() set view(v: View) {
    this.v = v as IFrameView;
  }
  get view() {
    return this.v;
  }

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
