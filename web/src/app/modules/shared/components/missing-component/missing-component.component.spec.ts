import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MissingComponentComponent } from './missing-component.component';
import { BrowserModule } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';

describe('MissingComponentComponent', () => {
  let component: MissingComponentComponent;
  let fixture: ComponentFixture<MissingComponentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [MissingComponentComponent],
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
