module.exports = {
  stories: ['../src/**/*.@(stories|story).@(js|ts|mdx)'],
  addons: [
    '@storybook/addon-actions',
    '@storybook/addon-links',
    '@storybook/addon-knobs/preset',
    {
      name: '@storybook/addon-docs/preset',
      options: {
        configureJSX: true,
      },
    },
  ],
  typescript: {
    check: false,
    checkOptions: {},
    // also valid 'react-docgen-typescript' | false
    reactDocgen: 'react-docgen',
    reactDocgenTypescriptOptions: {
      shouldExtractLiteralValuesFromEnum: true,
      propFilter: (prop) => (prop.parent ? !/node_modules/.test(prop.parent.fileName) : true),
    },
  },
  // logLevel: 'debug',
};
