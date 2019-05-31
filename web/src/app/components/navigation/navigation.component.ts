import { Component, Input } from '@angular/core';
import { Navigation, NavigationChild } from '../../models/navigation';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {
  @Input() navigation: Navigation = {
    sections: []
  };

  constructor() { }

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }
}
