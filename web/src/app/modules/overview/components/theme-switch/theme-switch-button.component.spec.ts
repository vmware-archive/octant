// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ThemeSwitchButtonComponent } from './theme-switch-button.component';
import { ThemeService } from './theme-switch.service';
import { themeServiceStub } from 'src/app/testing/theme-service-stub';
import { By } from '@angular/platform-browser';

describe('ThemeSwitchButtonComponent', () => {
  let component: ThemeSwitchButtonComponent;
  let fixture: ComponentFixture<ThemeSwitchButtonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ThemeSwitchButtonComponent],
      providers: [{ provide: ThemeService, useValue: themeServiceStub }]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ThemeSwitchButtonComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load light theme when component mounts if there is not a stored value', () => {
    spyOn(component, 'loadTheme');
    component.ngOnInit();

    expect(component.theme).toBe('light');
    expect(component.loadTheme).toHaveBeenCalled();
  });

  it('should load right theme when component mounts if there is a stored value', () => {
    spyOn(component, 'loadTheme');
    localStorage.setItem('theme', 'dark');
    component.ngOnInit();

    expect(component.theme).toBe('dark');
    expect(component.loadTheme).toHaveBeenCalled();
  });

  it('should indicate if the light theme is active or not', () => {
    component.theme = 'light';
    expect(component.isLightThemeEnabled()).toBe(true);

    component.theme = 'dark';
    expect(component.isLightThemeEnabled()).toBe(false);
  });

  it('should switch theme', () => {
    component.theme = 'light';
    component.switchTheme();

    expect(component.theme).toBe('dark');
    expect(localStorage.getItem('theme')).toBe('dark');

    component.switchTheme();

    expect(component.theme).toBe('light');
    expect(localStorage.getItem('theme')).toBe('light');
  });

  it('should render the right button', () => {
    component.theme = 'light';
    fixture.detectChanges();
    const switchButton = fixture.debugElement.query(By.css('#switchButton')).nativeElement;
    
    expect(switchButton.innerHTML).toContain('dark');
  });
});
