import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { HomeComponent } from './home.component';
import { themeServiceStub } from 'src/app/testing/theme-service-stub';
import { ThemeService } from 'src/app/modules/shared/services/theme/theme.service';
import { SharedModule } from 'src/app/modules/shared/shared.module';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let themeService: ThemeService;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        providers: [
          {
            provide: ThemeService,
            useValue: themeServiceStub,
          },
          SharedModule,
        ],
        declarations: [HomeComponent],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    themeService = TestBed.inject(ThemeService);

    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize when component mounts', () => {
    spyOn(themeService, 'loadTheme');
    component.ngOnInit();

    expect(themeService.loadTheme).toHaveBeenCalled();
  });
});
