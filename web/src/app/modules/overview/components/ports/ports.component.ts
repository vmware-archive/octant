import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import _ from 'lodash';
import { Port, PortsView } from 'src/app/models/content';
import { NotifierService } from 'src/app/services/notifier/notifier.service';
import { PortForwardService } from 'src/app/services/port-forward/port-forward.service';

@Component({
  selector: 'app-ports',
  templateUrl: './ports.component.html',
  styleUrls: ['./ports.component.scss']
})
export class PortsComponent implements OnChanges {
  @Input() view: PortsView;
  submittedPFCreation: string;
  submittedPFRemoval: string;

  constructor(
    private portForwardService: PortForwardService,
    private notifierService: NotifierService,
  ) { }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as PortsView;

      if (this.submittedPFCreation) {
        const foundPort = _.find(view.config.ports, (port: Port) => {
          return this.submittedPFCreation === port.config.name;
        }) as Port;
        if (foundPort && foundPort.config.state.isForwarded) {
          this.notifierService.loading.next(false);
          this.submittedPFCreation = '';
        }
      } else if (this.submittedPFRemoval) {
        const foundPort = _.find(view.config.ports, (port: Port) => {
          return this.submittedPFRemoval === port.config.name;
        }) as Port;
        if (foundPort && !foundPort.config.state.isForwarded) {
          this.notifierService.loading.next(false);
          this.submittedPFCreation = '';
        }
      }
    }
  }

  identifyPort(index: number, item: Port) {
    return item.config.name;
  }

  startPortForward(port: Port) {
    this.notifierService.loading.next(true);
    this.submittedPFCreation = port.config.name;

    this.portForwardService.create(port).subscribe(() => {
      // TODO: handle success
    }, () => {
      this.notifierService.error.next('There was an issue starting your port-forward');
      this.notifierService.loading.next(false);
    });
  }

  removePortForward(port: Port) {
    this.notifierService.loading.next(true);
    this.submittedPFRemoval = port.config.name;

    this.portForwardService.remove(port).subscribe(() => {
      // TODO: handle success
    }, () => {
      this.notifierService.error.next('There was an issue removing your port-forward');
      this.notifierService.loading.next(false);
    });
  }

  openPortForward(port: Port) {
    if (!port.config.state.isForwarded) {
      return;
    }
    const localhostUrl = `http://localhost:${port.config.state.port}`;
    window.open(localhostUrl, '_blank');
  }
}
