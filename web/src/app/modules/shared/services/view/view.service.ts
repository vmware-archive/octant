import { Injectable } from '@angular/core';
import { TextView, View } from '../../models/content';

@Injectable({
  providedIn: 'root',
})
export class ViewService {
  constructor() {}

  titleAsText(titleViews: View[]): string {
    if (!titleViews) {
      return '';
    }

    // assume it's a text title
    return titleViews
      .map((titleView: TextView) => titleView.config.value)
      .join(' / ');
  }

  viewTitleAsText(view: View): string {
    return this.titleAsText(view.metadata.title);
  }
}
