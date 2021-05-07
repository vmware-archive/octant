import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AlertComponent } from './alert.component';
import { ButtonGroupComponent } from '../button-group/button-group.component';
import { ButtonComponent } from '../button/button.component';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { OctantTooltipComponent } from '../octant-tooltip/octant-tooltip';
import { windowProvider, WindowToken } from '../../../../../window';

describe('AlertComponent', () => {
  let component: AlertComponent;
  let fixture: ComponentFixture<AlertComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [
          AlertComponent,
          OctantTooltipComponent,
          ButtonGroupComponent,
          ButtonComponent,
        ],
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(AlertComponent);
    component = fixture.componentInstance;

    component.alert = {
      message: 'message',
      type: 'default',
      status: 'error',
      closable: true,
      buttonGroup: {
        metadata: { type: 'buttonGroup' },
        config: {
          buttons: [
            {
              metadata: { type: 'button' },
              config: {
                payload: {},
                name: 'Alert Button',
              },
            },
          ],
        },
      },
    };

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('sets the alert type', () => {
    expect(component.status).toBe('danger');
  });

  it('sets the message', () => {
    const el: DebugElement = fixture.debugElement.query(By.css('cds-alert'));
    expect(el.nativeElement.textContent.trim()).toBe('message');
  });

  it('sets button group', () => {
    const el: DebugElement = fixture.debugElement.query(
      By.css('app-button-group app-button')
    );
    expect(el.nativeElement.textContent.trim()).toContain('Alert Button');
  });
});
