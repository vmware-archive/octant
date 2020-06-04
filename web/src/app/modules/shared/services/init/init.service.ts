import { Injectable, Renderer2, RendererFactory2 } from '@angular/core';
import {
  darkTheme,
  ThemeService,
} from '../../../sugarloaf/components/smart/theme-switch/theme-switch.service';
import { MonacoProviderService } from 'ng-monaco-editor';

@Injectable()
export class InitService {
  private renderer: Renderer2;

  constructor(
    rendererFactory: RendererFactory2,
    private themeService: ThemeService,
    private monacoService: MonacoProviderService
  ) {
    this.renderer = rendererFactory.createRenderer(null, null);
  }

  init(): Promise<any> {
    return new Promise((resolve, reject) => {
      try {
        this.themeService.loadCSS(darkTheme.assetPath);
        this.monacoService.changeTheme('vs-dark');
        this.renderer.addClass(document.body, 'dark');
        resolve();
      } catch (e) {
        reject(e);
      }
    });
  }
}
