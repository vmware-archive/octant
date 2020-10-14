import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AlertComponent } from './alert.component';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { OctantTooltipComponent } from '../octant-tooltip/octant-tooltip';

describe('AlertComponent', () => {
  let component: AlertComponent;
  let fixture: ComponentFixture<AlertComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [AlertComponent, OctantTooltipComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AlertComponent);
    component = fixture.componentInstance;

    component.alert = {
      message: 'message',
      type: 'warning',
    };

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('sets the alert style', () => {
    const el: DebugElement = fixture.debugElement.query(By.css('div'));
    expect(el.classes.hasOwnProperty('alert-warning')).toBeTruthy();
  });

  it('sets the proper icon', () => {
    const el: DebugElement = fixture.debugElement.query(By.css('.alert-icon'));
    expect(el.attributes.shape).toEqual('exclamation-triangle');
  });

  it('sets the message', () => {
    const el: DebugElement = fixture.debugElement.query(By.css('.alert-text'));
    expect(el.nativeElement.textContent.trim()).toBe('message');
  });
});
