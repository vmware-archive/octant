import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MissingComponentComponent } from './missing-component.component';
import { BrowserModule } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';
import { ApplyYAMLComponent } from 'src/app/modules/sugarloaf/components/smart/apply-yaml/apply-yaml.component';
import { OverlayScrollbarsComponent } from 'overlayscrollbars-ngx';

describe('MissingComponentComponent', () => {
  let component: MissingComponentComponent;
  let fixture: ComponentFixture<MissingComponentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        MissingComponentComponent,
        ApplyYAMLComponent,
        OverlayScrollbarsComponent,
      ],
      imports: [CommonModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MissingComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
