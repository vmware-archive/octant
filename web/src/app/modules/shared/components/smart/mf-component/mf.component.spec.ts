import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MfComponent } from './mf.component';
import { MfComponentView } from '../../../models/content';
import { windowProvider, WindowToken } from '../../../../../window';

describe('MfComponent', () => {
  let component: MfComponent;
  let fixture: ComponentFixture<MfComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [MfComponent],
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(MfComponent);
    component = fixture.componentInstance;

    component.view = {
      config: {},
    } as MfComponentView;

    fixture.detectChanges();
  });

  fit('should create', () => {
    expect(component).toBeTruthy();
  });
});
