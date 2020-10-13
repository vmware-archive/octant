import { Component, OnInit, OnDestroy, HostListener } from '@angular/core';
import { Subscription } from 'rxjs';
import { HelperService } from '../../../../shared/services/helper/helper.service';
import { TextView } from '../../../models/content';

@Component({
  selector: 'app-helper',
  templateUrl: './helper.component.html',
  styleUrls: ['./helper.component.scss'],
})
export class HelperComponent implements OnInit, OnDestroy {
  version = '';
  commit = '';
  time = '';
  releaseInfo: TextView = {
    metadata: {
      type: 'text',
    },
    config: {
      value: '',
      isMarkdown: true,
    },
  };
  isBuildModalOpen = false;
  isShortcutModalOpen = false;
  isReleaseModalOpen = false;
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

  openIssue(): void {
    window.open(
      'https://github.com/vmware-tanzu/octant/issues/new/choose',
      '_blank'
    );
  }

  getReleaseInfo(version: string) {
    const ver = version.substring(0, version.lastIndexOf('.'));
    const baseUrl =
      'https://raw.githubusercontent.com/vmware-tanzu/octant/master';
    const url =
      ver.length > 0
        ? `${baseUrl}/changelogs/CHANGELOG-${ver}.md`
        : `${baseUrl}/CHANGELOG.md`;

    fetch(url)
      .then(response => response.text())
      .then(data => {
        this.releaseInfo.config.value = data;
      });
  }

  showReleases(): void {
    this.getReleaseInfo(this.version);
    this.isReleaseModalOpen = true;
  }

  showDocs(): void {
    window.open('https://octant.dev/', '_blank');
  }

  ngOnDestroy(): void {
    if (this.buildInfoSubscription) {
      this.buildInfoSubscription.unsubscribe();
    }
  }

  @HostListener('window:keydown', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (this.isShortcutModalOpen) {
      return;
    }
    if (event.ctrlKey && event.key === '/') {
      event.preventDefault();
      event.cancelBubble = true;
      this.isShortcutModalOpen = true;
    }
  }
}
