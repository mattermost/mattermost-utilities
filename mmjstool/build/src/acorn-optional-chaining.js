"use strict";

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.acornOptionalChaining = acornOptionalChaining;
// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

var QUESTION_MARK_ASCII_CODE = 63;

function acornOptionalChaining(AcornParser) {
    AcornParser.prototype.realReadWord1 = AcornParser.prototype.readWord1;
    AcornParser.prototype.realParseIdent = AcornParser.prototype.parseIdent;
    AcornParser.prototype.realParseExprAtom = AcornParser.prototype.parseExprAtom;

    function readWord1() {
        for (var _len = arguments.length, args = Array(_len), _key = 0; _key < _len; _key++) {
            args[_key] = arguments[_key];
        }

        var word = AcornParser.prototype.realReadWord1.apply(this, args);
        if (AcornParser.prototype.fullCharCodeAtPos.call(this) === QUESTION_MARK_ASCII_CODE) {
            ++this.pos;
        }
        return word;
    }
    function parseIdent() {
        try {
            for (var _len2 = arguments.length, args = Array(_len2), _key2 = 0; _key2 < _len2; _key2++) {
                args[_key2] = arguments[_key2];
            }

            return AcornParser.prototype.realParseIdent.apply(this, args);
        } catch (e) {
            var startPos = this.start;
            var startLoc = this.startLoc;
            return this.parseSubscripts(this.parseExprAtom(), startPos, startLoc, true);
        }
    }

    AcornParser.prototype.readWord1 = readWord1;
    AcornParser.prototype.parseIdent = parseIdent;

    return AcornParser;
}