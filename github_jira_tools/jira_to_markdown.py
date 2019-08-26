# Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
# See LICENSE.txt for license information.

# Based on the jira2md npm pacakge

import re

def ordered_list(match):
    stars = len(match.group(1))
    return ("  " * stars) + '* ';

def unordered_list(match):
    nums = len(match.group(1))
    return ("  " * nums) + '1. ';

def headers(match):
    level = int(match.group(1), 10) + 1
    content = match.group(2)
    return ("#" * level) + content;

def table_headers(match):
    headers = match.group(1)
    singleBarred = re.sub(r'\|\|', '|', headers)
    return '\n' + singleBarred + '\n' + re.sub(r'\|[^|]+', '| --- ', singleBarred)

def jira_to_markdown(string):
    if string is None:
        return ''
    # Ordered Lists
    result = re.sub(r'^[ \t]*(\*+)\s+', ordered_list, string, flags=re.MULTILINE)
    # Un-ordered lists
    result = re.sub(r'^[ \t]*(#+)\s+', unordered_list, result, flags=re.MULTILINE)
    # Headers 1-6
    result = re.sub(r'^h([0-6])\.(.*)$', headers, result, flags=re.MULTILINE)
    # table header
    result = re.sub(r'^[ \t]*((?:\|\|.*?)+\|\|)[ \t]*$', table_headers, result, flags=re.MULTILINE)
    # Code Block
    result = re.sub(r'\{code(:([a-z]+))?([:|]?(title|borderStyle|borderColor|borderWidth|bgColor|titleBGColor)=.+?)*\}(.*?)\{code\}', r'```\2\n\5\n```', result, flags=re.MULTILINE|re.DOTALL)
    # Bold
    result = re.sub(r'\*(\S.*)\*', r'**\1**', result)
    # Italic
    result = re.sub(r'\_(\S.*)\_', r'*\1*', result)
    # Monospaced text
    result = re.sub(r'\{\{([^}]+)\}\}', r'`\1`', result)
    # Citations (buggy)
    # result = re.sub(r'\?\?((?:.[^?]|[^?].)+)\?\?', r'<cite>\1</cite>', result)
    # Inserts
    result = re.sub(r'\+([^+]*)\+', r'<ins>\0</ins>', result)
    # Superscript
    result = re.sub(r'\^([^^]*)\^', r'<sup>\1</sup>', result)
    # Subscript
    result = re.sub(r'~([^~]*)~', r'<sub>\1</sub>', result)
    # Strikethrough
    result = re.sub(r'(\s+)-(\S+.*?\S)-(\s+)', r'\1~~\2~~\3', result)
    # Pre-formatted text
    result = re.sub(r'{noformat}', r'```', result)
    # Un-named Links
    result = re.sub(r'\[([^|]+)\]', r'<\1>', result)
    # Images
    result = re.sub(r'!(.+)!', r'![](\1)', result)
    # Named Links
    result = re.sub(r'\[(.+?)\|(.+)\]', r'[\1](\2)', result)
    # Single Paragraph Blockquote
    result = re.sub(r'^bq\.\s+', r'> ', result, flags=re.MULTILINE)
    # Remove color: unsupported in md
    result = re.sub(r'\{color:[^}]+\}(.*)\{color\}', r'\1', result, flags=re.MULTILINE)
    # panel into table
    result = re.sub(r'\{panel:title=([^}]*)\}\n?(.*?)\n?\{panel\}', r'\n| \1 |\n| --- |\n| \2 |', result, flags=re.MULTILINE)
    # remove leading-space of table headers and rows
    result = re.sub(r'^[ \t]*\|', r'|', result, flags=re.MULTILINE);
    return result
