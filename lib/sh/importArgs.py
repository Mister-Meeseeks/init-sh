#!/usr/bin/env python

import sys

namespace = sys.argv[1]
isDropSuffix = int(sys.argv[2]) > 0 \
    if len(sys.argv) > 2 else False

# Nested imports are called with a namespace with trailing ::.
# I.e. "myLib" namespace imports ./myLib/subLib/script as myLib::script,
# whereas "myLib::" import as myLib::subLib::script
def isNamespaceNested (namespace) :
    return namespace[-2:] == "::"

def convertPathToFlatVar (importPath):
    return importPath.split("/")[-1]

def convertPathNestedVar (importPath):
    importDirs = importPath.split("/")
    return "::".join(importDirs)

def dropSuffixForImportVar (importVar):
    dotDivisions = importVar.split(".")
    return ".".join(dotDivisions[:-1]) \
        if len(dotDivisions) > 1 else importVar

def convertPathToVar (importPath, isNested, isDropSuffix):
    importCmd = convertPathNestedVar(importPath) \
        if isNested else convertPathToFlatVar(importPath)
    return dropSuffixForImportVar(importCmd) \
        if isDropSuffix else importCmd

def prependNamespace (importVar, namespace):
    sanitNamespace = namespace.rstrip(":")
    return sanitNamespace + "::" + importVar \
        if sanitNamespace != "" else importVar

def formImportPathTarget (importDir, importPath):
    return "%s/%s" % (importDir, importPath)

def printImportArgs (importTarget, importVar):
    print "%s %s" % (importTarget, importVar)

isNested = isNamespaceNested(namespace)
for line in sys.stdin:
    [importDir, importPath] = line.rstrip().split(" ")
    importVar = convertPathToVar(importPath, isNested, isDropSuffix)
    namespaceVar = prependNamespace(importVar, namespace)
    importTarget = formImportPathTarget(importDir, importPath)
    printImportArgs(importTarget, namespaceVar)
