import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HomeComponent } from './home.component';
import { initServiceStub } from 'src/app/testing/init-service-stub';
import { InitService } from 'src/app/modules/shared/services/init/init.service';
import { OverlayScrollbarsComponent } from 'overlayscrollbars-ngx';
import { ApplyYAMLComponent } from 'src/app/modules/sugarloaf/components/smart/apply-yaml/apply-yaml.component';
import { EditorComponent } from 'src/app/modules/shared/components/smart/editor/editor.component';
import { SharedModule } from 'src/app/modules/shared/shared.module';

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
        SharedModule,
      ],
      declarations: [
        HomeComponent,
        ApplyYAMLComponent,
        EditorComponent,
        OverlayScrollbarsComponent,
      ],
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
