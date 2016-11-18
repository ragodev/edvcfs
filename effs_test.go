package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestEffs(t *testing.T) {
	wd, _ := os.Getwd()
	var fs effs
	fs.absolutePathToRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123"))
	fs.repoName = "repo"
	// fs.httpServer = "http://cowyo.com:8082"
	fs.fossilUser = "zns"
	fs.fossilPassword = "80b2af"
	os.MkdirAll(fs.absolutePathToRepo, 0755)
	fs.password = "testtest"
	fs.initRepo()
	if !exists(path.Join(fs.absolutePathToRepo, fs.repoName)) {
		t.Errorf("Initiation from file didn't work")
	}
	fs.addEntry("test2.txt", "6", entry{"2010-08-19 07:21:21", "first entry, number 6"})
	fs.addEntry("test2.txt", "5", entry{"2010-08-19 06:21:21", "sixth entry"})
	fs.addEntry("test2.txt", "4", entry{"2010-08-18 06:21:21", "fifth entry some new text"})
	fs.addEntry("test.txt", "", entry{"2009-08-18 07:21:21", "fourth entry"})
	fs.addEntry("test.txt", "3", entry{"2009-08-18 06:21:21", "third entry"})
	fs.addEntry("test.txt", "2", entry{"2008-08-18 07:21:21", "second entry"})
	fs.addEntry("test.txt", "1", entry{"2007-08-18 07:21:21", "first entry"})
	fs.parseTimeline()
	text, _ := fs.getText("test2.txt-==-6")
	if text != "first entry, number 6" {
		t.Errorf(text)
	}
	if fmt.Sprintf("%v", fs.ordering["test2.txt"]) != "[test2.txt-==-6 test2.txt-==-5 test2.txt-==-4]" {
		t.Errorf("'%s'", fmt.Sprintf("%v", fs.ordering["test2.txt"]))
	}
	if fmt.Sprintf("%v", fs.ordering["test.txt"][1:]) != "[test.txt-==-3 test.txt-==-2 test.txt-==-1]" {
		t.Errorf("'%s'", fmt.Sprintf("%v", fs.ordering["test.txt"][1:]))
	}

	// test edit entry
	fs.editEntry("test.txt", "2", entry{"", "second entry-edited"})
	fs.parseTimeline()
	text, _ = fs.getText("test.txt-==-2")
	if text != "second entry-edited" {
		t.Errorf(text)
	}
	// Test cloning
	var fs2 effs
	fs2.absolutePathToRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123-2"))
	fs2.repoName = "repo"
	fs2.pathToSourceRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123", "repo"))
	fs2.fossilUser = "zns"
	fs2.fossilPassword = "80b2af"
	os.MkdirAll(fs2.absolutePathToRepo, 0755)
	fs2.password = "testtest"
	fs2.clone()
	if !exists(path.Join(fs2.absolutePathToRepo, fs2.repoName)) {
		t.Errorf("Didnt get created")
	}

	var fs3 effs
	fs3.absolutePathToRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123-3"))
	fs3.repoName = "repo"
	fs3.pathToSourceRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123", "repo"))
	fs3.fossilUser = "zns"
	fs3.fossilPassword = "80b2af"
	os.MkdirAll(fs3.absolutePathToRepo, 0755)
	fs3.password = "testtest"
	fs3.clone()
	if !exists(path.Join(fs3.absolutePathToRepo, fs3.repoName)) {
		t.Errorf("Didnt get created")
	}

	// test the push
	fs3.addEntry("test2.txt", "7", entry{"2012-08-20 07:21:21", "first entry, number 7"})
	fs3.push()
	fs.parseTimeline()
	text, _ = fs.getText("test2.txt-==-7")
	if text != "first entry, number 7" {
		t.Errorf(text)
	}

	// test the pull
	fs2.pull()
	fs2.parseTimeline()
	text, _ = fs2.getText("test2.txt-==-7")
	if text != "first entry, number 7" {
		t.Errorf(text)
	}

	// test the push even though its not up to date
	// everything should still be in sync
	fs3.addEntry("test.txt", "8", entry{"2013-08-20 07:21:21", "number 8"})
	fs3.push()
	fs.parseTimeline()
	text, _ = fs.getText("test.txt-==-8")
	if text != "number 8" {
		t.Errorf(text)
	}
	fs2.addEntry("test.txt", "9", entry{"2013-07-20 07:21:21", "number 9"})
	fs2.push()
	fs2.pull()
	fs2.parseTimeline()
	text, _ = fs2.getText("test.txt-==-8")
	if text != "number 8" {
		t.Errorf(text)
	}
	fs.parseTimeline()
	text, _ = fs.getText("test.txt-==-9")
	if text != "number 9" {
		t.Errorf(text)
	}

	fmt.Println(wd)
	fmt.Println(fs.absolutePathToRepo)
	fmt.Println(fs2.absolutePathToRepo)
	fmt.Println(fs3.absolutePathToRepo)
	os.RemoveAll(fs.absolutePathToRepo)
	os.RemoveAll(fs2.absolutePathToRepo)
	os.RemoveAll(fs3.absolutePathToRepo)
}

func TestMerge(t *testing.T) {
	ENCRYPTION_ENABLED = true
	wd, _ := os.Getwd()
	var fs effs
	fs.absolutePathToRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123"))
	fs.repoName = "repo"
	// fs.httpServer = "http://cowyo.com:8082"
	fs.fossilUser = "zns"
	fs.fossilPassword = "80b2af"
	os.MkdirAll(fs.absolutePathToRepo, 0755)
	fs.password = "testtest"
	fs.initRepo()
	if !exists(path.Join(fs.absolutePathToRepo, fs.repoName)) {
		t.Errorf("Initiation from file didn't work")
	}
	fs.addEntry("test.txt", "3", entry{"2008-08-18 07:21:21", "second entry"})
	fs.addEntry("test.txt", "1", entry{"2007-08-18 07:21:21", "first entry"})
	fs.parseTimeline()

	// Test cloning
	var fs2 effs
	fs2.absolutePathToRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123-2"))
	fs2.repoName = "repo"
	fs2.pathToSourceRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123", "repo"))
	fs2.fossilUser = "zns"
	fs2.fossilPassword = "80b2af"
	os.MkdirAll(fs2.absolutePathToRepo, 0755)
	fs2.password = "testtest"
	fs2.clone()
	fs2.parseTimeline()
	if !exists(path.Join(fs2.absolutePathToRepo, fs2.repoName)) {
		t.Errorf("Didnt get created")
	}

	var fs3 effs
	fs3.absolutePathToRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123-3"))
	fs3.repoName = "repo"
	fs3.pathToSourceRepo, _ = filepath.Abs(path.Join(wd, "..", "..", "..", "temp123", "repo"))
	fs3.fossilUser = "zns"
	fs3.fossilPassword = "80b2af"
	os.MkdirAll(fs3.absolutePathToRepo, 0755)
	fs3.password = "testtest"
	fs3.clone()
	fs3.parseTimeline()
	if !exists(path.Join(fs3.absolutePathToRepo, fs3.repoName)) {
		t.Errorf("Didnt get created")
	}

	// test editing entries simultaneously
	fs2.editEntry("test.txt", "3", entry{"", "second entry\nbottomedit2"})
	fs2.push()
	fs3.editEntry("test.txt", "3", entry{"", "topedit3\nsecond entry"})
	fs3.push()
	fs3.parseTimeline()
	text, _ := fs.getText("test.txt-==-3")
	if text != "topedit3\nsecond entry" {
		t.Errorf(text)
	}

	// handle the merge
	fs2.parseTimeline()
	fs2.pull()
	fs.parseTimeline()
	text, _ = fs.getText("test.txt-==-3")
	if text != "topedit3\nsecond entry\nbottomedit2" {
		t.Errorf("'%s'", text)
	}
	os.RemoveAll(fs.absolutePathToRepo)
	os.RemoveAll(fs2.absolutePathToRepo)
	os.RemoveAll(fs3.absolutePathToRepo)

}
