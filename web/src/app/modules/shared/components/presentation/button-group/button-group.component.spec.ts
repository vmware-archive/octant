import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ButtonGroupComponent } from './button-group.component';
import { ButtonGroupView } from '../../../models/content';

describe('ButtonGroupComponent', () => {
  let component: ButtonGroupComponent;
  let fixture: ComponentFixture<ButtonGroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ButtonGroupComponent],
    }).compileComponents();
  }));

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
});
