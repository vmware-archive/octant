import { Component, OnInit, Renderer2 } from '@angular/core';
import {
  darkTheme,
  ThemeService,
} from '../../../../shared/services/theme/theme.service';
import { MonacoProviderService } from 'ng-monaco-editor';

@Component({
  selector: 'app-root',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {
  constructor(
    private renderer: Renderer2,
    private themeService: ThemeService,
    private monacoService: MonacoProviderService
  ) {}

  ngOnInit() {
    this.loadTheme();
  }

  loadTheme() {
    // TODO: enable theme switching or denali
    this.themeService.loadCSS(darkTheme.assetPath);
    this.monacoService.changeTheme('vs-dark');
    this.renderer.addClass(document.body, 'dark');
  }
}
