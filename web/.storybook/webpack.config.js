const { merge } = require('webpack-merge');

const maxAssetSize = 1024 * 1024;

module.exports = async ({ config }) => {
  return merge(config, {
    optimization: {
      splitChunks: {
        chunks: 'all',
        minSize: 30 * 1024,
        maxSize: maxAssetSize,
      }
    },
    performance: {
      maxAssetSize: maxAssetSize
    }
  });
};
