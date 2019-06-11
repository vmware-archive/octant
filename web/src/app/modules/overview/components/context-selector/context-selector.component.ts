import { Component, OnInit } from '@angular/core';
import {
  ContextDescription,
  KubeContextService,
} from '../../services/kube-context/kube-context.service';

@Component({
  selector: 'app-context-selector',
  templateUrl: './context-selector.component.html',
  styleUrls: ['./context-selector.component.scss'],
})
export class ContextSelectorComponent implements OnInit {
  contexts: ContextDescription[];
  selected: string;

  constructor(private kubeContext: KubeContextService) {}

  ngOnInit() {
    this.kubeContext
      .contexts()
      .subscribe(contexts => (this.contexts = contexts));
    this.kubeContext
      .selected()
      .subscribe(selected => (this.selected = selected));
  }

  contextClass(context: ContextDescription) {
    const active = this.selected === context.name ? ['active'] : [];
    return ['context-button', ...active];
  }

  selectContext(context: ContextDescription) {
    this.kubeContext.select(context);
  }
}
