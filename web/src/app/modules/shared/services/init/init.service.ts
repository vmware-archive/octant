import { Injectable } from '@angular/core';
import { ThemeService } from '../theme/theme.service';
import { MonacoProviderService } from 'ng-monaco-editor';

@Injectable({
  providedIn: 'root',
})
export class InitService {
  private syncMonacoTheme: () => void;

  constructor(
    private themeService: ThemeService,
    private monacoService: MonacoProviderService
  ) {
    // we want a new instance of the handler for each component instance
    this.syncMonacoTheme = () => {
      const theme = this.themeService.isLightThemeEnabled() ? 'vs' : 'vs-dark';
      this.monacoService.changeTheme(theme);
    };
  }

  init(): void {
    this.themeService.loadTheme();
    // TODO remove this once we are able to define the theme before loading monaco
    this.monacoService.initMonaco().then(this.syncMonacoTheme);

    this.themeService.onChange(this.syncMonacoTheme);
    this.syncMonacoTheme();
  }
}
