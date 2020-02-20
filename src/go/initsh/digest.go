
package initsh

import "os"
import "bufio"
import "sort"
import "bytes"
import "crypto/sha256"
import "encoding/base32"

type DigestLookup struct {
	isMatch bool
	space string
}

type DigestPreImage struct {
	directive string
	space string
	contents []string
}

type DigestEntry struct {
	checkSum string
	spaceName string
}

func DigestTarget (directive string, paths []string, lobby ImportLobby) DigestLookup {
	digestKey := deriveKey(directive)
	checkSum := hashContents(paths)
	return toggleSpace(digestKey, checkSum, lobby)
}

func WriteDigest (image DigestPreImage, lobby ImportLobby) error {
	digestKey := deriveKey(image.directive)
	checkSum := hashContents(image.contents)
	entry := DigestEntry{checkSum, image.space}
	
	path, err := lobby.initDigest(digestKey)
	if (err != nil) {
		return err
	}
	return writeEntry(path, entry)
}

func deriveKey (directive string) string {
	return hashString(directive)
}

func hashString (value string) string {
	content := []string{value}
	return hashContents(content)
}

func hashContents (contents []string) string {
	preImage := collectBytes(contents)
	checkSum := sha256.Sum256(preImage)
	return base32.StdEncoding.EncodeToString(checkSum[:])
}

func collectBytes (contents []string) []byte {
	sort.Strings(contents)
	var buffer bytes.Buffer
	for _, row := range contents {
		buffer.WriteString(row)
	}
	return buffer.Bytes()
}

func toggleSpace (digestKey string, checkSum string, lobby ImportLobby) DigestLookup {
	path, tryErr := lobby.pullDigest(digestKey)
	if (tryErr != nil) {
		return DigestLookup{false, ""}
	} else {
		return comparePrior(path, checkSum)
	}
}

func comparePrior (path string, checkSum string) DigestLookup {
	entry, err := readDigest(path)
	if (err != nil) {
		return DigestLookup{false, ""}
	} else if (checkSum != entry.checkSum) {
		return DigestLookup{false, ""}
	} else {
		return DigestLookup{true, entry.spaceName}
	}
}

func readDigest (path string) (DigestEntry, error) {
	if (doesPathExist(path)) {
		return readExtant(path)
	} else {
		return emptyEntry(), nil
	}
}

func readExtant (stampPath string) (DigestEntry, error) {
	input, err := os.Open(stampPath)
	if (err != nil) {
		return emptyEntry(), err
	}
	scanner := bufio.NewScanner(input)
	return readScanner(scanner)
}

func readScanner (scanner *bufio.Scanner) (DigestEntry, error) {
	if (scanner.Scan()) {
		spaceName := scanner.Text()
		if (scanner.Scan()) {
			checkSum := scanner.Text()
			return DigestEntry{checkSum, spaceName}, nil
		}
	}
	return emptyEntry(), nil
}


func writeEntry (path string, entry DigestEntry) error {
	file, err := os.Create(path)
	defer file.Close()
	if (err != nil) {
		return err
	}
	buffer := bufio.NewWriter(file)
	return writeOn(buffer, entry)
}

func writeOn (writer *bufio.Writer, entry DigestEntry) error {
	_, err := writer.WriteString(entry.spaceName)
	if (err != nil) {
		return err
	}
	_, err = writer.WriteString(entry.checkSum)
	return err
}

func emptyEntry() DigestEntry {
	return DigestEntry{"", ""}
}
