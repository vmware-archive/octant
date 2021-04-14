import { Component, OnInit, OnDestroy, HostListener } from '@angular/core';
import '@cds/core/button/register.js';
import '@cds/core/modal/register';
import { ClarityIcons, helpIcon } from '@cds/core/icon';
import { Subscription } from 'rxjs';
import { HelperService } from '../../../services/helper/helper.service';
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
  private buildInfoSubscription: Subscription;

  constructor(private helperService: HelperService) {
    ClarityIcons.addIcons(helpIcon);
  }

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

  toggleReleases(): void {
    this.getReleaseInfo(this.version);
    const releaseModal = document.getElementById('release-modal');
    releaseModal.hidden = !releaseModal.hidden;
  }

  showDocs(): void {
    window.open('https://octant.dev/', '_blank');
  }

  ngOnDestroy(): void {
    if (this.buildInfoSubscription) {
      this.buildInfoSubscription.unsubscribe();
    }
  }

  toggleBuildInfo(): void {
    const buildModal = document.getElementById('build-modal');
    buildModal.hidden = !buildModal.hidden;
  }

  toggleShortcut(): void {
    const shortcutModal = document.getElementById('shortcut-modal');
    shortcutModal.hidden = !shortcutModal.hidden;
  }

  @HostListener('window:keydown', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (event.ctrlKey && event.key === '/') {
      event.preventDefault();
      event.cancelBubble = true;
      this.toggleShortcut();
    }
  }
}
