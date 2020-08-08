import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HomeComponent } from './home.component';
import {
  MonacoEditorComponent,
  MonacoProviderService,
  MonacoEditorConfig,
} from 'ng-monaco-editor';
import { themeServiceStub } from 'src/app/testing/theme-service-stub';
import { ThemeService } from 'src/app/modules/sugarloaf/components/smart/theme-switch/theme-switch.service';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let themeService: ThemeService;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [
        {
          provide: ThemeService,
          useValue: themeServiceStub,
        },
        MonacoProviderService,
        MonacoEditorComponent,
        MonacoEditorConfig,
      ],
      declarations: [HomeComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    themeService = TestBed.inject(ThemeService);

    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load light theme when component mounts', () => {
    spyOn(themeService, 'loadTheme');
    component.ngOnInit();

    expect(themeService.isLightThemeEnabled()).toBeTruthy();
    expect(themeService.loadTheme).toHaveBeenCalled();
  });
});
