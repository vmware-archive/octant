import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export interface ContextDescription {
  name: string;
}

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

  constructor() {
    this.contextsSource.next([
      {
        name: 'kubernetes-admin@service-account',
      },
      {
        name: 'kubernetes-admin@workload-test',
      },
    ]);

    this.selectedSource.next('kubernetes-admin@service-account');
  }

  select(context: ContextDescription) {
    console.log(`settings active context to ${context.name}`);
    this.selectedSource.next(context.name);

    // TODO: let backend know the context has changed
  }

  selected() {
    return this.selectedSource.asObservable();
  }

  contexts() {
    return this.contextsSource.asObservable();
  }
}
