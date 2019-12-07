import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SingleStatComponent } from './single-stat.component';
import { SingleStatView } from '../../../../models/content';

describe('SingleStatComponent', () => {
  let component: SingleStatComponent;
  let fixture: ComponentFixture<SingleStatComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [SingleStatComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SingleStatComponent);
    component = fixture.componentInstance;
    const view: SingleStatView = {
      metadata: {
        type: 'singleStat',
      },
      config: {
        title: 'title',
        value: {
          color: 'red',
          text: 'text',
        },
      },
    };
    component.view = view;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
