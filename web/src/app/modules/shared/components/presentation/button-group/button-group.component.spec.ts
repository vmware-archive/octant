import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ButtonGroupComponent } from './button-group.component';
import { ButtonGroupView } from '../../../models/content';
import { windowProvider, WindowToken } from '../../../../../window';

describe('ButtonGroupComponent', () => {
  let component: ButtonGroupComponent;
  let fixture: ComponentFixture<ButtonGroupComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ButtonGroupComponent],
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ButtonGroupComponent);
    component = fixture.componentInstance;

    component.view = {
      config: {
        buttons: [],
      },
    } as ButtonGroupView;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('creates button class', () => {
    const cases = [
      {
        style: 'outline',
        size: 'block',
        status: 'info',
        expected: 'btn btn-block btn-info-outline',
      },
      {
        style: 'flat',
        status: 'disabled',
        expected: 'btn btn-sm disabled btn-flat',
      },
      {
        style: 'solid',
        size: 'lg',
        expected: 'btn btn-solid',
      },
      {
        expected: 'btn btn-sm btn-outline',
      },
    ];
    cases.forEach(test => {
      const result = component.buttonClass(test.style, test.size, test.status);
      expect(result).toEqual(test.expected);
    });
  });
});
