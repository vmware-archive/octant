// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

// Karma configuration file, see link for more information
// https://karma-runner.github.io/1.0/config/configuration-file.html

module.exports = function (config) {
  config.set({
    basePath: '',
    frameworks: ['jasmine', '@angular-devkit/build-angular'],
    proxies: {
      '/vs/': 'vs/lib/',
    },
    plugins: [
      require('karma-jasmine'),
      require('karma-chrome-launcher'),
      require('karma-coverage'),
      require('karma-electron'),
      require('karma-electron-launcher'),
      require('karma-jasmine-html-reporter'),
      require('karma-spec-reporter'),
      require('@angular-devkit/build-angular/plugins/karma'),
    ],
    client: {
      clearContext: false, // leave Jasmine Spec Runner output visible in browser
    },
    coverageIstanbulReporter: {
      dir: require('path').join(__dirname, '../coverage/octant'),
      reports: ['html', 'lcovonly', 'text-summary'],
      fixWebpackSourcePaths: true,
    },
    reporters: ['spec', 'progress', 'kjhtml'],
    specReporter: {
      maxLogLines: 5,
      suppressErrorSummary: true,
      suppressFailed: false,
      suppressPassed: false,
      suppressedSkipped: false,
      showSpectiming: true,
      failFast: false,
    },
    port: 9876,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: true,
    browsers: ['Chrome', 'Electron'],
    singleRun: false,
    restartOnFileChange: true,
    customLaunchers: {
      ChromeDebug: {
        base: 'Chrome',
        flags: [
          '--remote-debugging-port=9333',
          '--no-sandbox',
          '--disable-extensions'
        ],
      },
      ElectronDebug: {
        base: 'Electron',
        browserWindowOptions: {
          webPreferences: {
            nodeIntegration: true,
            allowRunningInsecureContent: true
          }
        }
      }
    },
  });
};
