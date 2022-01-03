// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

const path = require('path');

const webpack = require('webpack');

module.exports = {
    target: 'node',
    entry: './mmjstool/src/index.js',
    mode: 'production',
    output: {
        path: path.resolve(__dirname, 'bin'),
        filename: 'mmjstool',
    },
    plugins: [
        new webpack.BannerPlugin({banner: '#!/usr/bin/env node', raw: true}),
    ],
};
