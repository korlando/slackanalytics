package slackanalytics

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "strings"
)

var (
  symbols = []rune{'.', ',', '"', '\'', '(', ')', '!', '?', '$', ';', ':'}
)

type Message struct {
  User      string `json:"user"`
  Type      string `json:"type"`
  SubType   string `json:"subtype"`
  Text      string `json:"text"`
  TimeStamp string `json:"ts"`
}

// ReadAllMessages takes in a path to the data folder and returns
// all messages from all channels in no particular order
func ReadAllMessages(dataPath string) (messages []Message, err error) {
  fileInfos, err := ioutil.ReadDir(dataPath)
  if err != nil {
    return
  }
  for _, f := range fileInfos {
    // only channels are dirs
    if !f.IsDir() {
      continue
    }
    channelPath := dataPath + "/" + f.Name()
    channelMessages, err := ReadChannelMessages(channelPath)
    if err != nil {
      continue
    }
    messages = append(messages, channelMessages...)
  }
  return
}

// ReadChannelMessages takes in a path to a channel folder
// and returns all messages from that channel
func ReadChannelMessages(channelPath string) (messages []Message, err error) {
  jsonFiles, err := ioutil.ReadDir(channelPath)
  if err != nil {
    return
  }
  // look at each json file in channel (1 per day)
  for _, j := range jsonFiles {
    file, err := os.Open(channelPath + "/" + j.Name())
    if err != nil {
      continue
    }
    defer file.Close()
    jsonBytes, err := ioutil.ReadAll(file)
    if err != nil {
      continue
    }
    var dayMessages []Message
    err = json.Unmarshal(jsonBytes, &dayMessages)
    if err != nil {
      continue
    }
    messages = append(messages, dayMessages...)
  }
  return
}

// FilterMessagesByUser takes in a slice of messages and user ID
// and returns a filtered slice of those messages by the user
func FilterMessagesByUser(messages []Message, userId string) (filteredMessages []Message) {
  filteredMessages = []Message{}
  for _, m := range messages {
    if m.User == userId {
      filteredMessages = append(filteredMessages, m)
    }
  }
  return
}

// MessageToWords takes in a message and returns a slice of words (strings);
// optionally trims symbols from individual words and converts to lowercase
func MessageToWords(m Message, trimSymbols, lower bool) (words []string) {
  words = strings.Fields(m.Text)
  if !trimSymbols {
    if lower {
      for i, w := range words {
        words[i] = strings.ToLower(w)
      }
    }
    return
  }
  // trim symbols
  for i, w := range words {
    start := 0
    end := len(w)
    for j, b := range w {
      if !isSymbol(b) {
        start = j
        break
      } 
    }
    for j := len(w) - 1; j >= 0; j-- {
      if !isSymbol(rune(w[j])) {
        end = j + 1
        break
      }
    }
    trimmed := string(w[start:end])
    if lower {
      words[i] = strings.ToLower(trimmed)
    } else {
      words[i] = trimmed
    }
  }
  return
}

// isSymbol decides whether a
// rune is in the symbol list
func isSymbol(r rune) bool {
  for _, s := range symbols {
    if s == r {
      return true
    }
  }
  return false
}
