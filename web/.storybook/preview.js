import { addDecorator, moduleMetadata } from '@storybook/angular';
import { SharedModule } from '../src/app/modules/shared/shared.module';
import { setConsoleOptions } from '@storybook/addon-console';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { APP_INITIALIZER } from '@angular/core';
import { ThemeService } from '../src/app/modules/shared/services/theme/theme.service';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from '../src/app/app-routing.module';
import { setCompodocJson } from '@storybook/addon-docs/angular';
import docJson from '../documentation.json';
import { MonacoEditorModule } from '@materia-ui/ngx-monaco-editor';
import { windowProvider, WindowToken } from '../src/app/window';

import theme from './theme';

setCompodocJson(docJson);

setConsoleOptions({
  panelExclude: [
    /Angular is running in the development mode/,
    /Ignored an update to unaccepted module/,
  ],
});

addDecorator(
  moduleMetadata({
    imports: [
      AppRoutingModule,
      BrowserAnimationsModule,
      BrowserModule,
      HttpClientModule,
      MonacoEditorModule,
      RouterModule,
      SharedModule,
    ],
    providers: [
      ThemeService,
      { provide: WindowToken, useFactory: windowProvider },
      {
        provide: APP_INITIALIZER,
        useFactory: themeService => () => themeService.loadTheme(),
        deps: [ThemeService],
        multi: true,
      },
    ],
  })
);

export const parameters = {
  docs: {
    theme: theme,
  },
  options: {
    storySort: (a, b) => {
      // Show component stories on top
      let leftId = a[1].id;
      let rightId = b[1].id;

      if (leftId.startsWith('docs')) {
        leftId = '1' + leftId;
        if (leftId.includes('intro')) {
          leftId = '0' + leftId;
        }
      }

      if (rightId.startsWith('docs')) {
        rightId = '1' + rightId;
        if (rightId.includes('intro')) {
          rightId = '0' + rightId;
        }
      }

      return leftId.localeCompare(rightId, {
        numeric: true,
        ignorePunctuation: true,
      });
    },
  },
};
