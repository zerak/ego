package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
type Section struct {
	data         map[string]string // key:value
	dataOrder    []string
	dataComments map[string][]string // key:comments
	Name         string
	comments     []string
	Comment      string
}
*/

const (
	CRLF     = '\n'
	Delimit  = ","
	Split    = " "
	Comment  = "#"
	SectionB = "["
	SectionE = "]"
	Include  = "include"
)

/*
# include other file
include ../common

# comment
commonKey commonVal
commonKey commonVal1,commonVal2

[sector]
sectorKey sectorVal1,sectorVal2
Section
*/
type Section struct {
	delimit string
	sector  string
	val     map[string]string // key val1,val2
}

// An NoKeyError describes a key that was not found in the section.
type NoKeyError struct {
	Key     string
	Section string
}

func (e *NoKeyError) Error() string {
	return fmt.Sprintf("key: \"%s\" not found in [%s]", e.Key, e.Section)
}

// String get config string value.
func (s *Section) String(key string) (string, error) {
	if v, ok := s.val[key]; ok {
		return v, nil
	} else {
		return "", &NoKeyError{Key: key, Section: s.sector}
	}
}

// Strings get config []string value.
func (s *Section) Strings(key string) ([]string, error) {
	if v, ok := s.val[key]; ok {
		return strings.Split(v, s.delimit), nil
	} else {
		return nil, &NoKeyError{Key: key, Section: s.sector}
	}
}

// Int get config int value.
func (s *Section) Int(key string) (int64, error) {
	if v, ok := s.val[key]; ok {
		return strconv.ParseInt(v, 10, 64)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.sector}
	}
}

// Uint get config uint value.
func (s *Section) Uint(key string) (uint64, error) {
	if v, ok := s.val[key]; ok {
		return strconv.ParseUint(v, 10, 64)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.sector}
	}
}

// Float get config float value.
func (s *Section) Float(key string) (float64, error) {
	if v, ok := s.val[key]; ok {
		return strconv.ParseFloat(v, 64)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.sector}
	}
}

// Bool get config boolean value.
//
// "yes", "1", "y", "true", "enable" means true.
//
// "no", "0", "n", "false", "disable" means false.
//
// if the specified value unknown then return false.
func (s *Section) Bool(key string) (bool, error) {
	if v, ok := s.val[key]; ok {
		parseBool := func(v string) bool {
			if v == "true" || v == "yes" || v == "1" || v == "y" || v == "enable" {
				return true
			} else if v == "false" || v == "no" || v == "0" || v == "n" || v == "disable" {
				return false
			} else {
				return false
			}
		}
		return parseBool(strings.ToLower(v)), nil
	} else {
		return false, &NoKeyError{Key: key, Section: s.sector}
	}
}

// MemSize Byte get config byte number value.
//
// 1kb = 1k = 1024.
//
// 1mb = 1m = 1024 * 1024.
//
// 1gb = 1g = 1024 * 1024 * 1024.
func (s *Section) MemSize(key string) (int, error) {
	if v, ok := s.val[key]; ok {
		return parseMemory(v)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.sector}
	}
}

func parseMemory(v string) (int, error) {
	unit := Byte
	subIdx := len(v)
	if strings.HasSuffix(v, "k") {
		unit = KB
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "kb") {
		unit = KB
		subIdx = subIdx - 2
	} else if strings.HasSuffix(v, "m") {
		unit = MB
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "mb") {
		unit = MB
		subIdx = subIdx - 2
	} else if strings.HasSuffix(v, "g") {
		unit = GB
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "gb") {
		unit = GB
		subIdx = subIdx - 2
	}
	b, err := strconv.ParseInt(v[:subIdx], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(b) * unit, nil
}

// Duration get config second value.
//
// 1s = 1sec = 1.
//
// 1m = 1min = 60.
//
// 1h = 1hour = 60 * 60.
func (s *Section) Duration(key string) (time.Duration, error) {
	if v, ok := s.val[key]; ok {
		if t, err := parseTime(v); err != nil {
			return 0, err
		} else {
			return time.Duration(t), nil
		}
	} else {
		return 0, &NoKeyError{Key: key, Section: s.sector}
	}
}

func parseTime(v string) (int64, error) {
	unit := int64(time.Nanosecond)
	subIdx := len(v)
	if strings.HasSuffix(v, "ms") {
		unit = int64(time.Millisecond)
		subIdx = subIdx - 2
	} else if strings.HasSuffix(v, "s") {
		unit = int64(time.Second)
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "sec") {
		unit = int64(time.Second)
		subIdx = subIdx - 3
	} else if strings.HasSuffix(v, "m") {
		unit = int64(time.Minute)
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "min") {
		unit = int64(time.Minute)
		subIdx = subIdx - 3
	} else if strings.HasSuffix(v, "h") {
		unit = int64(time.Hour)
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "hour") {
		unit = int64(time.Hour)
		subIdx = subIdx - 4
	}
	b, err := strconv.ParseInt(v[:subIdx], 10, 64)
	if err != nil {
		return 0, err
	}
	return b * unit, nil
}

// Keys return all the section keys.
func (s *Section) Keys() []string {
	var keys []string
	for k := range s.val {
		keys = append(keys, k)
	}
	return keys
}

// Config config
type Config struct {
	// common config
	Common map[string]string // commonKey commonVal1,commonVal2

	// sectors
	Sector map[string]*Section

	// config file path
	File string

	// default config
	Comment string
	Split   string
	Delimit string
}

// New return a new default Config object (Comment = '#', Split = ' ', Delimit = ',')
func New() *Config {
	return &Config{
		Common: make(map[string]string),
		Sector: make(map[string]*Section),

		// default config
		Comment: Comment,
		Split:   Split,
		Delimit: Delimit,
	}
}

// ParseReader parse from io.Reader
func (c *Config) ParseReader(reader io.Reader) error {
	var (
		line      int
		r         = bufio.NewReader(reader)
		sector    *Section
		sectorKey string
		key       string
		val       string
	)
	for {
		// process include
		// process common config
		// process sector
		line++
		row, err := r.ReadString(CRLF)
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(row) == 0 {
			break
		}
		row = strings.TrimSpace(row)

		// process comment
		if len(row) == 0 || strings.HasPrefix(row, c.Comment) {
			// comment or empty line
			continue
		}
		if strings.HasPrefix(row, SectionB) {
			if !strings.HasSuffix(row, SectionE) {
				return fmt.Errorf("no end sector %s at file:%v line:%v", SectionE, c.File, line)
			}

			sectorKey = row[1 : len(row)-1]
			if _, ok := c.Sector[sectorKey]; ok {
				return fmt.Errorf("sector key %v already exists at file:%v line:%v", sectorKey, c.File, line)
			} else {
				sector = &Section{
					delimit: c.Delimit,
					sector:  sectorKey,
					val:     make(map[string]string),
				}
				c.Sector[sectorKey] = sector
			}
			continue
		}

		// key/val in a row
		idx := strings.Index(row, c.Split)
		if idx > 0 {
			key = strings.TrimSpace(row[:idx])
			if len(row) > idx {
				val = strings.TrimSpace(row[idx+1:])
			}
		} else {
			return fmt.Errorf("no split in key row %v at file:%v line:%v", row, c.File, line)
		}

		if sector == nil {
			// process include
			if strings.Contains(row, Include) {
				abs, _ := filepath.Abs(c.File)
				file := path.Join(path.Dir(abs), val)
				if err = c.Parse(file); err != nil {
					return err
				}
			}
			// process common config
			if _, ok := c.Common[key]; ok {
				return fmt.Errorf("same common key %v at file:%v line:%v", key, c.File, line)
			}
			c.Common[key] = val
		} else {
			if c.Sector[sectorKey].val == nil {
				c.Sector[sectorKey].val = make(map[string]string)
			}
			if _, ok := c.Sector[sectorKey].val[key]; ok {
				return fmt.Errorf("section %s already has key: %s at file:%v line:%d", sectorKey, key, c.File, line)
			} else {
				c.Sector[sectorKey].val[key] = val
			}
		}
	}

	return nil
}

// Parse parse file
func (c *Config) Parse(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	c.File = file
	return c.ParseReader(f)
}

// Reload reload config
func (c *Config) Reload() (*Config, error) {
	nc := &Config{
		Common: make(map[string]string),
		Sector: make(map[string]*Section),
		File:   c.File,

		// config
		Comment: c.Comment,
		Split:   c.Split,
		Delimit: c.Delimit,
	}
	err := nc.Parse(c.File)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

// Get get a config section by key.
func (c *Config) Get(section string) *Section {
	s, _ := c.Sector[section]
	return s
}

// GetKey get common key
func (c *Config) GetKey(key string) string {
	return c.Common[key]
}

// GetKeys get common key slice
func (c *Config) GetKeys(key string) []string {
	return strings.Split(c.Common[key], c.Delimit)
}

// Unmarshal unmarshal struct
// memory
const (
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

// timer
const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
)

// Unmarshal unmarshal
func (c *Config) Unmarshal(v interface{}) error {
	// todo
	return nil
}
