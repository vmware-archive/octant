import { Component, OnInit, Renderer2 } from '@angular/core';
import {
  darkTheme,
  ThemeService,
} from '../../../../sugarloaf/components/smart/theme-switch/theme-switch.service';

@Component({
  selector: 'app-root',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {
  constructor(
    private renderer: Renderer2,
    private themeService: ThemeService
  ) {}

  ngOnInit() {
    this.loadTheme();
  }

  loadTheme() {
    // TODO: enable theme switching or denali
    this.themeService.loadCSS(darkTheme.assetPath);

    this.renderer.addClass(document.body, 'dark');
  }
}
