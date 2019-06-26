import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OlympusStacksComponent } from './olympus-stacks.component';

describe('OlympusStacksComponent', () => {
  let component: OlympusStacksComponent;
  let fixture: ComponentFixture<OlympusStacksComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OlympusStacksComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OlympusStacksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
