import { Component, OnInit, Renderer2 } from '@angular/core';
import { ThemeService } from 'src/app/modules/sugarloaf/components/smart/theme-switch/theme-switch.service';
import { ElectronService } from '../../../modules/shared/services/electron/electron.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.sass'],
})
export class HomeComponent implements OnInit {
  constructor(
    private renderer: Renderer2,
    private themeService: ThemeService,
    private electronService: ElectronService
  ) {}

  ngOnInit() {
    this.themeService.loadTheme();

    if (this.electronService.isElectron()) {
      this.renderer.addClass(document.body, 'electron');
      this.renderer.addClass(
        document.body,
        `platform-${this.electronService.platform()}`
      );
    }
  }
}
