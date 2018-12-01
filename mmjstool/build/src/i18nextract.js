// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
'use strict';

Object.defineProperty(exports, '__esModule', {
    value: true,
});
exports.extractFromDirectory = extractFromDirectory;
exports.extractFromFile = extractFromFile;

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
var acorn = require('acorn');
var astwalk = require('acorn/dist/walk');

var _require = require('acorn-jsx-walk');

var base = _require.base;

var fs = require('fs');

var walk = require('walk');

var path = require('path');

var injectAcornStage3 = require('acorn-stage3/inject');
var injectAcornJsx = require('acorn-jsx/inject');
var injectAcornStaticClassPropertyInitializer = require('acorn-static-class-property-initializer/inject');

injectAcornStage3(acorn);
injectAcornJsx(acorn);
injectAcornStaticClassPropertyInitializer(acorn);

astwalk.base.ClassProperty = function(node) {};
astwalk.base.FieldDefinition = function(node) {};
astwalk.base.Pattern = function(node) {};
astwalk.base.Import = function(node) {};
astwalk.base.JSXElement = base.JSXElement;
astwalk.base.JSXText = function(node) {};
astwalk.base.JSXExpressionContainer = base.JSXExpressionContainer;
astwalk.base.JSXEmptyExpression = function(node) {};

function extractFromDirectory(dirPath) {
    var walker = walk.walk(dirPath, {filters: ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'android', 'ios', 'builds']});
    var translations = {};

    return new Promise(((resolve, reject) => {
        walker.on('file', (root, fileStats, next) => {
            if (fileStats.name.endsWith('.js') || fileStats.name.endsWith('.jsx')) {
                Object.assign(translations, extractFromFile(path.join(root, fileStats.name)));
            }
            next();
        });

        walker.on('end', () => {
            resolve(translations);
        });

        walker.on('errors', (root, nodeStatsArray, next) => {
            reject('Error accessing to your files or directories');
        });
    }));
}

function extractFromFile(path) {
    var translations = {};

    var code = fs.readFileSync(path, 'utf-8');
    var ast = acorn.parse(code, {
        plugins: {
            stage3: true,
            jsx: true,
            staticClassPropertyInitializer: true,
        },
        ecmaVersion: 10,
        sourceType: 'module',
    });

    // Make it compatible with our source code
    astwalk.simple(ast, {
        CallExpression: function CallExpression(node) {
            if (node.callee.type === 'MemberExpression' && node.callee.property.name === 'localizeMessage' || node.callee.name === 'localizeMessage') {
                var id = node.arguments[0] && node.arguments[0].value;
                var defaultMessage = node.arguments[1] && node.arguments[1].value;

                if (id && id !== '') {
                    translations[id] = defaultMessage;
                } else {
                    console.log(node.arguments);
                }
            } else if (node.callee.type === 'MemberExpression' && node.callee.property.name === 'formatMessage' || node.callee.name === 'formatMessage') {
                for (var idx in node.arguments[0].properties) {
                    var prop = node.arguments[0].properties[idx];
                    var _id = '';
                    var _defaultMessage = '';
                    if (prop.value && prop.key && prop.key.name === 'id') {
                        _id = prop.value.value;
                    }
                    if (prop.value && prop.key && prop.key.name === 'defaultMessage') {
                        _defaultMessage = prop.value.value;
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
        JSXElement: function JSXElement(node) {
            if (node.openingElement.name.name === 'FormattedMessage' || node.openingElement.name.name === 'FormattedHTMLMessage' || node.openingElement.name.name === 'FormattedMarkdownMessage') {
                var id = '';
                var defaultMessage = '';
                var _iteratorNormalCompletion = true;
                var _didIteratorError = false;
                var _iteratorError = undefined;

                try {
                    for (var _iterator = node.openingElement.attributes[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
                        var attribute = _step.value;

                        if (attribute.value && attribute.value.expression && attribute.name && attribute.name.name === 'id') {
                            id = attribute.value.expression.value;
                        }
                        if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === 'id') {
                            id = attribute.value.value;
                        }
                        if (attribute.value && attribute.value.expression && attribute.name && attribute.name.name === 'defaultMessage') {
                            defaultMessage = attribute.value.expression.value;
                        }
                        if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === 'defaultMessage') {
                            defaultMessage = attribute.value.value;
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

                if (id && id !== '') {
                    translations[id] = defaultMessage;
                }
            } else {}
        },
    });
    return translations;
}