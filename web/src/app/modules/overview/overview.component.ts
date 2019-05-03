import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DataService } from 'src/app/services/data/data.service';
import { ContentResponse, View } from 'src/app/models/content';
import { titleAsText } from 'src/app/util/view';

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
  styleUrls: ['./overview.component.scss'],
})
export class OverviewComponent implements OnInit, OnDestroy {
  private previousUrl = '';
  hasTabs = false;
  hasReceivedContent = false;
  title: string = null;
  views: View[] = null;
  singleView: View = null;

  constructor(private route: ActivatedRoute, private dataService: DataService) {}

  ngOnInit() {
    this.route.url.subscribe((url) => {
      const currentPath = url.map((u) => u.path).join('/');
      if (currentPath !== this.previousUrl) {
        this.singleView = null;
        this.views = null;
        this.previousUrl = currentPath;
        this.dataService.openStream(currentPath);
        this.dataService.content.subscribe((contentResponse: ContentResponse) => {
          this.connect(contentResponse);
        });
      }
    });
  }

  connect(contentResponse: ContentResponse) {
    const views = contentResponse.content.viewComponents;
    if (views.length === 0) {
      this.hasReceivedContent = false;
      return;
    }

    this.hasTabs = views.length > 1;
    if (this.hasTabs) {
      this.views = views;
      this.title = titleAsText(contentResponse.content.title);
    } else if (views.length === 1) {
      this.singleView = views[0];
    }

    this.hasReceivedContent = true;
  }

  ngOnDestroy() {
    this.dataService.closeStream();
  }
}
