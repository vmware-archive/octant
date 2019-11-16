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
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements OnInit, OnDestroy, AfterViewInit {
  private terminalStream: TerminalOutputStreamer;
  selectedTerminal: TerminalDetail;
  @Input() view: TerminalView;
  @ViewChild('terminal', { static: true }) child: NgTerminal;
  trackByIdentity = trackByIdentity;

  constructor(
    private terminalService: TerminalOutputService,
    private wss: WebsocketService
  ) {}

  compareFn(c1: TerminalDetail, c2: TerminalDetail): boolean {
    return c1 && c2 ? c1.uuid === c2.uuid : c1 === c2;
  }

  ngOnInit() {
    this.initSelectedTerminal();
  }

  ngOnDestroy(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream = null;
    }
  }

  ngAfterViewInit() {
    if (this.selectedTerminal && this.view) {
      this.initSize();
      this.initStream();
    }
    this.enableResize();

    this.child.keyEventInput.subscribe(e => {
      const ev = e.domEvent;
      const printable = !ev.altKey && !ev.ctrlKey && !ev.metaKey;

      this.wss.sendMessage('sendTerminalCommand', {
        terminalID: this.selectedTerminal.uuid,
        key: e.key,
      });

      if (ev.keyCode === 8) {
        // Do not delete the prompt
        if (this.child.underlying.buffer.cursorX > 78) {
          this.child.write('\b \b');
        }
      }
    });
  }

  onTerminalChange(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream = null;
    }
    this.child.underlying.clear();
    this.initStream();
  }

  initSelectedTerminal() {
    if (this.view) {
      if (this.view.config.terminals && this.view.config.terminals.length > 0) {
        this.selectedTerminal = this.view.config.terminals[0];
      }
    }
  }

  enableResize() {
    let timeOut = null;
    const resizeDebounce = (e: { cols: number; rows: number }) => {
      const resize = () => {
        this.wss.sendMessage('sendTerminalResize', {
          terminalID: this.selectedTerminal.uuid,
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
      terminalID: this.selectedTerminal.uuid,
      rows: this.child.underlying.rows,
      cols: this.child.underlying.cols,
    });
  }

  initStream() {
    const namespace = this.view.config.namespace;
    const name = this.view.config.name;

    if (
      namespace &&
      name &&
      this.selectedTerminal.container &&
      this.selectedTerminal.uuid
    ) {
      this.terminalStream = this.terminalService.createStream(
        namespace,
        name,
        this.selectedTerminal.container,
        this.selectedTerminal.uuid
      );
      this.terminalStream.scrollback.subscribe((scrollback: string) => {
        if (scrollback && scrollback.length !== 0) {
          this.child.write(atob(scrollback).replace(/\n/g, '\n\r'));
        }
      });
      this.terminalStream.line.subscribe((line: string) => {
        if (line && line.length !== 0) {
          this.child.write(atob(line).replace(/\n/g, '\n\r'));
        }
      });
    }
  }
}
