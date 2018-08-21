// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

const fs = require('fs');

const FileHound = require('filehound');

const acorn = require('acorn');
const astwalk = require('acorn/dist/walk');
const {base} = require('acorn-jsx-walk');
const injectAcornStage3 = require('acorn-stage3/inject');
const injectAcornJsx = require('acorn-jsx/inject');
const injectAcornStaticClassPropertyInitializer = require('acorn-static-class-property-initializer/inject');

injectAcornStage3(acorn);
injectAcornJsx(acorn);
injectAcornStaticClassPropertyInitializer(acorn);

function patchAstWalk() {
    astwalk.base.ClassProperty = (node, st, c) => {
        c(node.key, st);
        c(node.value, st);
    };
    astwalk.base.FieldDefinition = (node, st, c) => {
        c(node.key, st);
        c(node.value, st);
    };
    astwalk.base.Import = () => { /* empty function */ };
    astwalk.base.JSXElement = (node, st, c) => {
        c(node.openingElement, st);
        node.children.forEach((n) => {
            c(n, st);
        });
    };

    astwalk.base.JSXOpeningElement = (node, st, c) => {
        node.attributes.forEach((n) => {
            c(n, st);
        });
    };

    astwalk.base.JSXAttribute = (node, st, c) => {
        c(node.name, st);
        c(node.value, st);
    };
    astwalk.base.JSXSpreadAttribute = (node, st, c) => {
        c(node.argument, st);
    };
    astwalk.base.JSXText = () => { /* empty function */ };
    astwalk.base.JSXIdentifier = () => { /* empty function */ };
    astwalk.base.JSXExpressionContainer = base.JSXExpressionContainer;
    astwalk.base.JSXEmptyExpression = () => { /* empty function */ };
}

patchAstWalk()

export function extractFromDirectory(dirPath, filters = []) {
    return new Promise((resolve) => {
        const translations = {};
        FileHound.create().
            paths(dirPath).
            discard(filters).
            ext('js', 'jsx').
            find().
            then((files) => {
                for (const file of files) {
                    try {
                        Object.assign(translations, extractFromFile(file));
                    } catch (e) {
                        console.log("Unable to parse file:", file);
                        console.log("Error in: line", e.loc.line, "column", e.loc.column);
                        return;
                    }
                }
                resolve(translations);
            });
    });
}

function extractFromFile(path) {
    const translations = {};

    var code = fs.readFileSync(path, 'utf-8');
    const ast = acorn.parse(code, {
        plugins: {
            stage3: true,
            jsx: true,
            staticClassPropertyInitializer: true,
        },
        ecmaVersion: 10,
        sourceType: 'module',
    });

    // Make it compatible with our source code
    astwalk.full(ast, (node, st, type) => {
        if (type === 'CallExpression') {
            if ((node.callee.type === 'MemberExpression' && node.callee.property.name === 'localizeMessage') ||
                node.callee.name === 'localizeMessage') {
                const id = node.arguments[0] && node.arguments[0].value;
                const defaultMessage = node.arguments[1] && node.arguments[1].value;

                if (id && id !== '') {
                    translations[id] = defaultMessage;
                } else {
                    // console.log(node.arguments);
                }
            } else if ((node.callee.type === 'MemberExpression' && node.callee.property.name === 'formatMessage') ||
                node.callee.name === 'formatMessage') {
                if (node.arguments && node.arguments[0] && node.arguments[0].properties) {
                    let id = '';
                    let defaultMessage = '';

                    for (const prop of node.arguments[0].properties) {
                        // let prop = node.arguments[0].properties[idx]
                        if (prop.value && prop.key && prop.key.name === 'id') {
                            id = prop.value.value;
                        }
                        if (prop.value && prop.key && prop.key.name === 'defaultMessage') {
                            defaultMessage = prop.value.value;
                        }
                    }
                    if (id && id !== '') {
                        translations[id] = defaultMessage;
                    }
                }
            } else if (node.callee.name === 't') {
                const id = node.arguments[0] && node.arguments[0].value;
                translations[id] = '';
            }
        }

        if (type === 'JSXOpeningElement') {
            if (node.name.name === 'FormattedText' || node.name.name === 'FormattedMessage' || node.name.name === 'FormattedHTMLMessage' || node.name.name === 'FormattedMarkdownMessage' || node.name.name === 'FormattedMarkdownText') {
                let id = '';
                let defaultMessage = '';
                for (var attribute of node.attributes) {
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

                if (id && id !== '') {
                    translations[id] = defaultMessage;
                }
            }
        }
    });
    return translations;
}
