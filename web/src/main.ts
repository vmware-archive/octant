// Note(marlon): zone.js seems to add a thread blocking event listener
// to mouse scroll events: https://github.com/angular/zone.js/issues/1097
// This makes all event listeners passive by default. To note when attaching
// custom event listeners in our app.
import 'default-passive-events';
import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AppModule } from './app/app.module';
import { environment } from './environments/environment';

if (environment.production) {
  enableProdMode();
}

platformBrowserDynamic()
  .bootstrapModule(AppModule)
  .catch((err) => console.error(err));
