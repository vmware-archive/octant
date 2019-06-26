import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { WorkloadListComponent } from '../workload-list/workload-list.component';
import { WorkloadCardComponent } from '../workload-card/workload-card.component';
import { OlympusStacksComponent } from './olympus-stacks.component';

describe('OlympusStacksComponent', () => {
  let component: OlympusStacksComponent;
  let fixture: ComponentFixture<OlympusStacksComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OlympusStacksComponent, WorkloadListComponent, WorkloadCardComponent ]
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
