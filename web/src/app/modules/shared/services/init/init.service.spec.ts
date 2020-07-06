import { inject, TestBed } from '@angular/core/testing';
import { InitService } from './init.service';
import { ThemeService } from '../../../sugarloaf/components/smart/theme-switch/theme-switch.service';
import { MonacoEditorConfig, MonacoProviderService } from 'ng-monaco-editor';

describe('InitService', () => {
  let init: InitService;
  let service: ThemeService;
  let monaco: MonacoProviderService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        InitService,
        ThemeService,
        MonacoProviderService,
        MonacoEditorConfig,
        Document,
      ],
    });
    init = TestBed.inject(InitService);
    service = TestBed.inject(ThemeService);
    monaco = TestBed.inject(MonacoProviderService);
  });

  it('should be created', () => {
    expect(init).toBeTruthy();
    expect(service).toBeTruthy();
    expect(monaco).toBeTruthy();
  });
});
