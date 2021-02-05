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

    i18nExtractLib.extractFromDirectory([argv['webapp-dir']], ['storybook-static', 'dist', 'node_modules', 'non_npm_dependencies', 'tests', 'components/gif_picker/static/gif.worker.js']).then((translationsWebapp) => {
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

export function i18nCheckEmptySrc(argv) {
    const wCode = i18nCheckEmptySrcWebapp(argv);
    const mCode = i18nCheckEmptySrcMobile(argv);
    process.exit(wCode || mCode);
}

export function i18nCheckEmptySrcWebapp(argv) {
    const webappDir = argv['webapp-dir'];
    const fPath = path.join(webappDir, 'i18n');
    const counter = countEmptyItems(fPath, 'en.json');
    if (counter > 0) {
        const msg = 'Found ' + counter + ' empty translations in ' + fPath + '/en.json\n';
        console.info(msg);
        return 1;
    }
    return 0;
}

export function i18nCheckEmptySrcMobile(argv) {
    const mobileDir = argv['mobile-dir'];
    const fPath = path.join(mobileDir, 'assets', 'base', 'i18n');
    const counter = countEmptyItems(fPath, 'en.json');
    if (counter > 0) {
        const msg = 'Found ' + counter + ' empty translations in ' + fPath + '/en.json\n';
        console.info(msg);
        return 1;
    }
    return 0;
}

function countEmptyItems(filePath, file) {
    let count = 0;
    JSON.parse(fs.readFileSync(path.join(filePath, file)).toString(), (k, v) => {
        if (v === '') {
            count++;
        }
        return v;
    });
    return count;
}

export function i18nCleanEmpty(argv) {
    const wCode = i18nCleanEmptyWebapp(argv);
    const mCode = i18nCleanEmptyMobile(argv);
    process.exit(wCode || mCode);
}

export function i18nCleanEmptyWebapp(argv) {
    const webappDir = argv['webapp-dir'];
    const dryRun = argv['dry-run'];
    const check = argv.check;
    const fPath = path.join(webappDir, 'i18n');
    return cleanAll(fPath, dryRun, check);
}

export function i18nCleanEmptyMobile(argv) {
    const mobileDir = argv['mobile-dir'];
    const dryRun = argv['dry-run'];
    const check = argv.check;
    const fPath = path.join(mobileDir, 'assets', 'base', 'i18n');
    return cleanAll(fPath, dryRun, check);
}

function cleanAll(filePath, dryRun, check) {
    const files = fs.readdirSync(filePath);
    let results = '';
    for (const file of files) {
        if (file.split('.').pop() !== 'json' || file === 'en.json') {
            continue;
        }
        const result = removeItems(filePath, file, dryRun, check);
        results += result;
    }
    if (results === '') {
        return 0;
    }
    console.info(filePath);
    console.info(results);
    if (check) {
        return 1;
    }
    return 0;
}

export function removeItems(filePath, file, dryRun, check) {
    let count = 0;
    const obj = JSON.parse(fs.readFileSync(path.join(filePath, file)).toString(), (k, v) => {
        if (v === '') {
            count++;
            return undefined; // eslint-disable-line no-undefined
        }
        return v;
    });

    if (count === 0) {
        return '';
    }
    const msg = file + ' has ' + count + ' empty translations\n';
    if (dryRun || check) {
        return msg;
    }
    fs.writeFileSync(path.join(filePath, file), JSON.stringify(obj, null, 2) + '\n');
    return msg;
}
