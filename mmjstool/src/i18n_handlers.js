// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable no-console,no-process-exit */

const fs = require('fs');
const path = require('path');

const sortJson = require('sort-json');

const i18nExtractLib = require('./i18n_extract');

function difference(setA, setB) {
    var differenceSet = new Set(setA);
    for (var elem of setB) {
        differenceSet.delete(elem);
    }
    return differenceSet;
}

function getCurrentTranslationsWebapp(webappDir) {
    const currentWebappTranslationsJson = fs.readFileSync(path.join(webappDir, 'i18n', 'en.json'));
    return JSON.parse(currentWebappTranslationsJson);
}

function getCurrentTranslationsMobile(mobileDir) {
    const currentMobileTranslationsJson = fs.readFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', 'en.json'));
    return JSON.parse(currentMobileTranslationsJson);
}

export function i18nCheckWebapp(argv) {
    const webappDir = argv['webapp-dir'];

    const currentTranslations = getCurrentTranslationsWebapp(webappDir);
    const currentWebappKeys = new Set(Object.keys(currentTranslations));

    i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']).then((translationsWebapp) => {
        const webappKeys = new Set(Object.keys(translationsWebapp));

        let changed = false;
        for (const key of difference(currentWebappKeys, webappKeys)) {
            // eslint-disable-next-line no-console
            console.log('Removed from webapp:', key);
            changed = true;
        }
        for (const key of difference(webappKeys, currentWebappKeys)) {
            // eslint-disable-next-line no-console
            console.log('Added to webapp:', key);
            changed = true;
        }
        if (changed) {
            console.log('Changes found');
            process.exit(1);
        }
    });
}

export function i18nCheck(argv) {
    i18nCheckWebapp(argv);
    i18nCheckMobile(argv);
}

export function i18nCheckMobile(argv) {
    const mobileDir = argv['mobile-dir'];

    const currentTranslations = getCurrentTranslationsMobile(mobileDir);
    const currentMobileKeys = new Set(Object.keys(currentTranslations));

    i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []).then((translationsMobile) => {
        const mobileKeys = new Set(Object.keys(translationsMobile));

        let changed = false;
        for (const key of difference(currentMobileKeys, mobileKeys)) {
            // eslint-disable-next-line no-console
            console.log('Removed from mobile:', key);
            changed = true;
        }
        for (const key of difference(mobileKeys, currentMobileKeys)) {
            // eslint-disable-next-line no-console
            console.log('Added to mobile:', key);
            changed = true;
        }

        if (changed) {
            console.log('Changes found');
            process.exit(1);
        }
    });
}

export function i18nExtractWebapp(argv) {
    const webappDir = argv['webapp-dir'];

    const currentTranslations = getCurrentTranslationsWebapp(webappDir);
    const currentWebappKeys = new Set(Object.keys(currentTranslations));

    i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']).then((translationsWebapp) => {
        const webappKeys = new Set(Object.keys(translationsWebapp));

        for (const key of difference(currentWebappKeys, webappKeys)) {
            delete currentTranslations[key];
        }
        for (const key of difference(webappKeys, currentWebappKeys)) {
            currentTranslations[key] = translationsWebapp[key];
        }

        const options = {ignoreCase: true, reverse: false, depth: 1};
        const sortedWebappTranslations = sortJson(currentTranslations, options);
        fs.writeFileSync(path.join(webappDir, 'i18n', 'en.json'), JSON.stringify(sortedWebappTranslations, null, 2) + '\n');
    });
}

export function i18nExtractMobile(argv) {
    const mobileDir = argv['mobile-dir'];

    const currentTranslations = getCurrentTranslationsMobile(mobileDir);
    const currentMobileKeys = new Set(Object.keys(currentTranslations));

    i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []).then((translationsMobile) => {
        const mobileKeys = new Set(Object.keys(translationsMobile));

        for (const key of difference(currentMobileKeys, mobileKeys)) {
            delete currentTranslations[key];
        }
        for (const key of difference(mobileKeys, currentMobileKeys)) {
            currentTranslations[key] = translationsMobile[key];
        }

        const options = {ignoreCase: true, reverse: false, depth: 1};
        const sortedMobileTranslations = sortJson(currentTranslations, options);
        fs.writeFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', 'en.json'), JSON.stringify(sortedMobileTranslations, null, 2) + '\n');
    });
}

export function i18nCombine(argv) {
    const outputFile = argv.output;

    const translations = {};

    for (const file of argv._.slice(2)) {
        const itemTranslationsJson = fs.readFileSync(file);
        const itemTranslations = JSON.parse(itemTranslationsJson);

        for (const key in itemTranslations) {
            if ({}.hasOwnProperty.call(itemTranslations, key)) {
                translations[key] = itemTranslations[key];
            }
        }
    }

    const options = {ignoreCase: true, reverse: false, depth: 1};
    const sortedTranslations = sortJson(translations, options);
    fs.writeFileSync(outputFile, JSON.stringify(sortedTranslations, null, 2) + '\n');
}

export function i18nSort(argv) {
    const outputFile = argv.output;

    const file = argv._[2];
    const itemTranslationsJson = fs.readFileSync(file);
    const itemTranslations = JSON.parse(itemTranslationsJson);

    const options = {ignoreCase: true, reverse: false, depth: 1};
    const sortedTranslations = sortJson(itemTranslations, options);
    fs.writeFileSync(outputFile, JSON.stringify(sortedTranslations, null, 2) + '\n');
}

export function i18nSplit(argv) {
    const webappDir = argv['webapp-dir'];
    const mobileDir = argv['mobile-dir'];
    const inputFiles = argv.inputs.split(',');

    const promise1 = i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']);
    const promise2 = i18nExtractLib.extractFromDirectory([argv['mobile-dir'] + '/app', argv['mobile-dir'] + '/share_extension'], []);
    Promise.all([promise1, promise2]).then(([translationsWebapp, translationsMobile]) => {
        for (const inputFile of inputFiles) {
            const filename = path.basename(inputFile.trim());
            const allTranslationsJson = fs.readFileSync(inputFile.trim());
            const allTranslations = JSON.parse(allTranslationsJson);

            const webappKeys = new Set(Object.keys(translationsWebapp));
            const mobileKeys = new Set(Object.keys(translationsMobile));

            const translationsWebappOutput = {};
            for (const key of webappKeys) {
                translationsWebappOutput[key] = allTranslations[key];
            }

            const translationsMobileOutput = {};
            for (const key of mobileKeys) {
                translationsMobileOutput[key] = allTranslations[key];
            }

            const options = {ignoreCase: true, reverse: false, depth: 1};
            const sortedWebappTranslations = sortJson(translationsWebappOutput, options);
            const sortedMobileTranslations = sortJson(translationsMobileOutput, options);
            fs.writeFileSync(path.join(webappDir, 'i18n', filename), JSON.stringify(sortedWebappTranslations, null, 2) + '\n');
            fs.writeFileSync(path.join(mobileDir, 'assets', 'base', 'i18n', filename), JSON.stringify(sortedMobileTranslations, null, 2) + '\n');
        }
    });
}

export function i18nClean(argv) {
    i18nCleanWebapp(argv);
    i18nCleanMobile(argv);
}

export function i18nCleanWebapp(argv) {
    const webappDir = argv['webapp-dir']
    const file = argv['file']
    const fPath = path.join(webappDir, 'i18n')
    const dryRun = argv['dry-run']
    const check = argv['check']
    const r = removeItems(fPath, file, dryRun)
    if (r !== '') {
        console.info(r)
    }
    if (check && r !== '') {
        return process.exit(1);
    }
}

export function i18nCleanMobile(argv) {
    const mobileDir = argv['mobile-dir']
    const file = argv['file']
    const fPath = path.join(mobileDir, 'assets', 'base', 'i18n')
    const dryRun = argv['dry-run']
    const check = argv['check']
    const r = removeItems(fPath, file, dryRun)
    if (r !== '') {
        console.info(r)
    }
    if (check && r !== '') {
        return process.exit(1);
    }
}

export function i18nCleanAll(argv) {
    i18nCleanAllWebapp(argv);
    i18nCleanAllMobile(argv);
}

export function i18nCleanAllWebapp(argv) {
    const webappDir = argv['webapp-dir']
    const dryRun = argv['dry-run']
    const check = argv['check']
    const fPath = path.join(webappDir, 'i18n')
    return cleanAll(fPath, dryRun, check);
}

export function i18nCleanAllMobile(argv) {
    const mobileDir = argv['mobile-dir']
    const dryRun = argv['dry-run']
    const check = argv['check']
    const fPath = path.join(mobileDir, 'assets', 'base', 'i18n')
    return cleanAll(fPath, dryRun, check);
}

function cleanAll(fPath, dryRun, check) {
    const files = fs.readdirSync(fPath)
    let rs = ''
    for (const f of files) {
        const r = removeItems(fPath, f, dryRun)
        rs += r
    }
    if (rs !== '') {
        console.info(fPath)
        console.info(rs)
    }
    if (check && rs !== '') {
        return process.exit(1);
    }
}

export function removeItems(fPath, f, dryRun) {
    if (f.split('.').pop() !== 'json' || f === 'en.json') {
        return ''
    }
    let count = 0
    const obj = JSON.parse(fs.readFileSync(path.join(fPath, f)).toString(), (k, v) => {
        if (v === '') {
            count++
        } else {
            return v
        }
    })

    if (count === 0) {
        return ''
    }
    if (!dryRun) {
        fs.writeFileSync(path.join(fPath, f), JSON.stringify(obj, null, 2) + '\n')
    }
    return f + ' has ' + count + ' empty translations'
}
