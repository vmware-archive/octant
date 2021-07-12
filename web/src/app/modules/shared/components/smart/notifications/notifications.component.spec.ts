import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { NotificationsComponent } from './notifications.component';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import { instance, mock } from 'ts-mockito';

describe('NotificationsComponent', () => {
  let component: NotificationsComponent;
  let fixture: ComponentFixture<NotificationsComponent>;
  const mockWebsocketService: WebsocketService = mock(WebsocketService);

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [NotificationsComponent],
        providers: [
          {
            provide: WebsocketService,
            useValue: instance(mockWebsocketService),
          },
        ],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(NotificationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
