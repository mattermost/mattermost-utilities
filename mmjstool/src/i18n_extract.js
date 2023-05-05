// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable no-console */

import * as fs from 'fs';

import {parse} from '@typescript-eslint/typescript-estree';
import walk from 'estree-walk';
import * as FileHound from 'filehound';

const translatableComponents = {
    FormattedMessage: [{id: 'id', default: 'defaultMessage'}],
    FormattedHTMLMessage: [{id: 'id', default: 'defaultMessage'}],
    FormattedMarkdownMessage: [{id: 'id', default: 'defaultMessage'}],
    FormattedAdminHeader: [{id: 'id', default: 'defaultMessage'}],
    LocalizedInput: ['placeholder'],
    LocalizedIcon: ['title'],

    // Used in mattermost-mobile exclusively
    FormattedText: [{id: 'id', default: 'defaultMessage'}],
    FormattedMarkdownText: [{id: 'id', default: 'defaultMessage'}],
};

export function extractFromDirectory(dirPaths, filters = []) {
    return new Promise((resolve) => {
        const promises = dirPaths.map((dirPath) => {
            return new Promise((innerResolve) => {
                const translations = {};
                FileHound.create().
                    paths(dirPath).
                    discard(filters).
                    ext('js', 'jsx', 'ts', 'tsx').
                    find().
                    then((files) => {
                        for (const file of files) {
                            try {
                                Object.assign(translations, extractFromFile(file));
                            } catch (e) {
                                console.log(e);
                                console.log('Unable to parse file:', file);
                                console.log('Error in: line', e.loc && e.loc.line, 'column', e.loc && e.loc.column);
                                return;
                            }
                        }
                        innerResolve(translations);
                    });
            });
        });

        Promise.all(promises).then((translations) => {
            resolve(Object.assign({}, ...translations));
        });
    });
}

function extractFromFile(path) {
    const translations = {};

    var code = fs.readFileSync(path, 'utf-8');
    const ast = parse(code, {
        filePath: path,
        jsx: path.endsWith('.tsx') || path.endsWith('.jsx'),
    });

    walk(ast, {
        CallExpression: (node) => {
            if ((node.callee.type === 'MemberExpression' && node.callee.property.name === 'localizeMessage') ||
                node.callee.name === 'localizeMessage') {
                const id = node.arguments[0] && node.arguments[0].value;
                const defaultMessage = node.arguments[1] && node.arguments[1].value;

                if (id && id !== '') {
                    translations[id] = defaultMessage;
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
            } else if ((node.callee.type === 'MemberExpression' && node.callee.property.name === 'defineMessages') ||
                node.callee.name === 'defineMessages') {
                if (node.arguments && node.arguments[0] && node.arguments[0].properties && node.arguments[0].properties.length !== 0) {
                    for (const nodeProperty of node.arguments[0].properties) {
                        if (nodeProperty.type && nodeProperty.type === 'Property' && nodeProperty.key && nodeProperty.key.name !== '' &&
                            nodeProperty.value && nodeProperty.value.type === 'ObjectExpression' &&
                            nodeProperty.value.properties && nodeProperty.value.properties.length !== 0) {
                            let id = '';
                            let defaultMessage = '';

                            for (const property of nodeProperty.value.properties) {
                                if (property && property.type && property.type === 'Property' &&
                                    (property.key.name === 'id' || property.key.name === 'defaultMessage') &&
                                    property.key && property.key.type && property.key.type === 'Identifier' &&
                                    property.value && property.value.type && property.value.type === 'Literal' &&
                                    property.value.value !== '') {
                                    if (property.key.name === 'id') {
                                        id = property.value.value;
                                    } else if (property.key.name === 'defaultMessage') {
                                        defaultMessage = property.value.value;
                                    }
                                }
                            }

                            if (id && id !== '' && defaultMessage && defaultMessage !== '') {
                                translations[id] = defaultMessage;
                            }
                        }
                    }
                }
            }
        },
        JSXOpeningElement: (node) => {
            const translatableProps = translatableComponents[node.name.name] || [];
            for (const translatableProp of translatableProps) {
                let id = '';
                let defaultMessage = '';

                if (typeof translatableProp === 'string') {
                    for (const attribute of node.attributes) {
                        if (attribute.value && attribute.value.expression && attribute.value.expression.value && attribute.name && attribute.name.name === translatableProp) {
                            id = attribute.value.expression.value.id;
                            defaultMessage = attribute.value.expression.value.defaultMessage;
                        }
                        if (attribute.value && attribute.value.value && attribute.name && attribute.name.name === translatableProp) {
                            id = attribute.value.value.id;
                            defaultMessage = attribute.value.value.defaultMessage;
                        }
                    }
                } else {
                    for (const attribute of node.attributes) {
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
                }
                if (id) {
                    translations[id] = defaultMessage;
                }
            }
        },
    });
    return translations;
}
