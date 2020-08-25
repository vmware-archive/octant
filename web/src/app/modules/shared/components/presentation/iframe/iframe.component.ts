import { Component } from '@angular/core';
import { IFrameView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-iframe',
  templateUrl: './iframe.component.html',
  styleUrls: ['./iframe.component.scss'],
})
export class IFrameComponent extends AbstractViewComponent<IFrameView> {
  url: string;
  title: string;

  constructor() {
    super();
  }

  update() {
    this.url = this.v.config.url;
    this.title = this.v.config.title;
  }
}
