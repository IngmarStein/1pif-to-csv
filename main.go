package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

var (
	inputFile    = flag.String("input", "", "input file in 1pif format")
	outputFile   = flag.String("output", "", "output CSV file")
	printVersion = flag.Bool("version", false, "print version")
	version      = "0.3"
)

type onepifRow struct {
	Uuid           string `json:"uuid"`
	UpdatedAt      int    `json:"updatedAt"`
	LocationKey    string `json:"locationKey"`
	SecurityLevel  string `json:"securityLevel"`
	ContentsHash   string `json:"contentsHash"`
	Title          string `json:"title"`
	Location       string `json:"location"`
	SecureContents struct {
		Fields []struct {
			Value       string `json:"value"`
			Name        string `json:"name"`
			Type        string `json:"type"`
			Id          string `json:"id,omitempty"`
			Designation string `json:"designation,omitempty"`
		} `json:"fields"`
		PasswordHistory []struct {
			Value string `json:"value"`
			Time  int    `json:"time"`
		} `json:"passwordHistory"`
		NotesPlain string `json:"notesPlain"`
		HtmlMethod string `json:"htmlMethod"`
		Sections   []struct {
			Fields []struct {
				K string          `json:"k"`
				N string          `json:"n"`
				V json.RawMessage `json:"v"`
				T string          `json:"t"`
			} `json:"fields,omitempty"`
			Name  string `json:"name"`
			Title string `json:"title,omitempty"`
		} `json:"sections"`
		URLs []struct {
			Url string `json:"url"`
		} `json:"URLs"`
		Password string `json:"password"`
	} `json:"secureContents"`
	TxTimestamp int    `json:"txTimestamp"`
	CreatedAt   int    `json:"createdAt"`
	TypeName    string `json:"typeName"`
}

func (r *onepifRow) fieldWithDesignation(designation string) (string, error) {
	for _, field := range r.SecureContents.Fields {
		if field.Designation == designation {
			return field.Value, nil
		}
	}

	return "", fmt.Errorf("no %v found", designation)
}

func (r *onepifRow) username() (string, error) {
	return r.fieldWithDesignation("username")
}

func (r *onepifRow) password() (string, error) {
	if r.SecureContents.Password != "" {
		return r.SecureContents.Password, nil
	}
	return r.fieldWithDesignation("password")
}

func (r *onepifRow) otpAuth() (string, error) {
	for _, section := range r.SecureContents.Sections {
		for _, field := range section.Fields {
			var v string
			err := json.Unmarshal(field.V, &v)
			if err != nil {
				continue
			}
			if strings.HasPrefix(v, "otpauth://") {
				// transform "otpauth://" into "apple-otpauth://"
				return "apple-" + v, nil
			}
		}
	}

	return "", errors.New("no otpauth found")
}

// openInputFile attempts to open a given 1pif file by path name. If the path points
// to a directory, `${name}/data.1pif` is used instead.
func openInputFile(name string) (*os.File, error) {
	ifile, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := ifile.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return ifile, nil
	}

	ifile.Close()
	ifile, err = os.Open(path.Join(name, "data.1pif"))
	if err != nil {
		return nil, err
	}

	return ifile, nil
}

func onepifToCSV(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	csvWriter := csv.NewWriter(out)

	// write header
	csvWriter.Write([]string{"Title", "Url", "Username", "Password", "OTPAuth"})

	lineNr := 0
	for scanner.Scan() {
		lineNr++
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "***") && strings.HasSuffix(line, "***") {
			continue
		}

		var row onepifRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			log.Printf("failed to parse line %d: %v %q\n", lineNr, err, line)
			continue
		}

		if len(row.SecureContents.URLs) == 0 {
			continue
		}

		username, _ := row.username()
		password, _ := row.password()
		otpauth, _ := row.otpAuth()

		seen := make(map[string]bool)
		for _, rowURL := range row.SecureContents.URLs {
			u, err := url.Parse(rowURL.Url)
			if err != nil {
				log.Printf("failed to parse URL %q (line %d): %v\n", rowURL.Url, lineNr, err)
				continue
			}
			// assume https if no scheme is set
			if u.Scheme == "" {
				u, err = url.Parse("https://" + rowURL.Url)
				if err != nil {
					log.Printf("failed to parse URL %q (line %d): %v\n", rowURL.Url, lineNr, err)
					continue
				}
			}
			// Monterey's password manager discards the URL path
			u.Path = "/"
			u.RawQuery = ""
			u.Fragment = ""
			domain := u.String()
			if seen[domain] {
				// skip domains we've already emitted
				continue
			}
			seen[domain] = true
			record := []string{
				row.Title,
				domain,
				username,
				password,
				otpauth,
			}

			csvWriter.Write(record)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	csvWriter.Flush()

	return nil
}

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	ifile, err := openInputFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer ifile.Close()

	ofile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer ofile.Close()

	if err := onepifToCSV(ifile, ofile); err != nil {
		log.Fatal(err)
	}
}
