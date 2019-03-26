import { Component, OnInit } from '@angular/core';
import { NamespaceService } from 'src/app/services/namespace/namespace.service';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
})
export class NamespaceComponent implements OnInit {
  namespaces: string[];
  currentNamespace: string;

  constructor(private namespaceService: NamespaceService) {}

  ngOnInit() {
    this.namespaceService.current.subscribe((namespace: string) => {
      this.currentNamespace = namespace;
    });

    this.namespaceService.list.subscribe((namespaces: string[]) => {
      this.namespaces = namespaces;
    });
  }

  selectNamespace(namespace: string) {
    this.namespaceService.setNamespace(namespace);
  }
}
