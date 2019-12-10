// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  ViewChild,
  OnInit,
  OnDestroy,
  AfterViewInit,
  Input,
} from '@angular/core';
import { NgTerminal } from 'ng-terminal';
import {
  TerminalOutputStreamer,
  TerminalOutputService,
} from 'src/app/services/terminals/terminals.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { TerminalView, TerminalDetail } from 'src/app/models/content';
import { WebsocketService } from '../../services/websocket/websocket.service';

@Component({
  selector: 'app-terminal',
  styleUrls: ['./terminal.component.scss'],
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements OnDestroy, AfterViewInit {
  private terminalStream: TerminalOutputStreamer;
  @Input() view: TerminalView;
  @ViewChild('terminal', { static: true }) child: NgTerminal;
  trackByIdentity = trackByIdentity;

  constructor(
    private terminalService: TerminalOutputService,
    private wss: WebsocketService
  ) {
}

  compareFn(c1: TerminalDetail, c2: TerminalDetail): boolean {
    return c1 && c2 ? c1.uuid === c2.uuid : c1 === c2;
  }

  ngOnDestroy(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream = null;
    }
  }

  ngAfterViewInit() {
    if (this.view) {
      this.initSize();
      this.initStream();
    }
    this.enableResize();

    this.child.keyEventInput.subscribe(e => {
      this.wss.sendMessage('sendTerminalCommand', {
        terminalID: this.view.config.terminal.uuid,
        key: e.key,
      });
    });
  }

  onTerminalChange(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream = null;
    }
    this.child.underlying.clear();
    this.child.underlying.reset();
    this.initStream();
  }

  enableResize() {
    let timeOut = null;
    const resizeDebounce = (e: { cols: number; rows: number }) => {
      const resize = () => {
        this.wss.sendMessage('sendTerminalResize', {
          terminalID: this.view.config.terminal.uuid,
          rows: e.rows,
          cols: e.cols,
        });
      };

      if (timeOut != null) {
        clearTimeout(timeOut);
      }
      timeOut = setTimeout(resize, 1000);
    };
    this.child.underlying.onResize(resizeDebounce);
  }

  initSize() {
    this.wss.sendMessage('sendTerminalResize', {
      terminalID: this.view.config.terminal.uuid,
      rows: this.child.underlying.rows,
      cols: this.child.underlying.cols,
    });
  }

  initStream() {
    const namespace = this.view.config.namespace;
    const name = this.view.config.name;
    const terminal = this.view.config.terminal;

    if (namespace && name && terminal.container && terminal.uuid) {
      this.terminalStream = this.terminalService.createStream(
        namespace,
        name,
        terminal.container,
        terminal.uuid
      );
      this.terminalStream.scrollback.subscribe((scrollback: string) => {
        if (scrollback && scrollback.length !== 0) {
          this.child.write(atob(scrollback).replace(/\n/g, '\n\r'));
        }
      });
      if (terminal.active) {
        this.terminalStream.line.subscribe((line: string) => {
          if (line && line.length !== 0) {
            this.child.write(atob(line).replace(/\n/g, '\n\r'));
          }
        });
      }
    }
  }
}
