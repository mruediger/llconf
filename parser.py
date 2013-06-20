#!/bin/python

import sys


class Promise(object):

    def __init__(self):
        self.bracket_sum = 0
        self.line = ""
        self.name = ""
        self.accepts_name = True

    def add(self, char):
        if (self.accepts_name):
            if (char.isspace() or char in {'(',')'} ):
                self.accepts_name = False
            else:
                self.name += char
        
        if (char != '\n'):
            self.line += char

        if ( char == '('):
            if ( not self.name ):
                self.accepts_name = True
            self.bracket_sum += 1
            return

        if ( char == ')'):
            self.bracket_sum -= 1

    def isEnded(self):
        return self.bracket_sum == 0



    def __str__(self):
        return self.name + ':' + self.line.strip()


def consume(line, rules):
    
    word = str()
    rule = Promise()

    for c in line:
        rule.add(c)
        if ( c == ")" and rule.isEnded()):
            rules.append(rule)
            rule = Promise()

    return rules    
                


if __name__ == "__main__":
    
    rules = []
    for line in (sys.stdin):
        rules = consume(line, rules)

    for rule in rules:
        print rule
