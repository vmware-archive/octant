/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Meta } from '@storybook/angular/types-6-0';

export default {
  title: 'Other/Bottom Panel',
} as Meta;

export const panelStory = () => {
  return {
    styles: [],
    template: `
    <div style="display: flex; flex-direction: column; height: 500px; background: hsl(198, 83%, 94%)">
        <div style="flex: 1 1 auto; background: hsl(198, 0%, 98%); overflow-y: scroll">
            <p *ngFor="let item of [].constructor(30); index as i">row {{i+1}}</p>
        </div>
        <app-bottom-panel>
            bottom content
        </app-bottom-panel>
    </div>
        `,
  };
};
panelStory.storyName = 'Panel';
