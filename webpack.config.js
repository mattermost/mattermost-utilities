const path = require('path');

module.exports = {
    target: 'node',
    entry: './mmjstool/src/index.js',
    output: {
        path: path.resolve(__dirname, 'bin'),
        filename: 'mmjstool',
    },
};
