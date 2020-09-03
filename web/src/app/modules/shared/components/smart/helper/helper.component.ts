import { Component, OnInit, OnDestroy, HostListener } from '@angular/core';
import { Subscription } from 'rxjs';
import { HelperService } from '../../../../shared/services/helper/helper.service';

@Component({
  selector: 'app-helper',
  templateUrl: './helper.component.html',
  styleUrls: ['./helper.component.scss'],
})
export class HelperComponent implements OnInit, OnDestroy {
  version = '';
  commit = '';
  time = '';
  isBuildModalOpen = false;
  isShortcutModalOpen = false;
  private buildInfoSubscription: Subscription;

  constructor(private helperService: HelperService) {}

  ngOnInit() {
    this.buildInfoSubscription = this.helperService
      .buildVersion()
      .subscribe(version => (this.version = version));
    this.buildInfoSubscription = this.helperService
      .buildCommit()
      .subscribe(commit => (this.commit = commit));
    this.buildInfoSubscription = this.helperService
      .buildTime()
      .subscribe(time => (this.time = time));
  }

  ngOnDestroy(): void {
    this.buildInfoSubscription.unsubscribe();
  }

  @HostListener('window:keyup', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (this.isShortcutModalOpen) {
      return;
    }
    if (event.ctrlKey && event.key === '/') {
      this.isShortcutModalOpen = true;
    }
  }
}
