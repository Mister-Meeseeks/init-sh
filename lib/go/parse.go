
package initsh

import "strings"
import "errors"

func parseImportStr (cmd string, dir ImportDirector) (*PathIngester, string, error) {
	lex, err := lexImportStr(cmd)
	if err != nil {
		return nil, "", err
	}
	ing, err := parseImportLex(lex, dir)
	return ing, lex.importPath, err
}

type importCmdLex struct {
	cmdType string
	importPath string
	namespace *string
}

func lexImportStr (cmd string) (importCmdLex, error) {
	fields := strings.Split(cmd, ":")
	switch len(fields) {
	case 1:
		return raiseNoImportTgt(cmd)
	case 2:
		return importCmdLex{fields[0], fields[1], nil}, nil
	case 3:
		return importCmdLex{fields[0], fields[1], &(fields[2])}, nil
	default:
		return raiseMalformedImport(cmd)
	}
}

func raiseNoImportTgt (cmd string) (importCmdLex, error) {
	empty := importCmdLex{"", "", nil}
	err := errors.New("initSh parse error: No import target specified in " +
		"directive=" + cmd)
	return empty, err
}

func raiseMalformedImport (cmd string) (importCmdLex, error) {
	empty := importCmdLex{"", "", nil}
	err := errors.New("initSh parse error: Malformed import directive=" + cmd)
	return empty, err
}

func parseImportLex (cmd importCmdLex, dir ImportDirector) (*PathIngester, error) {
	switch cmd.cmdType {
	case "shell":
		return liftDir(dir.importShell(cmd.importPath, cmd.namespace))
	case "nested":
		return liftDir(dir.importNested(cmd.importPath, cmd.namespace))
	case "subcmd":
		if cmd.namespace == nil {
			return nil, raiseSubcmdNamespace(cmd)
		}
		return liftDir(dir.importSubcmd(cmd.importPath, *(cmd.namespace)))
	case "subcmdNest":
		if cmd.namespace == nil {
			return nil, raiseSubcmdNamespace(cmd)
		}
		return liftDir(dir.importNestSubcmd(cmd.importPath, *(cmd.namespace)))
	case "data":
		return liftDir(dir.importData(cmd.importPath, cmd.namespace))
	case "dataNest":
		return liftDir(dir.importNestData(cmd.importPath, cmd.namespace))
	default:
		return nil, raiseBadCmdType(cmd)
	}
}

func liftDir (ing PathIngester) (*PathIngester, error) {
	return &ing, nil
}

func raiseBadCmdType (cmd importCmdLex) error {
	return errors.New("Unrecognized directive command=" + cmd.cmdType +
		" on import target=" + cmd.importPath)
}

func raiseSubcmdNamespace (cmd importCmdLex) error {
	return errors.New("initSh directive error: subcmd imports must specify " +
		"a namespace. import target=" + cmd.importPath)
}
