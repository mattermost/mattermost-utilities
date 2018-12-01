'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});

var _slicedToArray = function () { function sliceIterator(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"]) _i["return"](); } finally { if (_d) throw _e; } } return _arr; } return function (arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { return sliceIterator(arr, i); } else { throw new TypeError("Invalid attempt to destructure non-iterable instance"); } }; }();

exports.i18nCheck = i18nCheck;
exports.i18nExtractWebapp = i18nExtractWebapp;
exports.i18nExtractMobile = i18nExtractMobile;
exports.i18nCombine = i18nCombine;
exports.i18nSort = i18nSort;
exports.i18nSplit = i18nSplit;
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
            if (!_iteratorNormalCompletion && _iterator.return) {
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
                if (!_iteratorNormalCompletion2 && _iterator2.return) {
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
                if (!_iteratorNormalCompletion3 && _iterator3.return) {
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
                if (!_iteratorNormalCompletion4 && _iterator4.return) {
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
                if (!_iteratorNormalCompletion5 && _iterator5.return) {
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

function i18nExtractWebapp(argv) {
    var webappDir = argv['webapp-dir'];
    var mobileDir = argv['mobile-dir'];

    var currentTranslations = getCurrentTranslations(webappDir, mobileDir);
    var currentWebappKeys = new Set(Object.keys(currentTranslations.webapp));

    i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']).then(function (translationsWebapp) {
        var webappKeys = new Set(Object.keys(translationsWebapp));

        var _iteratorNormalCompletion6 = true;
        var _didIteratorError6 = false;
        var _iteratorError6 = undefined;

        try {
            for (var _iterator6 = difference(currentWebappKeys, webappKeys)[Symbol.iterator](), _step6; !(_iteratorNormalCompletion6 = (_step6 = _iterator6.next()).done); _iteratorNormalCompletion6 = true) {
                var key = _step6.value;

                delete currentTranslations.webapp[key];
            }
        } catch (err) {
            _didIteratorError6 = true;
            _iteratorError6 = err;
        } finally {
            try {
                if (!_iteratorNormalCompletion6 && _iterator6.return) {
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

                currentTranslations.webapp[_key4] = translationsWebapp[_key4];
            }
        } catch (err) {
            _didIteratorError7 = true;
            _iteratorError7 = err;
        } finally {
            try {
                if (!_iteratorNormalCompletion7 && _iterator7.return) {
                    _iterator7.return();
                }
            } finally {
                if (_didIteratorError7) {
                    throw _iteratorError7;
                }
            }
        }

        var options = { ignoreCase: true, reverse: false, depth: 1 };
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

        var _iteratorNormalCompletion8 = true;
        var _didIteratorError8 = false;
        var _iteratorError8 = undefined;

        try {
            for (var _iterator8 = difference(currentMobileKeys, mobileKeys)[Symbol.iterator](), _step8; !(_iteratorNormalCompletion8 = (_step8 = _iterator8.next()).done); _iteratorNormalCompletion8 = true) {
                var key = _step8.value;

                delete currentTranslations.mobile[key];
            }
        } catch (err) {
            _didIteratorError8 = true;
            _iteratorError8 = err;
        } finally {
            try {
                if (!_iteratorNormalCompletion8 && _iterator8.return) {
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

                currentTranslations.mobile[_key5] = translationsMobile[_key5];
            }
        } catch (err) {
            _didIteratorError9 = true;
            _iteratorError9 = err;
        } finally {
            try {
                if (!_iteratorNormalCompletion9 && _iterator9.return) {
                    _iterator9.return();
                }
            } finally {
                if (_didIteratorError9) {
                    throw _iteratorError9;
                }
            }
        }

        var options = { ignoreCase: true, reverse: false, depth: 1 };
        var sortedMobileTranslations = sortJson(currentTranslations.mobile, options);
        fs.writeFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', 'en.json'), JSON.stringify(sortedMobileTranslations, null, 2));
    });
}

function i18nCombine(argv) {
    var outputFile = argv.output;

    var translations = {};

    var _iteratorNormalCompletion10 = true;
    var _didIteratorError10 = false;
    var _iteratorError10 = undefined;

    try {
        for (var _iterator10 = argv._.slice(2)[Symbol.iterator](), _step10; !(_iteratorNormalCompletion10 = (_step10 = _iterator10.next()).done); _iteratorNormalCompletion10 = true) {
            var file = _step10.value;

            var itemTranslationsJson = fs.readFileSync(file);
            var itemTranslations = JSON.parse(itemTranslationsJson);

            for (var key in itemTranslations) {
                if ({}.hasOwnProperty.call(itemTranslations, key)) {
                    translations[key] = itemTranslations[key];
                }
            }
        }
    } catch (err) {
        _didIteratorError10 = true;
        _iteratorError10 = err;
    } finally {
        try {
            if (!_iteratorNormalCompletion10 && _iterator10.return) {
                _iterator10.return();
            }
        } finally {
            if (_didIteratorError10) {
                throw _iteratorError10;
            }
        }
    }

    var options = { ignoreCase: true, reverse: false, depth: 1 };
    var sortedTranslations = sortJson(translations, options);
    fs.writeFileSync(outputFile, JSON.stringify(sortedTranslations, null, 2));
}

function i18nSort(argv) {
    var outputFile = argv.output;

    var file = argv._[2];
    var itemTranslationsJson = fs.readFileSync(file);
    var itemTranslations = JSON.parse(itemTranslationsJson);

    var options = { ignoreCase: true, reverse: false, depth: 1 };
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

        var _iteratorNormalCompletion11 = true;
        var _didIteratorError11 = false;
        var _iteratorError11 = undefined;

        try {
            for (var _iterator11 = inputFiles[Symbol.iterator](), _step11; !(_iteratorNormalCompletion11 = (_step11 = _iterator11.next()).done); _iteratorNormalCompletion11 = true) {
                var inputFile = _step11.value;

                var filename = path.basename(inputFile.trim());
                var allTranslationsJson = fs.readFileSync(inputFile.trim());
                var allTranslations = JSON.parse(allTranslationsJson);

                var webappKeys = new Set(Object.keys(translationsWebapp));
                var mobileKeys = new Set(Object.keys(translationsMobile));

                var translationsWebappOutput = {};
                var _iteratorNormalCompletion12 = true;
                var _didIteratorError12 = false;
                var _iteratorError12 = undefined;

                try {
                    for (var _iterator12 = webappKeys[Symbol.iterator](), _step12; !(_iteratorNormalCompletion12 = (_step12 = _iterator12.next()).done); _iteratorNormalCompletion12 = true) {
                        var key = _step12.value;

                        translationsWebappOutput[key] = allTranslations[key];
                    }
                } catch (err) {
                    _didIteratorError12 = true;
                    _iteratorError12 = err;
                } finally {
                    try {
                        if (!_iteratorNormalCompletion12 && _iterator12.return) {
                            _iterator12.return();
                        }
                    } finally {
                        if (_didIteratorError12) {
                            throw _iteratorError12;
                        }
                    }
                }

                var translationsMobileOutput = {};
                var _iteratorNormalCompletion13 = true;
                var _didIteratorError13 = false;
                var _iteratorError13 = undefined;

                try {
                    for (var _iterator13 = mobileKeys[Symbol.iterator](), _step13; !(_iteratorNormalCompletion13 = (_step13 = _iterator13.next()).done); _iteratorNormalCompletion13 = true) {
                        var _key6 = _step13.value;

                        translationsMobileOutput[_key6] = allTranslations[_key6];
                    }
                } catch (err) {
                    _didIteratorError13 = true;
                    _iteratorError13 = err;
                } finally {
                    try {
                        if (!_iteratorNormalCompletion13 && _iterator13.return) {
                            _iterator13.return();
                        }
                    } finally {
                        if (_didIteratorError13) {
                            throw _iteratorError13;
                        }
                    }
                }

                var options = { ignoreCase: true, reverse: false, depth: 1 };
                var sortedWebappTranslations = sortJson(translationsWebappOutput, options);
                var sortedMobileTranslations = sortJson(translationsMobileOutput, options);
                fs.writeFileSync(path.join(webappDir, 'i18n', filename), JSON.stringify(sortedWebappTranslations, null, 2));
                fs.writeFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', filename), JSON.stringify(sortedMobileTranslations, null, 2));
            }
        } catch (err) {
            _didIteratorError11 = true;
            _iteratorError11 = err;
        } finally {
            try {
                if (!_iteratorNormalCompletion11 && _iterator11.return) {
                    _iterator11.return();
                }
            } finally {
                if (_didIteratorError11) {
                    throw _iteratorError11;
                }
            }
        }
    });
}