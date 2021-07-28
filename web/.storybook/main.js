module.exports = {
  core: {
    builder: 'webpack5',
  },
  stories: ['../src/**/*.@(stories|story).@(js|ts|mdx)'],
  addons: [
    '@storybook/addon-a11y',
    '@storybook/addon-actions',
    '@storybook/addon-links',
    {
      name: '@storybook/addon-docs/preset',
      options: {
        configureJSX: true,
      },
    },
    '@storybook/addon-controls',
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
