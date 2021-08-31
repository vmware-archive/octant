import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import {
  IndicatorComponent,
  Status,
  statusLookup,
  iconLookup,
} from './indicator.component';
import { CdsModule } from '@cds/angular';
import '@cds/core/icon/register.js';
import { Component } from '@angular/core';
import { View } from '../../../models/content';

@Component({
  template:
    '<app-indicator [status]="status" [detail]="detail"></app-indicator>',
})
class WrapperComponent {
  status: number;
  detail: View;
}

describe('IndicatorComponent', () => {
  let component: WrapperComponent;
  let fixture: ComponentFixture<WrapperComponent>;

  let element: HTMLDivElement;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [WrapperComponent, IndicatorComponent],
        imports: [CdsModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(WrapperComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  [Status.Ok, Status.Warning, Status.Error].forEach(v => {
    const name = statusLookup[v];

    describe(`with ${name} status`, () => {
      beforeEach(() => {
        element = fixture.nativeElement;
        component.status = v;
        component.detail = {
          metadata: {
            type: 'test',
          },
        };
        fixture.detectChanges();
      });

      it(`shows ${name} indicator`, () => {
        const e = element.querySelector(`app-indicator cds-icon`);
        expect(e.getAttribute('ng-reflect-shape')).toBe(iconLookup[v]);
      });
    });
  });

  describe('with unknown status', () => {
    beforeEach(() => {
      element = fixture.nativeElement;
      component.status = 0;
      fixture.detectChanges();
    });

    it('does not show an indicator', () => {
      expect(element.querySelector('app-indicator cds-icon')).toBeNull();
    });
  });
});
