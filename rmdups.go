package main

import (
	"github.com/carlogit/phash"
	"crypto/md5"
	"encoding/gob"
	"path/filepath"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type MetaFile struct {
	Path	string
	Size	int64
	Hash	[]byte
	HashImg	string
}
type MetaFileList map[string]MetaFile

type Mode int
const (
	ModeInit  Mode = iota
	ModeDedup      = iota
)

var (
	flagHashSize = flag.Int64("hashsize", 1024*1024, "count mask (verbose logging every N items")
	flagVerbose  = flag.Bool("verbose", false, "verbose logging")
	flagCount    = flag.Int("count", 1000, "count mask (verbose logging every N items")
	flagSilos    = flag.String("silos", ".silos", "silos path for md5 sums")
	flagInit     = flag.Bool("init", false, "init silos")
	flagDedup    = flag.Bool("dedup", false, "dedup stuff to silos")
	flagFirst    = flag.Bool("1", false, "if duplicate, operate 1st file")
	flagSecond   = flag.Bool("2", false, "if duplicate, operate 2nd file")
	flagRemove   = flag.Bool("remove", false, "remove file")
	flagImg      = flag.Bool("img", false, "operate on image")
)

func main() {
	flag.Parse()

	if (len(flag.Args()) != 1) {
		flag.Usage()
		log.Fatal("bleh")
	}
	searchDir := flag.Args()[0]
	if (*flagImg) {
		*flagSilos += ".img"
	}
	fmt.Println("==>", searchDir, *flagSilos, *flagVerbose, *flagCount)

	if (*flagInit == false && *flagDedup == false) {
		flag.Usage()
	}

	fmt.Println("# indexing....")
	rawPathList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if (f.IsDir()) {
			return nil
		}
		rawPathList = append(rawPathList, path)
		if (len(rawPathList) % *flagCount == 0) {
			fmt.Printf("# %d files indexed..\n", len(rawPathList))
		}
		return nil
	})
	if (err != nil) {
		log.Fatal(err)
	}

	if (*flagInit) {
		silosOp(*flagSilos, rawPathList, ModeInit)
	} else if (*flagDedup) {
		silosOp(*flagSilos, rawPathList, ModeDedup)
	}
}

func silosOp(silosFile string, pathList []string, whichMode Mode) {
	var metaFileList = make(MetaFileList)
	if (whichMode == ModeDedup) {
		metaFileList = silosRead(silosFile)
		fmt.Println("# will dedup now")
	}

	for path_i, path := range(pathList) {
		silosFileIterator(path, &metaFileList)
		if ((path_i % *flagCount) == 0) {
			fmt.Printf("# %d/%d hashing done %s\n", path_i, len(pathList), path)
		}
	}

	silosWrite(silosFile, metaFileList)
}

func silosFileIterator(path string, ma *MetaFileList) error {
	var metaFile, _ = processFile(path)
	var key = fmt.Sprintf("%x", metaFile.Hash)
	if (metaFile.HashImg != "") {
		key = metaFile.HashImg
	}

	maybeDupFile, key_exists := (*ma)[key]
	if (key_exists && (maybeDupFile.Path == path)) {
		return nil
	}

	if (key_exists) {
		fmt.Printf("EXISTS: %s!\n", path)
		fmt.Printf("      : %s\n", maybeDupFile.Path)
		handleDuplicate(path, maybeDupFile.Path)
	} else {
		(*ma)[key] = *metaFile
		if (*flagVerbose) {
			fmt.Println(path)
		}
	}
	return nil
}

func handleDuplicate(path1 string, path2 string) {
	path := path1
	if (*flagSecond) {
		path = path2
	}
	if (*flagRemove) {
		fmt.Printf("rm %s\n", path)
	}
}

func processFile(fileName string) (*MetaFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	md5hash := md5.New()
	lenToHash := stat.Size()
	if (lenToHash > *flagHashSize) {
		lenToHash = *flagHashSize
	}
	io.CopyN(md5hash, file, lenToHash)

	var fileNameExt = fileName[len(fileName) - 3:]

	var hashImg = ""
	if (fileNameExt == "jpg" || fileNameExt == "JPG") {
		file.Seek(0, os.SEEK_SET)
		hashImg, err = phash.GetHash(file)
		if (err != nil) {
			fmt.Println(fileName, " is messed up. will skip")
			hashImg = ""
		}
	}

	file.Close()

	return &MetaFile{
		Path: fileName,
		Size: stat.Size(),
		Hash: md5hash.Sum(nil),
		HashImg: hashImg,
	}, nil
}

func silosWrite(silosFile string, metaFileList MetaFileList) {
	fmt.Println("writing")
	encodeFile, err := os.Create(silosFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("writing")
	encoder := gob.NewEncoder(encodeFile)
	if err := encoder.Encode(metaFileList); err != nil {
		panic(err)
	}
	fmt.Println("writing")
	encodeFile.Close()
}

func silosRead(silosFile string) MetaFileList {
	decodeFile, err := os.Open(silosFile)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	metaFileList := make(map[string]MetaFile)
	decoder.Decode(&metaFileList)
	return metaFileList
}
