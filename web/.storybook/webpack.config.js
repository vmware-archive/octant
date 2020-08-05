const DisableWarnings = require('./disable-warnings.js');

module.exports = async ({ config }) => {
  config.plugins.push(new DisableWarnings());
  return config;
};
