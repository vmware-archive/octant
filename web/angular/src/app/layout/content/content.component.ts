import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subject } from 'rxjs';
import { DataService } from 'src/app/data.service';
import { ContentResponse, View } from 'src/app/models/content';
import { titleAsText } from 'src/app/util/view';

@Component({
  selector: 'app-content',
  templateUrl: './content.component.html',
  styleUrls: ['./content.component.scss'],
})
export class ContentComponent implements OnInit, OnDestroy {
  hasTabs = false;

  title: string = null;
  views: View[] = null;
  singleView: View = null;

  hasReceivedContent = false;
  constructor(private route: ActivatedRoute, private dataService: DataService) {}

  currentPath: Subject<string> = new Subject<string>();

  private previousUrl = '';

  ngOnInit() {
    this.route.url.subscribe((url) => {
      const currentPath = `${url.map((u) => u.path).join('/')}`;
      if (currentPath !== this.previousUrl) {
        this.previousUrl = currentPath;
        this.dataService.startPoller(currentPath);
        this.dataService.pollContent().subscribe((contentResponse: ContentResponse) => {
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
    this.dataService.stopPoller();
  }
}
