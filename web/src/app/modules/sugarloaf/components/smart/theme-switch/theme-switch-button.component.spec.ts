// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ThemeSwitchButtonComponent } from './theme-switch-button.component';
import { ThemeService } from './theme-switch.service';
import { themeServiceStub } from 'src/app/testing/theme-service-stub';
import { By } from '@angular/platform-browser';
import { MonacoEditorConfig, MonacoProviderService } from 'ng-monaco-editor';

describe('ThemeSwitchButtonComponent', () => {
  let component: ThemeSwitchButtonComponent;
  let fixture: ComponentFixture<ThemeSwitchButtonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ThemeSwitchButtonComponent],
      providers: [
        { provide: ThemeService, useValue: themeServiceStub },
        MonacoEditorConfig,
        MonacoProviderService,
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ThemeSwitchButtonComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should switch theme', () => {
    component.switchTheme();

    expect(localStorage.getItem('theme')).toBe('dark');
    expect(component.lightThemeEnabled).toBe(false);

    component.switchTheme();

    expect(localStorage.getItem('theme')).toBe('light');
    expect(component.lightThemeEnabled).toBe(true);
  });

  it('should render the right button', () => {
    component.lightThemeEnabled = true;
    fixture.detectChanges();
    const switchButton = fixture.debugElement.query(By.css('#switchButton'))
      .nativeElement;

    expect(switchButton.innerHTML).toContain('dark');
  });
});
