package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"time"
)

type effs struct {
	absolutePathToRepo string
	pathToSourceRepo   string // use if you want to clone from file
	httpServer         string // use if you want to clone from server
	repoName           string
	fossilUser         string
	fossilPassword     string
	password           string
	entries            map[string]entry    // branch -> entry
	branches           map[string]string   // branch plaintext -> branch encrypted (from the repo)
	ordering           map[string][]string // document -> list of entries
}

type entry struct {
	date    string
	content string
}

// TODO: Add cleanup to close the repo

// initRepo should initialize a new repo (if no httpServer)
// otherwise it should clone/pull the httpServer repo
// then it should open and checkout the trunk
func (fs *effs) initRepo() {
	fs.run("fossil", "init", fs.repoName)
}

func (fs *effs) clone() {
	var stdout, stderr string
	if fs.httpServer != "" {
		fmt.Println("fossil", "clone", fs.httpServer+"/"+fs.repoName, fs.repoName)
		_, stdout, stderr = fs.run("fossil", "clone", fs.httpServer+"/"+fs.repoName, fs.repoName)
	} else if fs.pathToSourceRepo != "" {
		fmt.Println("fossil", "clone", fs.pathToSourceRepo, fs.repoName)
		_, stdout, stderr = fs.run("fossil.exe", "clone", fs.pathToSourceRepo, fs.repoName)
	}
	fmt.Println(stdout, stderr)
}

// push
func (fs effs) push() {
	var stdout, stderr string
	if len(fs.httpServer) != 0 {
		_, stdout, stderr = fs.run("fossil", "push", strings.Replace(fs.httpServer, "://", "://"+fs.fossilUser+":"+fs.fossilPassword+"@", -1)+"/"+fs.repoName, "-R", fs.repoName)
	} else if fs.pathToSourceRepo != "" {
		_, stdout, stderr = fs.run("fossil", "push", "-R", fs.repoName)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
}

// pull
func (fs effs) pull() {
	var stdout, stderr string
	if len(fs.httpServer) != 0 {
		_, stdout, stderr = fs.run("fossil", "pull", strings.Replace(fs.httpServer, "://", "://"+fs.fossilUser+":"+fs.fossilPassword+"@", -1)+"/"+fs.repoName, "-R", fs.repoName)
	} else if fs.pathToSourceRepo != "" {
		_, stdout, stderr = fs.run("fossil", "pull", "-R", fs.repoName)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	if strings.Contains(stderr, "a fork has occurred") {
		fs.handleMerge()
	}
}

func (fs effs) handleMerge() {
	var stdout, stderr string
	fs.run("fossil", "open", fs.repoName)
	fmt.Println("HANDLING MERGE!!!!!!!")
	_, stdout, stderr = fs.run("fossil", "leaves", "-multiple")
	fmt.Println(stdout, stderr)
	branchEncrypted := StrExtract(stdout, "tags: ", ")", 1)
	fmt.Println(branchEncrypted)
	fmt.Println("fossil", "checkout", branchEncrypted, "--force")
	_, stdout, stderr = fs.run("fossil", "checkout", branchEncrypted, "--force")
	fmt.Println(stdout, stderr)
	_, stdout, stderr = fs.run("fossil", "merge")
	fmt.Println(stdout, stderr)
	fmt.Println("--------")
	fmt.Println("--------")
	mergeData, _ := ioutil.ReadFile(path.Join(fs.absolutePathToRepo, "data.aes"))
	fmt.Println(string(mergeData))
	part1 := StrExtract(string(mergeData), "first <<<<<<<<<<<<<<<", "=======", 1)
	part2 := StrExtract(string(mergeData), "==================================", ">>>>>>>", 1)
	fmt.Println(part1, part2)
	part1D, _ := decryptString(strings.TrimSpace(part1), fs.password)
	part2D, _ := decryptString(strings.TrimSpace(part2), fs.password)
	fmt.Println("MERGED:")
	mergeText := strings.TrimSpace(MergeText(part1D, part2D))
	fmt.Println(mergeText)
	fs.run("fossil", "close")
	branch, _ := decryptString(branchEncrypted, fs.password)
	documentName := strings.Split(branch, "-==-")[0]
	entryName := strings.Split(branch, "-==-")[1]
	fs.editEntry(documentName, entryName, entry{"", mergeText})
	fs.push()
}

// getText returns the text of the branch in plaintext
// the branch in plaintext must be determined by using the timeline
// to load the map of plaintext branches -> encrypted branches
func (fs *effs) getText(branch string) (string, error) {
	var stdout, stderr string
	if len(fs.branches) == 0 {
		return "", errors.New("Must run timeline first")
	}
	if _, ok := fs.branches[branch]; !ok {
		return "", errors.New("Incorrect branch name")
	}
	fs.run("fossil", "open", fs.repoName)
	defer fs.run("fossil", "close")
	fmt.Println("fossil", "checkout", fs.branches[branch], "--force")
	_, stdout, stderr = fs.run("fossil", "checkout", fs.branches[branch], "--force")
	fmt.Println(stdout, stderr)
	contents, _ := openAndDecrypt(path.Join(fs.absolutePathToRepo, "data.aes"), fs.password)
	fmt.Println(contents)
	return contents, nil
}

// parseTimeline is used to gather all the files and determine the ordering / deleting (through commit messages)
func (fs *effs) parseTimeline() {
	if len(fs.entries) == 0 {
		fs.entries = make(map[string]entry)
	}
	fs.ordering = make(map[string][]string)
	fs.branches = make(map[string]string)
	fmt.Println("fossil", "timeline", "-n", "0", "-R", fs.repoName)
	_, stdout, stderr := fs.run("fossil", "timeline", "-n", "0", "-R", fs.repoName)
	fmt.Println(stdout, stderr)
	curDate := ""
	var entries []entry
	for _, line := range strings.Split(stdout, "\n") {
		if len(line) < 3 {
			continue
		}
		if line[0:3] == "===" {
			curDate = StrExtract(line, "=== ", " ===", 1)
			fmt.Println(curDate)
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if !(strings.Contains(fields[1], "[") && strings.Contains(fields[1], "]")) || strings.Contains(fields[2], "*") {
			continue
		}
		var e entry
		e.date = curDate + " " + fields[0]
		messageEncrypted := StrExtract(line, "] ", " (", 1)
		branchEncrypted := StrExtract(line, "tags: ", ")", 1)
		message, _ := decryptString(messageEncrypted, fs.password) // TODO: Do something if message = deleted
		branch, _ := decryptString(branchEncrypted, fs.password)
		fs.branches[branch] = branchEncrypted
		fmt.Println(line)
		fmt.Println(message, branch)
		if len(message) == 0 || len(branch) == 0 {
			continue
		}
		entries = append(entries, e)
		if _, ok := fs.entries[branch]; !ok {
			fs.entries[branch] = e
		} else {
			// TODO: Check if this is a newer version
			// if newer version, get the latest text
		}
		documentName := strings.Split(branch, "-==-")[0]
		if _, ok := fs.ordering[documentName]; !ok {
			fs.ordering[documentName] = []string{}
		}
		fs.ordering[documentName] = append(fs.ordering[documentName], branch)

	}
	fmt.Printf("%+v\n", entries)
	return
}

// addEntry adds a new entry to the document
// if no entryName is provided, a random one is provided
// if no date is provided, the current one is used
func (fs effs) addEntry(documentName string, entryName string, e entry) (err error) {
	fs.run("fossil", "open", fs.repoName)
	fs.run("fossil", "setting", "autosync", "off")
	defer fs.run("fossil", "close")
	if len(entryName) == 0 {
		entryName = GetRandomMD5Hash()
	}
	fmt.Println("encrypting")
	err = encryptAndWrite(path.Join(fs.absolutePathToRepo, "data.aes"), e.content, fs.password)
	// err = ioutil.WriteFile(path.Join(fs.absolutePathToRepo, "data.aes"), []byte(e.content), 0644)
	if err != nil {
		return
	}
	fmt.Println("fossil", "add", "data.aes")
	err, stdout, stderr := fs.run("fossil", "add", "--force", "data.aes")
	if err != nil {
		return err
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	branchName, _ := encryptString(documentName+"-==-"+entryName, fs.password)
	message, _ := encryptString("new", fs.password)
	fmt.Println("commiting")
	if len(e.date) > 0 {
		_, stdout, stderr = fs.run("fossil", "commit", "--branch", branchName, "-m", message, "--date-override", e.date, "--allow-older")
	} else {
		_, stdout, stderr = fs.run("fossil", "commit", "--branch", branchName, "-m", message)
	}
	for iterations := 0; iterations < 2; iterations++ {
		if len(stderr) > 0 {
			// There seems to be a problem with adding a file and then committing too quickly.
			// Take a break and try again
			time.Sleep(500 * time.Millisecond)
			_, stdout, stderr = fs.run("fossil", "commit", "--branch", branchName, "-m", message, "--date-override", e.date, "--allow-older")
		}
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	return nil
}

// editEntry adds a new entry to the document
// if no entryName is provided, a random one is provided
// if no date is provided, the current one is used
func (fs effs) editEntry(documentName string, entryName string, e entry) (err error) {
	fmt.Println("EDITING ENTRY")
	var stdout, stderr string
	fs.run("fossil", "open", fs.repoName)
	fs.run("fossil", "setting", "autosync", "off")
	defer fs.run("fossil", "close")
	branch := documentName + "-==-" + entryName
	if _, ok := fs.branches[branch]; !ok {
		return errors.New("Incorrect branch name")
	}
	branchName := fs.branches[branch]
	fmt.Println("fossil", "checkout", branchName, "--force")
	_, stdout, stderr = fs.run("fossil", "checkout", branchName, "--force")
	fmt.Println(stdout, stderr)
	err = encryptAndWrite(path.Join(fs.absolutePathToRepo, "data.aes"), e.content, fs.password)
	// err = ioutil.WriteFile(path.Join(fs.absolutePathToRepo, "data.aes"), []byte(e.content), 0644)
	if err != nil {
		return
	}
	fmt.Println("fossil", "add", "data.aes")
	err, stdout, stderr = fs.run("fossil", "add", "--force", "data.aes")
	if err != nil {
		return err
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	message, _ := encryptString("edited", fs.password)
	fmt.Println("commiting")
	if len(e.date) > 0 {
		_, stdout, stderr = fs.run("fossil", "commit", "--branch", branchName, "-m", message, "--date-override", e.date, "--allow-older")
	} else {
		_, stdout, stderr = fs.run("fossil", "commit", "--branch", branchName, "-m", message)
	}
	for iterations := 0; iterations < 2; iterations++ {
		if len(stderr) > 0 {
			// There seems to be a problem with adding a file and then committing too quickly.
			// Take a break and try again
			time.Sleep(500 * time.Millisecond)
			_, stdout, stderr = fs.run("fossil", "commit", "--branch", branchName, "-m", message, "--date-override", e.date, "--allow-older")
		}
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	return nil
}

// run is a function to run in a given directory and collect the output
func (effs effs) run(name string, args ...string) (error, string, string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = effs.absolutePathToRepo
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	cmd.Start()
	stderr, _ := ioutil.ReadAll(stderrPipe)
	stdout, err := ioutil.ReadAll(stdoutPipe)
	cmd.Wait()
	return err, string(stdout), string(stderr)
}
