// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ThemeSwitchButtonComponent } from './theme-switch-button.component';
import { ThemeService } from '../../../../shared/services/theme/theme.service';
import { themeServiceStub } from 'src/app/testing/theme-service-stub';
import { By } from '@angular/platform-browser';

describe('ThemeSwitchButtonComponent', () => {
  let component: ThemeSwitchButtonComponent;
  let fixture: ComponentFixture<ThemeSwitchButtonComponent>;
  let themeService: ThemeService;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ThemeSwitchButtonComponent],
        providers: [{ provide: ThemeService, useValue: themeServiceStub }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    themeService = TestBed.inject(ThemeService);

    fixture = TestBed.createComponent(ThemeSwitchButtonComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should switch theme', () => {
    spyOn(themeService, 'switchTheme');

    component.switchTheme();

    expect(themeService.switchTheme).toHaveBeenCalled();
  });

  it('should render the right button', () => {
    component.lightThemeEnabled = true;
    fixture.detectChanges();
    const switchButton = fixture.debugElement.query(
      By.css('#switchButton')
    ).nativeElement;

    expect(switchButton.innerHTML).toContain('dark');
  });
});
