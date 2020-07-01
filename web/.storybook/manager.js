import { addons } from '@storybook/addons';
import { themes } from '@storybook/theming';
import { STORY_CHANGED, SET_CURRENT_STORY } from '@storybook/core-events';

addons.setConfig({
  theme: themes.light,
  previewTabs: {
    'storybook/docs/panel': { index: -1 },
  },
});

let firstTime = true;
let notUpdated = true;

// Show doc panel initially
addons.register('Octant/SetDocsStory', api => {
  api.on(SET_CURRENT_STORY, storyId => {
    if (firstTime && storyId.viewMode !== 'docs') {
      storyId.viewMode = 'docs';
      if(!notUpdated) {
        api.emit(SET_CURRENT_STORY, {
          ...storyId,
          viewMode: 'docs'
        });
      }
      notUpdated= false;
    }
  });

  api.on(STORY_CHANGED, () => {
    firstTime= false;
  });
});
