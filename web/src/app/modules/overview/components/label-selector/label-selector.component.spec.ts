import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LabelSelectorComponent } from './label-selector.component';

describe('LabelSelectorComponent', () => {
  let component: LabelSelectorComponent;
  let fixture: ComponentFixture<LabelSelectorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LabelSelectorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LabelSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
