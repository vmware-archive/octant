// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  Output,
  EventEmitter,
  OnDestroy,
} from '@angular/core';
import _ from 'lodash';
import { Port, PortsView } from 'src/app/models/content';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from 'src/app/services/notifier/notifier.service';
import { PortForwardService } from 'src/app/services/port-forward/port-forward.service';

@Component({
  selector: 'app-ports',
  templateUrl: './ports.component.html',
  styleUrls: ['./ports.component.scss'],
})
export class PortsComponent implements OnChanges, OnDestroy {
  private notifierSession: NotifierSession;
  private submittedPFCreation: string;
  private submittedPFRemoval: string;

  @Input() view: PortsView;
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

    this.portForwardService.create(port).subscribe(
      () => {
        // TODO: handle success
      },
      () => {
        this.notifierSession.pushSignal(
          NotifierSignalType.ERROR,
          'There was an issue starting your port-forward'
        );
        this.portLoad.emit(false);
      }
    );
  }

  removePortForward(port: Port) {
    this.portLoad.emit(true);
    this.submittedPFRemoval = port.config.name;

    this.portForwardService.remove(port).subscribe(
      () => {
        // TODO: handle success
      },
      () => {
        this.notifierSession.pushSignal(
          NotifierSignalType.ERROR,
          'There was an issue removing your port-forward'
        );
        this.portLoad.emit(false);
      }
    );
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
