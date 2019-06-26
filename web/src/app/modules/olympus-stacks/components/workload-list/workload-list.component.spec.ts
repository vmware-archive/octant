import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { WorkloadCardComponent } from '../workload-card/workload-card.component';
import { WorkloadListComponent } from './workload-list.component';

describe('WorkloadListComponent', () => {
  let component: WorkloadListComponent;
  let fixture: ComponentFixture<WorkloadListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WorkloadListComponent, WorkloadCardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WorkloadListComponent);
    component = fixture.componentInstance;
    component.workloadList = {
      currentStack: 'olympusstackA',
      stackOptions: ['olympusstackA', 'releaseChannelB', 'olympusstackC'],
      channelFollowing: 'unstable',
      workloads: [],
    };
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
