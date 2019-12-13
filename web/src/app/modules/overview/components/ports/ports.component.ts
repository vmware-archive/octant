// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  Output,
  SimpleChanges,
} from '@angular/core';
import _ from 'lodash';
import { Port, PortsView, View } from 'src/app/models/content';
import {
  NotifierService,
  NotifierSession,
} from 'src/app/services/notifier/notifier.service';
import { PortForwardService } from 'src/app/services/port-forward/port-forward.service';

@Component({
  selector: 'app-view-ports',
  templateUrl: './ports.component.html',
  styleUrls: ['./ports.component.scss'],
})
export class PortsComponent implements OnChanges, OnDestroy {
  private notifierSession: NotifierSession;
  private submittedPFCreation: string;
  private submittedPFRemoval: string;

  v: PortsView;

  @Input() set view(v: View) {
    this.v = v as PortsView;
  }
  get view() {
    return this.v;
  }

  @Output() portLoad: EventEmitter<boolean> = new EventEmitter(true);

  constructor(
    private portForwardService: PortForwardService,
    notifierService: NotifierService
  ) {
    this.notifierSession = notifierService.createSession();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as PortsView;

      if (this.submittedPFCreation) {
        const foundPort = _.find(view.config.ports, (port: Port) => {
          return this.submittedPFCreation === port.config.name;
        }) as Port;
        if (foundPort && foundPort.config.state.isForwarded) {
          this.portLoad.emit(false);
          this.submittedPFCreation = '';
        }
      } else if (this.submittedPFRemoval) {
        const foundPort = _.find(view.config.ports, (port: Port) => {
          return this.submittedPFRemoval === port.config.name;
        }) as Port;
        if (foundPort && !foundPort.config.state.isForwarded) {
          this.portLoad.emit(false);
          this.submittedPFCreation = '';
        }
      }
    }
  }

  identifyPort(index: number, item: Port) {
    return item.config.name;
  }

  startPortForward(port: Port) {
    this.portLoad.emit(true);
    this.submittedPFCreation = port.config.name;

    this.portForwardService.create(port);
  }

  removePortForward(port: Port) {
    this.portLoad.emit(true);
    this.submittedPFRemoval = port.config.name;

    this.portForwardService.remove(port.config.state.id);
  }

  openPortForward(port: Port) {
    if (!port.config.state.isForwarded) {
      return;
    }
    const localhostUrl = `http://localhost:${port.config.state.port}`;
    window.open(localhostUrl, '_blank');
  }

  ngOnDestroy() {
    this.notifierSession.removeAllSignals();
  }
}
