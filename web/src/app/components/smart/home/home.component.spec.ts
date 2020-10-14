import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { HomeComponent } from './home.component';
import { initServiceStub } from 'src/app/testing/init-service-stub';
import { InitService } from 'src/app/modules/shared/services/init/init.service';
import { SharedModule } from 'src/app/modules/shared/shared.module';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let initService: InitService;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        providers: [
          {
            provide: InitService,
            useValue: initServiceStub,
          },
          SharedModule,
        ],
        declarations: [HomeComponent],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    initService = TestBed.inject(InitService);

    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize when component mounts', () => {
    spyOn(initService, 'init');
    component.ngOnInit();

    expect(initService.init).toHaveBeenCalled();
  });
});
