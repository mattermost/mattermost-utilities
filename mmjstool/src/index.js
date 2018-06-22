#!/usr/bin/env node

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
const yargs = require('yargs');

const i18nHandlers = require('./i18n_handlers');

/*eslint no-unused-vars: ["error", { "varsIgnorePattern": "[iI]gnored" }]*/
const ignored = yargs.
    usage('Usage: mmjstool <command> [options]').
    example('mmjstool i18n extract-webapp --webapp-dir ./', 'Extract all the i18n strings from the webapp source code').
    demandCommand(1).
    help('h').
    alias('h', 'help').
    command('i18n', 'I18n management commands', (i18nArgs) => {
        i18nArgs.
            demandCommand(1).
            command('extract-webapp',
                'Read the source code, find all the translations string from mattermost-webapp and write them to the file mattermost-webapp/i18n/en.json',
                () => { /* empty function */ },
                i18nHandlers.i18nExtractWebapp,
            ).
            command('extract-mobile',
                'Read the source code, find all the translations string from mattermost-mobile and write them to the file mattermost-mobile/assets/base/i18n/en.json.',
                () => { /* empty function */ },
                i18nHandlers.i18nExtractMobile,
            ).
            command('combine',
                'Read the translations string from mattermost-webapp and mattermost-mobile and combine them in a single file',
                (combineArgs) => {
                    combineArgs.
                        option('output', {
                            describe: 'File to store the combined translations',
                            default: 'en.json',
                        });
                },
                i18nHandlers.i18nCombine
            ).
            command('split',
                'Read a set of combined translation files, and split them in mattermost-server and mattermost-web translations',
                (splitArgs) => {
                    splitArgs.
                        option('inputs', {
                            describe: 'List of file to read the combined translations, splitted by ",". (e.g. en.json,es.json,fr.json)',
                            default: 'en.json',
                        });
                },
                i18nHandlers.i18nSplit,
            ).
            command('check',
                'Read the source code, find all the translations string, and show you the differences with the current i18n/en.json files',
                () => { /* empty function */ },
                i18nHandlers.i18nCheck,
            ).
            option('webapp-dir', {
                describe: 'webapp source code directory',
                default: '../mattermost-webapp',
            }).
            option('mobile-dir', {
                describe: 'webapp source code directory',
                default: '../mattermost-mobile',
            });
    }, () => { /* empty function */ }
    ).argv;
