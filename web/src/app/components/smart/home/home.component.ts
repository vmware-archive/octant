import { Component, OnInit, Renderer2 } from '@angular/core';
import { InitService } from '../../../modules/shared/services/init/init.service';
import { ElectronService } from '../../../modules/shared/services/electron/electron.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.sass'],
})
export class HomeComponent implements OnInit {
  constructor(
    private renderer: Renderer2,
    private initService: InitService,
    private electronService: ElectronService
  ) {}

  ngOnInit() {
    this.initService.init();

    if (this.electronService.isElectron()) {
      this.renderer.addClass(document.body, 'electron');
      this.renderer.addClass(
        document.body,
        `platform-${this.electronService.platform()}`
      );
    }
  }
}
