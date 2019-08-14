import { Injectable } from '@angular/core';
import { Router, Resolve } from '@angular/router';
import { Observable } from 'rxjs';
import { NamespaceService } from 'src/app/services/namespace/namespace.service';
import { Namespace } from 'src/app/models/namespace';

@Injectable({
   providedIn: 'root'
})
export class NamespaceResolver implements Resolve<Namespace> {
  namespaces: string[];
  initialNamespace: string;
  url = '/content/overview/namespace/'

  constructor(
    private namespaceService: NamespaceService,
    private router: Router,
  ) {
    this.namespaceService.getInitialNamespace().subscribe((initial: Namespace) => {
      this.initialNamespace = initial.namespace;
      this.router.navigate([this.url + this.initialNamespace]);
    });
  }

  resolve(): Observable<any> {
    this.router.navigate([this.url + this.namespaceService.current.getValue()]);
    return;
  }
}