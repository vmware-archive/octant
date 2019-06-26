import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WorkloadCardComponent } from './workload-card.component';

describe('WorkloadCardComponent', () => {
  let component: WorkloadCardComponent;
  let fixture: ComponentFixture<WorkloadCardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WorkloadCardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WorkloadCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
