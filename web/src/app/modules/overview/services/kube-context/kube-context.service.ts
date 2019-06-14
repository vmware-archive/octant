import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import getAPIBase from '../../../../services/common/getAPIBase';
import {
  ContentStreamService,
  ContextDescription,
} from '../../../../services/content-stream/content-stream.service';

@Injectable({
  providedIn: 'root',
})
export class KubeContextService {
  private contextsSource: BehaviorSubject<
    ContextDescription[]
  > = new BehaviorSubject<ContextDescription[]>([]);

  private selectedSource: BehaviorSubject<string> = new BehaviorSubject<string>(
    ''
  );

  constructor(
    private http: HttpClient,
    private contentStream: ContentStreamService
  ) {
    contentStream.kubeContext.subscribe(update => {
      this.contextsSource.next(update.contexts);
      this.selectedSource.next(update.currentContext);
    });
  }

  select(context: ContextDescription) {
    this.selectedSource.next(context.name);

    this.updateContext(context.name).subscribe();
  }

  selected() {
    return this.selectedSource.asObservable();
  }

  contexts() {
    return this.contextsSource.asObservable();
  }

  private updateContext(name: string) {
    const url = [
      getAPIBase(),
      'api/v1/content/configuration',
      'kube-contexts',
    ].join('/');

    const payload = {
      requestedContext: name,
    };

    console.log('updating?', url);

    return this.http.post(url, payload);
  }
}
