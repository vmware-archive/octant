import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { OverviewModule } from '../../overview.module';
import { ResourceViewerComponent } from './resource-viewer.component';

describe('ResourceViewerComponent', () => {
  let component: ResourceViewerComponent;
  let fixture: ComponentFixture<ResourceViewerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ OverviewModule ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
