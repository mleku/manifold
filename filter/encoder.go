package filter

import (
	"bufio"
	"bytes"
	"encoding/base64"

	"manifold.mleku.dev/errorf"
	"manifold.mleku.dev/ints"
	"manifold.mleku.dev/text"
)

const (
	IDS int = iota
	NOTIDS
	AUTHORS
	NOTAUTHORS
	TAGS
	NOTTAGS
	SINCE
	UNTIL
	SORT
)

var Sentinels = [][]byte{
	[]byte("IDS:"),
	[]byte("NOTIDS:"),
	[]byte("AUTHORS:"),
	[]byte("NOTAUTHORS:"),
	[]byte("TAGS:"),
	[]byte("NOTTAGS:"),
	[]byte("SINCE:"),
	[]byte("UNTIL:"),
	[]byte("SORT:"),
}

// Marshal encodes a filter.F into a byte slice.
// All caps sentinels at the start of lines, fields can appear in any order.
// If Ids are present, any other fields are invalid.
// Sort defaults to descending.
func (f *F) Marshal() (data []byte, err error) {
	buf := new(bytes.Buffer)
	
	// If Ids are present, only include Ids
	if f.Ids != nil && len(f.Ids) > 0 {
		for i, id := range f.Ids {
			if i > 0 {
				buf.WriteByte('\n')
			}
			buf.Write(Sentinels[IDS])
			b := make([]byte, base64.RawURLEncoding.EncodedLen(len(id)))
			base64.RawURLEncoding.Encode(b, id)
			buf.Write(b)
		}
		data = buf.Bytes()
		return
	}
	
	// Otherwise, include all other fields
	var lineCount int
	
	// NotIds
	if f.NotIds != nil && len(f.NotIds) > 0 {
		for _, id := range f.NotIds {
			if lineCount > 0 {
				buf.WriteByte('\n')
			}
			buf.Write(Sentinels[NOTIDS])
			b := make([]byte, base64.RawURLEncoding.EncodedLen(len(id)))
			base64.RawURLEncoding.Encode(b, id)
			buf.Write(b)
			lineCount++
		}
	}
	
	// Authors
	if f.Authors != nil && len(f.Authors) > 0 {
		for _, author := range f.Authors {
			if lineCount > 0 {
				buf.WriteByte('\n')
			}
			buf.Write(Sentinels[AUTHORS])
			b := make([]byte, base64.RawURLEncoding.EncodedLen(len(author)))
			base64.RawURLEncoding.Encode(b, author)
			buf.Write(b)
			lineCount++
		}
	}
	
	// NotAuthors
	if f.NotAuthors != nil && len(f.NotAuthors) > 0 {
		for _, author := range f.NotAuthors {
			if lineCount > 0 {
				buf.WriteByte('\n')
			}
			buf.Write(Sentinels[NOTAUTHORS])
			b := make([]byte, base64.RawURLEncoding.EncodedLen(len(author)))
			base64.RawURLEncoding.Encode(b, author)
			buf.Write(b)
			lineCount++
		}
	}
	
	// Tags
	if f.Tags != nil && len(f.Tags) > 0 {
		for key, values := range f.Tags {
			for _, value := range values {
				if lineCount > 0 {
					buf.WriteByte('\n')
				}
				buf.Write(Sentinels[TAGS])
				if err = text.Write(buf, []byte(key)); err != nil {
					return nil, err
				}
				buf.WriteByte(':')
				b := make([]byte, base64.RawURLEncoding.EncodedLen(len(value)))
				base64.RawURLEncoding.Encode(b, value)
				buf.Write(b)
				lineCount++
			}
		}
	}
	
	// NotTags
	if f.NotTags != nil && len(f.NotTags) > 0 {
		for key, values := range f.NotTags {
			for _, value := range values {
				if lineCount > 0 {
					buf.WriteByte('\n')
				}
				buf.Write(Sentinels[NOTTAGS])
				if err = text.Write(buf, []byte(key)); err != nil {
					return nil, err
				}
				buf.WriteByte(':')
				b := make([]byte, base64.RawURLEncoding.EncodedLen(len(value)))
				base64.RawURLEncoding.Encode(b, value)
				buf.Write(b)
				lineCount++
			}
		}
	}
	
	// Since
	if f.Since != 0 {
		if lineCount > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(Sentinels[SINCE])
		ts := ints.New(f.Since)
		b := ts.Marshal(nil)
		buf.Write(b)
		lineCount++
	}
	
	// Until
	if f.Until != 0 {
		if lineCount > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(Sentinels[UNTIL])
		ts := ints.New(f.Until)
		b := ts.Marshal(nil)
		buf.Write(b)
		lineCount++
	}
	
	// Sort (defaults to descending if not specified)
	if f.Sort != "" && f.Sort != "desc" {
		if lineCount > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(Sentinels[SORT])
		buf.WriteString(f.Sort)
		lineCount++
	}
	
	data = buf.Bytes()
	return
}

// Unmarshal decodes a byte slice into a filter.F.
// All caps sentinels at the start of lines, fields can appear in any order.
// If Ids are present, any other fields are invalid.
// Sort defaults to descending.
func (f *F) Unmarshal(data []byte) (err error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	buf := make([]byte, 1_000_000)
	scanner.Buffer(buf, len(buf))
	
	// Default Sort to descending
	f.Sort = "desc"
	
	var hasIds bool
	
	for scanner.Scan() {
		if scanner.Err() != nil {
			err = scanner.Err()
			return
		}
		
		line := scanner.Bytes()
		
		// Check for IDS sentinel
		if bytes.HasPrefix(line, Sentinels[IDS]) {
			// If we already have other fields and now found IDS, return error
			if hasOtherFields(f) {
				return errorf.E("IDS found but other fields already present")
			}
			
			hasIds = true
			
			// Decode the ID
			id := make([]byte, base64.RawURLEncoding.DecodedLen(len(line)-len(Sentinels[IDS])))
			n, decErr := base64.RawURLEncoding.Decode(id, line[len(Sentinels[IDS]):])
			if decErr != nil {
				return decErr
			}
			id = id[:n]
			
			// Add to Ids slice
			f.Ids = append(f.Ids, id)
			continue
		}
		
		// If we have IDS, other fields are invalid
		if hasIds {
			return errorf.E("other fields found but IDS already present")
		}
		
		// Process other fields
		switch {
		case bytes.HasPrefix(line, Sentinels[NOTIDS]):
			id := make([]byte, base64.RawURLEncoding.DecodedLen(len(line)-len(Sentinels[NOTIDS])))
			n, decErr := base64.RawURLEncoding.Decode(id, line[len(Sentinels[NOTIDS]):])
			if decErr != nil {
				return decErr
			}
			id = id[:n]
			f.NotIds = append(f.NotIds, id)
			
		case bytes.HasPrefix(line, Sentinels[AUTHORS]):
			author := make([]byte, base64.RawURLEncoding.DecodedLen(len(line)-len(Sentinels[AUTHORS])))
			n, decErr := base64.RawURLEncoding.Decode(author, line[len(Sentinels[AUTHORS]):])
			if decErr != nil {
				return decErr
			}
			author = author[:n]
			f.Authors = append(f.Authors, author)
			
		case bytes.HasPrefix(line, Sentinels[NOTAUTHORS]):
			author := make([]byte, base64.RawURLEncoding.DecodedLen(len(line)-len(Sentinels[NOTAUTHORS])))
			n, decErr := base64.RawURLEncoding.Decode(author, line[len(Sentinels[NOTAUTHORS]):])
			if decErr != nil {
				return decErr
			}
			author = author[:n]
			f.NotAuthors = append(f.NotAuthors, author)
			
		case bytes.HasPrefix(line, Sentinels[TAGS]):
			line = line[len(Sentinels[TAGS]):]
			keyEnd := bytes.IndexByte(line, ':')
			if keyEnd == -1 {
				return errorf.E("invalid TAGS format")
			}
			
			key, keyErr := text.Read(bytes.NewBuffer(line[:keyEnd]))
			if keyErr != nil {
				return keyErr
			}
			
			value := make([]byte, base64.RawURLEncoding.DecodedLen(len(line)-keyEnd-1))
			n, decErr := base64.RawURLEncoding.Decode(value, line[keyEnd+1:])
			if decErr != nil {
				return decErr
			}
			value = value[:n]
			
			// Initialize Tags if nil
			if f.Tags == nil {
				f.Tags = make(TagMap)
			}
			
			// Add to Tags map
			keyStr := string(key)
			f.Tags[keyStr] = append(f.Tags[keyStr], value)
			
		case bytes.HasPrefix(line, Sentinels[NOTTAGS]):
			line = line[len(Sentinels[NOTTAGS]):]
			keyEnd := bytes.IndexByte(line, ':')
			if keyEnd == -1 {
				return errorf.E("invalid NOTTAGS format")
			}
			
			key, keyErr := text.Read(bytes.NewBuffer(line[:keyEnd]))
			if keyErr != nil {
				return keyErr
			}
			
			value := make([]byte, base64.RawURLEncoding.DecodedLen(len(line)-keyEnd-1))
			n, decErr := base64.RawURLEncoding.Decode(value, line[keyEnd+1:])
			if decErr != nil {
				return decErr
			}
			value = value[:n]
			
			// Initialize NotTags if nil
			if f.NotTags == nil {
				f.NotTags = make(TagMap)
			}
			
			// Add to NotTags map
			keyStr := string(key)
			f.NotTags[keyStr] = append(f.NotTags[keyStr], value)
			
		case bytes.HasPrefix(line, Sentinels[SINCE]):
			ts := ints.New(int64(0))
			if _, tsErr := ts.Unmarshal(line[len(Sentinels[SINCE]):]); tsErr != nil {
				return tsErr
			}
			f.Since = ts.Int64()
			
		case bytes.HasPrefix(line, Sentinels[UNTIL]):
			ts := ints.New(int64(0))
			if _, tsErr := ts.Unmarshal(line[len(Sentinels[UNTIL]):]); tsErr != nil {
				return tsErr
			}
			f.Until = ts.Int64()
			
		case bytes.HasPrefix(line, Sentinels[SORT]):
			f.Sort = string(line[len(Sentinels[SORT]):])
			
		default:
			return errorf.E("unknown sentinel: '%s'", line)
		}
	}
	
	return
}

// hasOtherFields checks if the filter has any fields other than Ids
func hasOtherFields(f *F) bool {
	return (f.NotIds != nil && len(f.NotIds) > 0) ||
		(f.Authors != nil && len(f.Authors) > 0) ||
		(f.NotAuthors != nil && len(f.NotAuthors) > 0) ||
		(f.Tags != nil && len(f.Tags) > 0) ||
		(f.NotTags != nil && len(f.NotTags) > 0) ||
		f.Since != 0 ||
		f.Until != 0 ||
		(f.Sort != "" && f.Sort != "desc")
}