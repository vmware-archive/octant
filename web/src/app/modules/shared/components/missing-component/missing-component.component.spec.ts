import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MissingComponentComponent } from './missing-component.component';
import { CommonModule } from '@angular/common';

describe('MissingComponentComponent', () => {
  let component: MissingComponentComponent;
  let fixture: ComponentFixture<MissingComponentComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [
          MissingComponentComponent,
          ApplyYAMLComponent,
          OverlayScrollbarsComponent,
        ],
        imports: [CommonModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(MissingComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
