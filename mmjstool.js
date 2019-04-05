#!/usr/bin/env node
// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
"use strict";

var yargs = require('yargs');

var i18nHandlers = require('./i18n_handlers');
/*eslint no-unused-vars: ["error", { "varsIgnorePattern": "[iI]gnored" }]*/


var ignored = yargs.usage('Usage: mmjstool <command> [options]').example('mmjstool i18n extract-webapp --webapp-dir ./', 'Extract all the i18n strings from the webapp source code').demandCommand(1).help('h').alias('h', 'help').command('i18n', 'I18n management commands', function (i18nArgs) {
  i18nArgs.demandCommand(1).command('extract-webapp', 'Read the source code, find all the translations string from mattermost-webapp and write them to the file mattermost-webapp/i18n/en.json', function () {
    /* empty function */
  }, i18nHandlers.i18nExtractWebapp).command('extract-mobile', 'Read the source code, find all the translations string from mattermost-mobile and write them to the file mattermost-mobile/assets/base/i18n/en.json.', function () {
    /* empty function */
  }, i18nHandlers.i18nExtractMobile).command('combine', 'Read the translations string from mattermost-webapp and mattermost-mobile and combine them in a single file', function (combineArgs) {
    combineArgs.demandCommand(2).option('output', {
      describe: 'File to store the combined translations',
      default: 'en.json'
    });
  }, i18nHandlers.i18nCombine).command('split', 'Read a set of combined translation files, and split them in mattermost-server and mattermost-web translations', function (splitArgs) {
    splitArgs.option('inputs', {
      describe: 'List of file to read the combined translations, splitted by ",". (e.g. en.json,es.json,fr.json)',
      default: 'en.json'
    });
  }, i18nHandlers.i18nSplit).command('sort', 'read a file and sort the content', function (sortArgs) {
    sortArgs.demandCommand(1).option('output', {
      describe: 'File to store sorted translations',
      default: 'en.json'
    });
  }, i18nHandlers.i18nSort).command('check', 'Read the source code, find all the translations string, and show you the differences with the current i18n/en.json files', function () {
    /* empty function */
  }, i18nHandlers.i18nCheck).command('check-mobile', 'Read the source code, find all the translations string, and show you the differences with the current i18n/en.json files', function () {
    /* empty function */
  }, i18nHandlers.i18nCheckMobile).command('check-webapp', 'Read the source code, find all the translations string, and show you the differences with the current i18n/en.json files', function () {
    /* empty function */
  }, i18nHandlers.i18nCheckWebapp).option('webapp-dir', {
    describe: 'webapp source code directory',
    default: '../mattermost-webapp'
  }).option('mobile-dir', {
    describe: 'mobile source code directory',
    default: '../mattermost-mobile'
  });
}, function () {
  /* empty function */
}).argv;
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.extractFromDirectory = extractFromDirectory;

function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _nonIterableSpread(); }

function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance"); }

function _iterableToArray(iter) { if (Symbol.iterator in Object(iter) || Object.prototype.toString.call(iter) === "[object Arguments]") return Array.from(iter); }

function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = new Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } }

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
var fs = require('fs');

var FileHound = require('filehound');

var Parser = require('flow-parser');

var walk = require('estree-walk');

var translatableComponents = {
  FormattedText: [{
    id: 'id',
    default: 'defaultMessage'
  }],
  FormattedMessage: [{
    id: 'id',
    default: 'defaultMessage'
  }],
  FormattedHTMLMessage: [{
    id: 'id',
    default: 'defaultMessage'
  }],
  FormattedMarkdownMessage: [{
    id: 'id',
    default: 'defaultMessage'
  }],
  FormattedMarkdownText: [{
    id: 'id',
    default: 'defaultMessage'
  }],
  FormattedAdminHeader: [{
    id: 'id',
    default: 'defaultMessage'
  }],
  LocalizedInput: ['placeholder']
};

function extractFromDirectory(dirPaths) {
  var filters = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : [];
  return new Promise(function (resolve) {
    var promises = dirPaths.map(function (dirPath) {
      return new Promise(function (innerResolve) {
        var translations = {};
        FileHound.create().paths(dirPath).discard(filters).ext('js', 'jsx').find().then(function (files) {
          var _iteratorNormalCompletion = true;
          var _didIteratorError = false;
          var _iteratorError = undefined;

          try {
            for (var _iterator = files[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
              var file = _step.value;

              try {
                Object.assign(translations, extractFromFile(file));
              } catch (e) {
                console.log(e);
                console.log('Unable to parse file:', file);
                console.log('Error in: line', e.loc && e.loc.line, 'column', e.loc && e.loc.column);
                return;
              }
            }
          } catch (err) {
            _didIteratorError = true;
            _iteratorError = err;
          } finally {
            try {
              if (!_iteratorNormalCompletion && _iterator.return != null) {
                _iterator.return();
              }
            } finally {
              if (_didIteratorError) {
                throw _iteratorError;
              }
            }
          }

          innerResolve(translations);
        });
      });
    });
    Promise.all(promises).then(function (translations) {
      resolve(Object.assign.apply(Object, [{}].concat(_toConsumableArray(translations))));
    });
  });
}

function extractFromFile(path) {
  var translations = {};
  var code = fs.readFileSync(path, 'utf-8');
  var ast = Parser.parse(code, {
    esproposal_class_static_fields: true,
    esproposal_class_instance_fields: true,
    esproposal_optional_chaining: true
  });
  walk(ast, {
    CallExpression: function CallExpression(node) {
      if (node.callee.type === 'MemberExpression' && node.callee.property.name === 'localizeMessage' || node.callee.name === 'localizeMessage') {
        var id = node.arguments[0] && node.arguments[0].value;
        var defaultMessage = node.arguments[1] && node.arguments[1].value;

        if (id && id !== '') {
          translations[id] = defaultMessage;
        } else {// console.log(node.arguments);
        }
      } else if (node.callee.type === 'MemberExpression' && node.callee.property.name === 'formatMessage' || node.callee.name === 'formatMessage') {
        if (node.arguments && node.arguments[0] && node.arguments[0].properties) {
          var _id = '';
          var _defaultMessage = '';
          var _iteratorNormalCompletion2 = true;
          var _didIteratorError2 = false;
          var _iteratorError2 = undefined;

          try {
            for (var _iterator2 = node.arguments[0].properties[Symbol.iterator](), _step2; !(_iteratorNormalCompletion2 = (_step2 = _iterator2.next()).done); _iteratorNormalCompletion2 = true) {
              var prop = _step2.value;

              // let prop = node.arguments[0].properties[idx]
              if (prop.value && prop.key && prop.key.name === 'id') {
                _id = prop.value.value;
              }

              if (prop.value && prop.key && prop.key.name === 'defaultMessage') {
                _defaultMessage = prop.value.value;
              }
            }
          } catch (err) {
            _didIteratorError2 = true;
            _iteratorError2 = err;
          } finally {
            try {
              if (!_iteratorNormalCompletion2 && _iterator2.return != null) {
                _iterator2.return();
              }
            } finally {
              if (_didIteratorError2) {
                throw _iteratorError2;
              }
            }
          }

          if (_id && _id !== '') {
            translations[_id] = _defaultMessage;
          }
        }
      } else if (node.callee.name === 't') {
        var _id2 = node.arguments[0] && node.arguments[0].value;

        translations[_id2] = '';
      }
    },
    JSXOpeningElement: function JSXOpeningElement(node) {
      var translatableProps = translatableComponents[node.name.name] || [];
      var _iteratorNormalCompletion3 = true;
      var _didIteratorError3 = false;
      var _iteratorError3 = undefined;

      try {
        for (var _iterator3 = translatableProps[Symbol.iterator](), _step3; !(_iteratorNormalCompletion3 = (_step3 = _iterator3.next()).done); _iteratorNormalCompletion3 = true) {
          var translatableProp = _step3.value;
          var id = '';
          var defaultMessage = '';

          if (typeof translatableProp === 'string') {
            var _iteratorNormalCompletion4 = true;
            var _didIteratorError4 = false;
            var _iteratorError4 = undefined;

            try {
              for (var _iterator4 = node.attributes[Symbol.iterator](), _step4; !(_iteratorNormalCompletion4 = (_step4 = _iterator4.next()).done); _iteratorNormalCompletion4 = true) {
                var attribute = _step4.value;

                if (attribute.value && attribute.value.expression && attribute.value.expression.value && attribute.name && attribute.name.name === translatableProp) {
                  id = attribute.value.expression.value.id;
                  defaultMessage = attribute.value.expression.value.defaultMessage;
                }

                if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === translatableProp) {
                  id = attribute.value.value.id;
                  defaultMessage = attribute.value.value.defaultMessage;
                }
              }
            } catch (err) {
              _didIteratorError4 = true;
              _iteratorError4 = err;
            } finally {
              try {
                if (!_iteratorNormalCompletion4 && _iterator4.return != null) {
                  _iterator4.return();
                }
              } finally {
                if (_didIteratorError4) {
                  throw _iteratorError4;
                }
              }
            }
          } else {
            var _iteratorNormalCompletion5 = true;
            var _didIteratorError5 = false;
            var _iteratorError5 = undefined;

            try {
              for (var _iterator5 = node.attributes[Symbol.iterator](), _step5; !(_iteratorNormalCompletion5 = (_step5 = _iterator5.next()).done); _iteratorNormalCompletion5 = true) {
                var attribute = _step5.value;

                if (attribute.value && attribute.value.expression && attribute.name && attribute.name.name === translatableProp.id) {
                  id = attribute.value.expression.value;
                }

                if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === translatableProp.id) {
                  id = attribute.value.value;
                }

                if (attribute.value && attribute.value.expression && attribute.name && attribute.name.name === translatableProp.default) {
                  defaultMessage = attribute.value.expression.value;
                }

                if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === translatableProp.default) {
                  defaultMessage = attribute.value.value;
                }
              }
            } catch (err) {
              _didIteratorError5 = true;
              _iteratorError5 = err;
            } finally {
              try {
                if (!_iteratorNormalCompletion5 && _iterator5.return != null) {
                  _iterator5.return();
                }
              } finally {
                if (_didIteratorError5) {
                  throw _iteratorError5;
                }
              }
            }
          }

          if (id) {
            translations[id] = defaultMessage;
          }
        }
      } catch (err) {
        _didIteratorError3 = true;
        _iteratorError3 = err;
      } finally {
        try {
          if (!_iteratorNormalCompletion3 && _iterator3.return != null) {
            _iterator3.return();
          }
        } finally {
          if (_didIteratorError3) {
            throw _iteratorError3;
          }
        }
      }
    }
  });
  return translations;
}
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.i18nCheck = i18nCheck;
exports.i18nCheckWebapp = i18nCheckWebapp;
exports.i18nCheckMobile = i18nCheckMobile;
exports.i18nExtractWebapp = i18nExtractWebapp;
exports.i18nExtractMobile = i18nExtractMobile;
exports.i18nCombine = i18nCombine;
exports.i18nSort = i18nSort;
exports.i18nSplit = i18nSplit;

function _slicedToArray(arr, i) { return _arrayWithHoles(arr) || _iterableToArrayLimit(arr, i) || _nonIterableRest(); }

function _nonIterableRest() { throw new TypeError("Invalid attempt to destructure non-iterable instance"); }

function _iterableToArrayLimit(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"] != null) _i["return"](); } finally { if (_d) throw _e; } } return _arr; }

function _arrayWithHoles(arr) { if (Array.isArray(arr)) return arr; }

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
var fs = require('fs');

var path = require('path');

var sortJson = require('sort-json');

var i18nExtractLib = require('./i18n_extract');

function difference(setA, setB) {
  var differenceSet = new Set(setA);
  var _iteratorNormalCompletion = true;
  var _didIteratorError = false;
  var _iteratorError = undefined;

  try {
    for (var _iterator = setB[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
      var elem = _step.value;
      differenceSet.delete(elem);
    }
  } catch (err) {
    _didIteratorError = true;
    _iteratorError = err;
  } finally {
    try {
      if (!_iteratorNormalCompletion && _iterator.return != null) {
        _iterator.return();
      }
    } finally {
      if (_didIteratorError) {
        throw _iteratorError;
      }
    }
  }

  return differenceSet;
}

function getCurrentTranslations(webappDir, mobileDir) {
  var currentWebappTranslationsJson = fs.readFileSync(path.join(webappDir, 'i18n', 'en.json'));
  var currentWebappTranslations = JSON.parse(currentWebappTranslationsJson);
  var currentMobileTranslationsJson = fs.readFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', 'en.json'));
  var currentMobileTranslations = JSON.parse(currentMobileTranslationsJson);
  return {
    webapp: currentWebappTranslations,
    mobile: currentMobileTranslations
  };
}

function i18nCheck(argv) {
  var webappDir = argv['webapp-dir'];
  var mobileDir = argv['mobile-dir'];
  var currentTranslations = getCurrentTranslations(webappDir, mobileDir);
  var currentWebappKeys = new Set(Object.keys(currentTranslations.webapp));
  var currentMobileKeys = new Set(Object.keys(currentTranslations.mobile));
  var promise1 = i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']);
  var promise2 = i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []);
  Promise.all([promise1, promise2]).then(function (_ref) {
    var _ref2 = _slicedToArray(_ref, 2),
        translationsWebapp = _ref2[0],
        translationsMobile = _ref2[1];

    var webappKeys = new Set(Object.keys(translationsWebapp));
    var mobileKeys = new Set(Object.keys(translationsMobile));
    var _iteratorNormalCompletion2 = true;
    var _didIteratorError2 = false;
    var _iteratorError2 = undefined;

    try {
      for (var _iterator2 = difference(currentWebappKeys, webappKeys)[Symbol.iterator](), _step2; !(_iteratorNormalCompletion2 = (_step2 = _iterator2.next()).done); _iteratorNormalCompletion2 = true) {
        var key = _step2.value;
        // eslint-disable-next-line no-console
        console.log('Removed from webapp:', key);
      }
    } catch (err) {
      _didIteratorError2 = true;
      _iteratorError2 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion2 && _iterator2.return != null) {
          _iterator2.return();
        }
      } finally {
        if (_didIteratorError2) {
          throw _iteratorError2;
        }
      }
    }

    var _iteratorNormalCompletion3 = true;
    var _didIteratorError3 = false;
    var _iteratorError3 = undefined;

    try {
      for (var _iterator3 = difference(webappKeys, currentWebappKeys)[Symbol.iterator](), _step3; !(_iteratorNormalCompletion3 = (_step3 = _iterator3.next()).done); _iteratorNormalCompletion3 = true) {
        var _key = _step3.value;
        // eslint-disable-next-line no-console
        console.log('Added to webapp:', _key);
      }
    } catch (err) {
      _didIteratorError3 = true;
      _iteratorError3 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion3 && _iterator3.return != null) {
          _iterator3.return();
        }
      } finally {
        if (_didIteratorError3) {
          throw _iteratorError3;
        }
      }
    }

    var _iteratorNormalCompletion4 = true;
    var _didIteratorError4 = false;
    var _iteratorError4 = undefined;

    try {
      for (var _iterator4 = difference(currentMobileKeys, mobileKeys)[Symbol.iterator](), _step4; !(_iteratorNormalCompletion4 = (_step4 = _iterator4.next()).done); _iteratorNormalCompletion4 = true) {
        var _key2 = _step4.value;
        // eslint-disable-next-line no-console
        console.log('Removed from mobile:', _key2);
      }
    } catch (err) {
      _didIteratorError4 = true;
      _iteratorError4 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion4 && _iterator4.return != null) {
          _iterator4.return();
        }
      } finally {
        if (_didIteratorError4) {
          throw _iteratorError4;
        }
      }
    }

    var _iteratorNormalCompletion5 = true;
    var _didIteratorError5 = false;
    var _iteratorError5 = undefined;

    try {
      for (var _iterator5 = difference(mobileKeys, currentMobileKeys)[Symbol.iterator](), _step5; !(_iteratorNormalCompletion5 = (_step5 = _iterator5.next()).done); _iteratorNormalCompletion5 = true) {
        var _key3 = _step5.value;
        // eslint-disable-next-line no-console
        console.log('Added to mobile:', _key3);
      }
    } catch (err) {
      _didIteratorError5 = true;
      _iteratorError5 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion5 && _iterator5.return != null) {
          _iterator5.return();
        }
      } finally {
        if (_didIteratorError5) {
          throw _iteratorError5;
        }
      }
    }
  });
}

function i18nCheckWebapp(argv) {
  var webappDir = argv['webapp-dir'];
  var mobileDir = argv['mobile-dir'];
  var currentTranslations = getCurrentTranslations(webappDir, mobileDir);
  var currentWebappKeys = new Set(Object.keys(currentTranslations.webapp));
  i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']).then(function (translationsWebapp) {
    var webappKeys = new Set(Object.keys(translationsWebapp));
    var changed = false;
    var _iteratorNormalCompletion6 = true;
    var _didIteratorError6 = false;
    var _iteratorError6 = undefined;

    try {
      for (var _iterator6 = difference(currentWebappKeys, webappKeys)[Symbol.iterator](), _step6; !(_iteratorNormalCompletion6 = (_step6 = _iterator6.next()).done); _iteratorNormalCompletion6 = true) {
        var key = _step6.value;
        // eslint-disable-next-line no-console
        console.log('Removed from webapp:', key);
        changed = true;
      }
    } catch (err) {
      _didIteratorError6 = true;
      _iteratorError6 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion6 && _iterator6.return != null) {
          _iterator6.return();
        }
      } finally {
        if (_didIteratorError6) {
          throw _iteratorError6;
        }
      }
    }

    var _iteratorNormalCompletion7 = true;
    var _didIteratorError7 = false;
    var _iteratorError7 = undefined;

    try {
      for (var _iterator7 = difference(webappKeys, currentWebappKeys)[Symbol.iterator](), _step7; !(_iteratorNormalCompletion7 = (_step7 = _iterator7.next()).done); _iteratorNormalCompletion7 = true) {
        var _key4 = _step7.value;
        // eslint-disable-next-line no-console
        console.log('Added to webapp:', _key4);
        changed = true;
      }
    } catch (err) {
      _didIteratorError7 = true;
      _iteratorError7 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion7 && _iterator7.return != null) {
          _iterator7.return();
        }
      } finally {
        if (_didIteratorError7) {
          throw _iteratorError7;
        }
      }
    }

    if (changed) {
      console.log('Changes found');
      process.exit(1);
    }
  });
}

function i18nCheckMobile(argv) {
  var webappDir = argv['webapp-dir'];
  var mobileDir = argv['mobile-dir'];
  var currentTranslations = getCurrentTranslations(webappDir, mobileDir);
  var currentMobileKeys = new Set(Object.keys(currentTranslations.mobile));
  i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []).then(function (translationsMobile) {
    var mobileKeys = new Set(Object.keys(translationsMobile));
    var changed = false;
    var _iteratorNormalCompletion8 = true;
    var _didIteratorError8 = false;
    var _iteratorError8 = undefined;

    try {
      for (var _iterator8 = difference(currentMobileKeys, mobileKeys)[Symbol.iterator](), _step8; !(_iteratorNormalCompletion8 = (_step8 = _iterator8.next()).done); _iteratorNormalCompletion8 = true) {
        var key = _step8.value;
        // eslint-disable-next-line no-console
        console.log('Removed from mobile:', key);
        changed = true;
      }
    } catch (err) {
      _didIteratorError8 = true;
      _iteratorError8 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion8 && _iterator8.return != null) {
          _iterator8.return();
        }
      } finally {
        if (_didIteratorError8) {
          throw _iteratorError8;
        }
      }
    }

    var _iteratorNormalCompletion9 = true;
    var _didIteratorError9 = false;
    var _iteratorError9 = undefined;

    try {
      for (var _iterator9 = difference(mobileKeys, currentMobileKeys)[Symbol.iterator](), _step9; !(_iteratorNormalCompletion9 = (_step9 = _iterator9.next()).done); _iteratorNormalCompletion9 = true) {
        var _key5 = _step9.value;
        // eslint-disable-next-line no-console
        console.log('Added to mobile:', _key5);
        changed = true;
      }
    } catch (err) {
      _didIteratorError9 = true;
      _iteratorError9 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion9 && _iterator9.return != null) {
          _iterator9.return();
        }
      } finally {
        if (_didIteratorError9) {
          throw _iteratorError9;
        }
      }
    }

    if (changed) {
      console.log('Changes found');
      process.exit(1);
    }
  });
}

function i18nExtractWebapp(argv) {
  var webappDir = argv['webapp-dir'];
  var mobileDir = argv['mobile-dir'];
  var currentTranslations = getCurrentTranslations(webappDir, mobileDir);
  var currentWebappKeys = new Set(Object.keys(currentTranslations.webapp));
  i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']).then(function (translationsWebapp) {
    var webappKeys = new Set(Object.keys(translationsWebapp));
    var _iteratorNormalCompletion10 = true;
    var _didIteratorError10 = false;
    var _iteratorError10 = undefined;

    try {
      for (var _iterator10 = difference(currentWebappKeys, webappKeys)[Symbol.iterator](), _step10; !(_iteratorNormalCompletion10 = (_step10 = _iterator10.next()).done); _iteratorNormalCompletion10 = true) {
        var key = _step10.value;
        delete currentTranslations.webapp[key];
      }
    } catch (err) {
      _didIteratorError10 = true;
      _iteratorError10 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion10 && _iterator10.return != null) {
          _iterator10.return();
        }
      } finally {
        if (_didIteratorError10) {
          throw _iteratorError10;
        }
      }
    }

    var _iteratorNormalCompletion11 = true;
    var _didIteratorError11 = false;
    var _iteratorError11 = undefined;

    try {
      for (var _iterator11 = difference(webappKeys, currentWebappKeys)[Symbol.iterator](), _step11; !(_iteratorNormalCompletion11 = (_step11 = _iterator11.next()).done); _iteratorNormalCompletion11 = true) {
        var _key6 = _step11.value;
        currentTranslations.webapp[_key6] = translationsWebapp[_key6];
      }
    } catch (err) {
      _didIteratorError11 = true;
      _iteratorError11 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion11 && _iterator11.return != null) {
          _iterator11.return();
        }
      } finally {
        if (_didIteratorError11) {
          throw _iteratorError11;
        }
      }
    }

    var options = {
      ignoreCase: true,
      reverse: false,
      depth: 1
    };
    var sortedWebappTranslations = sortJson(currentTranslations.webapp, options);
    fs.writeFileSync(path.join(webappDir, 'i18n', 'en.json'), JSON.stringify(sortedWebappTranslations, null, 2));
  });
}

function i18nExtractMobile(argv) {
  var webappDir = argv['webapp-dir'];
  var mobileDir = argv['mobile-dir'];
  var currentTranslations = getCurrentTranslations(webappDir, mobileDir);
  var currentMobileKeys = new Set(Object.keys(currentTranslations.mobile));
  i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []).then(function (translationsMobile) {
    var mobileKeys = new Set(Object.keys(translationsMobile));
    var _iteratorNormalCompletion12 = true;
    var _didIteratorError12 = false;
    var _iteratorError12 = undefined;

    try {
      for (var _iterator12 = difference(currentMobileKeys, mobileKeys)[Symbol.iterator](), _step12; !(_iteratorNormalCompletion12 = (_step12 = _iterator12.next()).done); _iteratorNormalCompletion12 = true) {
        var key = _step12.value;
        delete currentTranslations.mobile[key];
      }
    } catch (err) {
      _didIteratorError12 = true;
      _iteratorError12 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion12 && _iterator12.return != null) {
          _iterator12.return();
        }
      } finally {
        if (_didIteratorError12) {
          throw _iteratorError12;
        }
      }
    }

    var _iteratorNormalCompletion13 = true;
    var _didIteratorError13 = false;
    var _iteratorError13 = undefined;

    try {
      for (var _iterator13 = difference(mobileKeys, currentMobileKeys)[Symbol.iterator](), _step13; !(_iteratorNormalCompletion13 = (_step13 = _iterator13.next()).done); _iteratorNormalCompletion13 = true) {
        var _key7 = _step13.value;
        currentTranslations.mobile[_key7] = translationsMobile[_key7];
      }
    } catch (err) {
      _didIteratorError13 = true;
      _iteratorError13 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion13 && _iterator13.return != null) {
          _iterator13.return();
        }
      } finally {
        if (_didIteratorError13) {
          throw _iteratorError13;
        }
      }
    }

    var options = {
      ignoreCase: true,
      reverse: false,
      depth: 1
    };
    var sortedMobileTranslations = sortJson(currentTranslations.mobile, options);
    fs.writeFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', 'en.json'), JSON.stringify(sortedMobileTranslations, null, 2));
  });
}

function i18nCombine(argv) {
  var outputFile = argv.output;
  var translations = {};
  var _iteratorNormalCompletion14 = true;
  var _didIteratorError14 = false;
  var _iteratorError14 = undefined;

  try {
    for (var _iterator14 = argv._.slice(2)[Symbol.iterator](), _step14; !(_iteratorNormalCompletion14 = (_step14 = _iterator14.next()).done); _iteratorNormalCompletion14 = true) {
      var file = _step14.value;
      var itemTranslationsJson = fs.readFileSync(file);
      var itemTranslations = JSON.parse(itemTranslationsJson);

      for (var key in itemTranslations) {
        if ({}.hasOwnProperty.call(itemTranslations, key)) {
          translations[key] = itemTranslations[key];
        }
      }
    }
  } catch (err) {
    _didIteratorError14 = true;
    _iteratorError14 = err;
  } finally {
    try {
      if (!_iteratorNormalCompletion14 && _iterator14.return != null) {
        _iterator14.return();
      }
    } finally {
      if (_didIteratorError14) {
        throw _iteratorError14;
      }
    }
  }

  var options = {
    ignoreCase: true,
    reverse: false,
    depth: 1
  };
  var sortedTranslations = sortJson(translations, options);
  fs.writeFileSync(outputFile, JSON.stringify(sortedTranslations, null, 2));
}

function i18nSort(argv) {
  var outputFile = argv.output;
  var file = argv._[2];
  var itemTranslationsJson = fs.readFileSync(file);
  var itemTranslations = JSON.parse(itemTranslationsJson);
  var options = {
    ignoreCase: true,
    reverse: false,
    depth: 1
  };
  var sortedTranslations = sortJson(itemTranslations, options);
  fs.writeFileSync(outputFile, JSON.stringify(sortedTranslations, null, 2));
}

function i18nSplit(argv) {
  var webappDir = argv['webapp-dir'];
  var mobileDir = argv['mobile-dir'];
  var inputFiles = argv.inputs.split(',');
  var promise1 = i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']);
  var promise2 = i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []);
  Promise.all([promise1, promise2]).then(function (_ref3) {
    var _ref4 = _slicedToArray(_ref3, 2),
        translationsWebapp = _ref4[0],
        translationsMobile = _ref4[1];

    var _iteratorNormalCompletion15 = true;
    var _didIteratorError15 = false;
    var _iteratorError15 = undefined;

    try {
      for (var _iterator15 = inputFiles[Symbol.iterator](), _step15; !(_iteratorNormalCompletion15 = (_step15 = _iterator15.next()).done); _iteratorNormalCompletion15 = true) {
        var inputFile = _step15.value;
        var filename = path.basename(inputFile.trim());
        var allTranslationsJson = fs.readFileSync(inputFile.trim());
        var allTranslations = JSON.parse(allTranslationsJson);
        var webappKeys = new Set(Object.keys(translationsWebapp));
        var mobileKeys = new Set(Object.keys(translationsMobile));
        var translationsWebappOutput = {};
        var _iteratorNormalCompletion16 = true;
        var _didIteratorError16 = false;
        var _iteratorError16 = undefined;

        try {
          for (var _iterator16 = webappKeys[Symbol.iterator](), _step16; !(_iteratorNormalCompletion16 = (_step16 = _iterator16.next()).done); _iteratorNormalCompletion16 = true) {
            var key = _step16.value;
            translationsWebappOutput[key] = allTranslations[key];
          }
        } catch (err) {
          _didIteratorError16 = true;
          _iteratorError16 = err;
        } finally {
          try {
            if (!_iteratorNormalCompletion16 && _iterator16.return != null) {
              _iterator16.return();
            }
          } finally {
            if (_didIteratorError16) {
              throw _iteratorError16;
            }
          }
        }

        var translationsMobileOutput = {};
        var _iteratorNormalCompletion17 = true;
        var _didIteratorError17 = false;
        var _iteratorError17 = undefined;

        try {
          for (var _iterator17 = mobileKeys[Symbol.iterator](), _step17; !(_iteratorNormalCompletion17 = (_step17 = _iterator17.next()).done); _iteratorNormalCompletion17 = true) {
            var _key8 = _step17.value;
            translationsMobileOutput[_key8] = allTranslations[_key8];
          }
        } catch (err) {
          _didIteratorError17 = true;
          _iteratorError17 = err;
        } finally {
          try {
            if (!_iteratorNormalCompletion17 && _iterator17.return != null) {
              _iterator17.return();
            }
          } finally {
            if (_didIteratorError17) {
              throw _iteratorError17;
            }
          }
        }

        var options = {
          ignoreCase: true,
          reverse: false,
          depth: 1
        };
        var sortedWebappTranslations = sortJson(translationsWebappOutput, options);
        var sortedMobileTranslations = sortJson(translationsMobileOutput, options);
        fs.writeFileSync(path.join(webappDir, 'i18n', filename), JSON.stringify(sortedWebappTranslations, null, 2));
        fs.writeFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', filename), JSON.stringify(sortedMobileTranslations, null, 2));
      }
    } catch (err) {
      _didIteratorError15 = true;
      _iteratorError15 = err;
    } finally {
      try {
        if (!_iteratorNormalCompletion15 && _iterator15.return != null) {
          _iterator15.return();
        }
      } finally {
        if (_didIteratorError15) {
          throw _iteratorError15;
        }
      }
    }
  });
}
