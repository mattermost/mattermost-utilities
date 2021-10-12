// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable no-console */

const fs = require('fs');

const FileHound = require('filehound');

const Parser = require('flow-parser');
const walk = require('estree-walk');

const translatableComponents = {
    FormattedText: [{id: 'id', default: 'defaultMessage'}],
    FormattedMessage: [{id: 'id', default: 'defaultMessage'}],
    FormattedHTMLMessage: [{id: 'id', default: 'defaultMessage'}],
    FormattedMarkdownMessage: [{id: 'id', default: 'defaultMessage'}],
    FormattedMarkdownText: [{id: 'id', default: 'defaultMessage'}],
    FormattedAdminHeader: [{id: 'id', default: 'defaultMessage'}],
    LocalizedInput: ['placeholder'],
    LocalizedIcon: ['title'],
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
    const ast = Parser.parse(code, {
        esproposal_class_static_fields: true,
        esproposal_class_instance_fields: true,
        esproposal_optional_chaining: true,
    });

    walk(ast, {
        CallExpression: (node) => {
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

export function translateToQA(translations) {
    qa = {};
    for (const key of translations) {
        let translation = translations[key];
        for (let i = translation.length; i >= 0; i--) {
            switch (translation[i]) {
                case 'a':
                    translation[i] = 'á'
                    break;
                case 'A':
                    translation[i] = 'Á'
                    break;
                case 'e':
                    translation[i] = 'é'
                    break;
                case 'E':
                    translation[i] = 'É'
                    break;
                case 'i':
                    translation[i] = 'í'
                    break;
                case 'I':
                    translation[i] = 'Í'
                    break;
                case 'o':
                    translation[i] = 'ó'
                    break;
                case 'O':
                    translation[i] = 'Ó'
                    break;
                case 'u':
                    translation[i] = 'ú'
                    break;
                case 'U':
                    translation[i] = 'Ú'
                    break;
                case 'n':
                    translation[i] = 'ñ'
                    break;
                case 'N':
                    translation[i] = 'Ñ'
                    break;
                case '?':
                    translation[i] = '¿'
                    break;
                case '!':
                    translation[i] = '¡'
                    break;
            }
        }
        qa[key] = translation
    }
    return qa
}
