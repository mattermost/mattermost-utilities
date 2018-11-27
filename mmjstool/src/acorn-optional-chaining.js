// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

const QUESTION_MARK_ASCII_CODE = 63;

export function acornOptionalChaining(AcornParser) {
    AcornParser.prototype.realReadWord1 = AcornParser.prototype.readWord1;
    AcornParser.prototype.realParseIdent = AcornParser.prototype.parseIdent;
    AcornParser.prototype.realParseExprAtom = AcornParser.prototype.parseExprAtom;

    function readWord1(...args) {
        const word = AcornParser.prototype.realReadWord1.apply(this, args);
        if (AcornParser.prototype.fullCharCodeAtPos.call(this) === QUESTION_MARK_ASCII_CODE) {
            ++this.pos;
        }
        return word;
    }
    function parseIdent(...args) {
        try {
            return AcornParser.prototype.realParseIdent.apply(this, args);
        } catch (e) {
            const startPos = this.start;
            const startLoc = this.startLoc;
            return this.parseSubscripts(this.parseExprAtom(), startPos, startLoc, true);
        }
    }

    AcornParser.prototype.readWord1 = readWord1;
    AcornParser.prototype.parseIdent = parseIdent;

    return AcornParser;
}

