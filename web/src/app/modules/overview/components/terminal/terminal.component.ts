// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, ViewChild, AfterViewInit, Input } from '@angular/core';
import { NgTerminal } from 'ng-terminal';
import {
  TerminalOutputStreamer,
  TerminalOutputService,
} from 'src/app/services/terminals/terminals.service';
import { BehaviorSubject } from 'rxjs';
import { TerminalView } from 'src/app/models/content';

@Component({
  selector: 'app-terminal',
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements AfterViewInit {
  private terminalStream: TerminalOutputStreamer;
  @Input() view: TerminalView;
  @ViewChild('terminal', { static: true }) child: NgTerminal;
  scrollback: BehaviorSubject<string>;

  constructor(private terminalService: TerminalOutputService) {}

  ngAfterViewInit() {
    this.initStream();

    this.child.keyEventInput.subscribe(e => {
      const ev = e.domEvent;
      const printable = !ev.altKey && !ev.ctrlKey && !ev.metaKey;

      if (ev.keyCode === 13) {
        this.child.write('\r\n$ ');
      } else if (ev.keyCode === 8) {
        // Do not delete the prompt
        if (this.child.underlying.buffer.cursorX > 2) {
          this.child.write('\b \b');
        }
      } else if (printable) {
        this.child.write(e.key);
      }
    });
  }

  initStream() {
    const namespace = this.view.config.namespace;
    const name = this.view.config.name;
    const container = this.view.config.container;
    const uuid = this.view.config.uuid;

    if (namespace && name && container && uuid) {
      this.terminalStream = this.terminalService.createStream(
        namespace,
        name,
        container,
        uuid
      );
      this.terminalStream.scrollback.subscribe((scrollback: string) => {
        this.child.write(scrollback);
      });
      this.terminalStream.line.subscribe((line: string) => {
        this.child.write(line);
      });
    }
  }
}
