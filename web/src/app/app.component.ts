import { Component, OnInit } from '@angular/core';
import { ContentStreamService } from './services/content-stream/content-stream.service';
import { Navigation } from './models/navigation';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  navigation: Navigation;

  constructor(private contentStreamService: ContentStreamService) {}

  ngOnInit(): void {
    this.contentStreamService.navigation.subscribe((navigation: Navigation) => {
      this.navigation = navigation;
    });
  }
}
