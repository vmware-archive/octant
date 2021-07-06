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
        this.newContent = true;
        this.upperCaseFirst(ie);
        this.notifications = ie.errors;
      }
    );
    this.websocketService.sendMessage('event.octant.dev/notification', {});
  }

  toggleModal(): void {
    this.newContent = false;
    this.hidden = !this.hidden;
  }

  upperCaseFirst(ie: { errors: InternalError[] }) {
    ie.errors.map(e => {
      e.error = e.error[0].toUpperCase() + e.error.substr(1);
      return e;
    });
  }
}
