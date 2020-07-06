import { Injectable, Renderer2, RendererFactory2 } from '@angular/core';
import {
  lightTheme,
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
        this.themeService.loadCSS(lightTheme.assetPath);
        this.monacoService.changeTheme('vs');
        this.renderer.addClass(document.body, 'light');
        resolve();
      } catch (e) {
        reject(e);
      }
    });
  }
}
