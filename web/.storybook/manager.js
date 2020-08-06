import { addons } from '@storybook/addons';
import { themes } from '@storybook/theming';

addons.setConfig({
  theme: themes.light,
  previewTabs: {
    'storybook/docs/panel': { index: -1 },
  },
});
