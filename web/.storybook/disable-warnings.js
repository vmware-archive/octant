function DisableWarnings(options) {
}

DisableWarnings.prototype.apply = function(compiler) {
  compiler.plugin("emit", function(compilation, cb) {
    compilation.warnings = [];
    cb();
  });
};

module.exports = DisableWarnings;
