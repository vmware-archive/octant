import { Component, OnInit, OnDestroy } from '@angular/core';
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
  isModalOpen = false;
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
}
