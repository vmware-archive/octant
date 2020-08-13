import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HomeComponent } from './home.component';
import { initServiceStub } from 'src/app/testing/init-service-stub';
import { InitService } from 'src/app/modules/shared/services/init/init.service';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let initService: InitService;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [
        {
          provide: InitService,
          useValue: initServiceStub,
        },
      ],
      declarations: [HomeComponent],
    }).compileComponents();
  }));

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
