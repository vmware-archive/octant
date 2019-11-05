// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, ViewChild, AfterViewInit } from '@angular/core';
import { NgTerminal } from 'ng-terminal';

@Component({
  selector: 'app-terminal',
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements AfterViewInit {
  @ViewChild('terminal', { static: true }) child: NgTerminal;

  ngAfterViewInit() {
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
}
