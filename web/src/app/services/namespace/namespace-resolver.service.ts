import { Injectable } from '@angular/core';
import { Router, Resolve } from '@angular/router';
import { Observable } from 'rxjs';
import { NamespaceService } from 'src/app/services/namespace/namespace.service';
import { Namespace } from 'src/app/models/namespace';
import { take, takeLast } from 'rxjs/operators';
import { ContentService } from '../../modules/overview/services/content/content.service';

@Injectable({
  providedIn: 'root',
})
export class NamespaceResolver implements Resolve<Namespace> {
  namespaces: string[];
  initialNamespace: string;
  url = '/content/overview/namespace/';

  constructor(
    private namespaceService: NamespaceService,
    private router: Router
  ) {}

  resolve(): Observable<any> {
    this.router.navigate([
      this.url + this.namespaceService.activeNamespace.getValue(),
    ]);
    return;
  }
}
