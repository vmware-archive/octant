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
    component.workload = {
      name: 'ui',
      isPinned: true,
      lastUpdated: new Date(),
      revision: '2f2932c83b7a401d960f4538bf787e12c44dfd666',
      sourceImage: 'sourceImageA',
    };
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
