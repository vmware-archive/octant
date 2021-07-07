import { Component, OnInit } from '@angular/core';
import { ClarityIcons, bellIcon } from '@cds/core/icon';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import { InternalError } from '../../../models/content';

@Component({
  selector: 'app-notifications',
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.scss'],
})
export class NotificationsComponent implements OnInit {
  hidden = true;
  notifications: InternalError[] = [];
  newContent = false;

  trackByFn = trackByIndex;

  constructor(private websocketService: WebsocketService) {
    ClarityIcons.addIcons(bellIcon);
  }

  ngOnInit(): void {
    this.websocketService.registerHandler(
      'event.octant.dev/notification',
      (ie: { errors: InternalError[] }) => {
        if (ie.errors == null) {
          return;
        }
        this.newContent = true;
        this.upperCaseFirst(ie);
        this.notifications = ie.errors;
      }
    );
    // Ask for the errors on the server, it only happends the first time
    this.websocketService.sendMessage('event.octant.dev/notification', {});
  }

  toggleModal(): void {
    this.newContent = false;
    this.hidden = !this.hidden;
  }

  upperCaseFirst(ie: { errors: InternalError[] }) {
    ie.errors.map(e => {
      if (e.error) {
        e.error = e.error[0].toUpperCase() + e.error.substr(1);
      }
      return e;
    });
  }
}
