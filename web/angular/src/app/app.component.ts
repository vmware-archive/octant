import { Component, OnInit } from '@angular/core';

import { DataService } from './services/data.service';
import { Navigation } from './models/navigation';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  navigation: Navigation;

  constructor(private dataService: DataService) {}

  ngOnInit(): void {
    this.dataService.getNavigation().subscribe((navigation: Navigation) => {
      this.navigation = navigation;
    });
  }
}
