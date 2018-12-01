'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.extractFromDirectory = extractFromDirectory;

function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

var fs = require('fs');

var FileHound = require('filehound');
var flowRemoveTypes = require('flow-remove-types');

var _require = require('acorn'),
    Parser = _require.Parser;

var astwalk = require('acorn-walk');

var _require2 = require('acorn-jsx-walk'),
    base = _require2.base;

var acornJsx = require('acorn-jsx');
var acornStage3 = require('acorn-stage3');

var _require3 = require('./acorn-optional-chaining'),
    acornOptionalChaining = _require3.acornOptionalChaining;

var acorn = Parser.extend(acornStage3, acornJsx(), acornOptionalChaining);

function patchAstWalk() {
    astwalk.base.ClassProperty = function (node, st, c) {
        c(node.key, st);
        c(node.value, st);
    };
    astwalk.base.FieldDefinition = function (node, st, c) {
        c(node.key, st);
        c(node.value, st);
    };
    astwalk.base.Import = function () {/* empty function */};
    astwalk.base.JSXElement = function (node, st, c) {
        c(node.openingElement, st);
        node.children.forEach(function (n) {
            c(n, st);
        });
    };

    astwalk.base.JSXOpeningElement = function (node, st, c) {
        node.attributes.forEach(function (n) {
            c(n, st);
        });
    };

    astwalk.base.JSXAttribute = function (node, st, c) {
        c(node.name, st);
        c(node.value, st);
    };
    astwalk.base.JSXSpreadAttribute = function (node, st, c) {
        c(node.argument, st);
    };
    astwalk.base.JSXText = function () {/* empty function */};
    astwalk.base.JSXIdentifier = function () {/* empty function */};
    astwalk.base.JSXExpressionContainer = base.JSXExpressionContainer;
    astwalk.base.JSXEmptyExpression = function () {/* empty function */};
}

patchAstWalk();

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
                                console.log("Unable to parse file:", file);
                                console.log("Error in: line", e.loc.line, "column", e.loc.column);
                                return;
                            }
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
    var ast = acorn.parse(flowRemoveTypes(code), {
        plugins: {
            stage3: true,
            jsx: true,
            staticClassPropertyInitializer: true,
            ignoreOptionalChaining: true
        },
        ecmaVersion: 10,
        sourceType: 'module'
    });

    // Make it compatible with our source code
    astwalk.full(ast, function (node, st, type) {
        if (type === 'CallExpression') {
            if (node.callee.type === 'MemberExpression' && node.callee.property.name === 'localizeMessage' || node.callee.name === 'localizeMessage') {
                var id = node.arguments[0] && node.arguments[0].value;
                var defaultMessage = node.arguments[1] && node.arguments[1].value;

                if (id && id !== '') {
                    translations[id] = defaultMessage;
                } else {
                    // console.log(node.arguments);
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
                            if (!_iteratorNormalCompletion2 && _iterator2.return) {
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
        }

        if (type === 'JSXOpeningElement') {
            if (node.name.name === 'FormattedText' || node.name.name === 'FormattedMessage' || node.name.name === 'FormattedHTMLMessage' || node.name.name === 'FormattedMarkdownMessage' || node.name.name === 'FormattedMarkdownText' || node.name.name === 'FormattedAdminHeader') {
                var _id3 = '';
                var _defaultMessage2 = '';
                var _iteratorNormalCompletion3 = true;
                var _didIteratorError3 = false;
                var _iteratorError3 = undefined;

                try {
                    for (var _iterator3 = node.attributes[Symbol.iterator](), _step3; !(_iteratorNormalCompletion3 = (_step3 = _iterator3.next()).done); _iteratorNormalCompletion3 = true) {
                        var attribute = _step3.value;

                        if (attribute.value && attribute.value.expression && attribute.name && attribute.name.name === 'id') {
                            _id3 = attribute.value.expression.value;
                        }
                        if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === 'id') {
                            _id3 = attribute.value.value;
                        }
                        if (attribute.value && attribute.value.expression && attribute.name && attribute.name.name === 'defaultMessage') {
                            _defaultMessage2 = attribute.value.expression.value;
                        }
                        if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === 'defaultMessage') {
                            _defaultMessage2 = attribute.value.value;
                        }
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

                if (_id3 && _id3 !== '') {
                    translations[_id3] = _defaultMessage2;
                }
            }
        }
    });
    return translations;
}