import { Component, OnInit, Renderer2 } from '@angular/core';
import { ThemeService } from 'src/app/modules/sugarloaf/components/smart/theme-switch/theme-switch.service';
import { MonacoProviderService } from 'ng-monaco-editor';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.sass'],
})
export class HomeComponent implements OnInit {
  constructor(
    private monacoService: MonacoProviderService,
    private renderer: Renderer2,
    private themeService: ThemeService
  ) {}

  ngOnInit() {
    this.themeService.loadTheme(this.monacoService, this.renderer);
  }
}
