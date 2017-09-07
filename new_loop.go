package main

import (
        "os"
        "fmt"
        "bufio"
        "strings"
        "sort"
        "path"
	"io/ioutil"
)


func make_table_dev(infile string, clusters []string) (string, []string) {

  // Equivelant to get_seq_name
  basename := path.Base(infile)

  // Open input file
  fi, err := os.Open(infile)
  if err != nil {
    fmt.Println("Error - there was a problem opening the input file. Does it exist?")
    os.Exit(1)
  }
  scanner := bufio.NewScanner(fi)
  defer fi.Close()

  // Keys used to sort dict
  var keys []string
  dict := make(map[string]string)


  // TODO add try/except block to make sure fields is at least len = 2
  for scanner.Scan() {
    fields := strings.Split(scanner.Text(), "\t")
    if len(fields) < 2 {
       fmt.Println("Abnormal number of fields")
       os.Exit(1)
    }
    dict[fields[0]] = fields[1]
    keys = append(keys, fields[0])
  }

  // Iterate over clusters, assigning "0" to the map key if it doesn't exist already
  for _, item := range clusters {
    if _, ok := dict[item]; !ok {
      keys = append(keys, item)
      dict[item] = "0"
    }
  }

  // Sort array
  sort.Strings(keys)

  // Write all values of map to "values" array
  values := make([]string, len(keys)+1)
  values[0] = basename
  for i := 1; i < len(keys)+1; i++ {
    values[i] = dict[keys[i-1]]
  }

  // Now returns basename alone instead of in 2 nested lists. Values should be the same
  return strings.Replace(basename, ".fasta.new_blast.out.filtered.unique", "", -1), values
}


func new_loop(to_iterate, clusters []string) ([]string, [][]string) {

  names := make([]string, len(to_iterate))
  table_list := make([][]string, len(to_iterate))

  ch1 := make(chan string)
  ch2 := make(chan []string)

  for _, file := range to_iterate {

    go func(f string, c []string) {
      out1, out2 := make_table_dev(f, c)
      ch1 <- out1
      ch2 <- out2
    }(file, clusters)
  }

  names_idx := 0
  table_idx := 0

  for i := 0; i < len(to_iterate) * 2; i++ {
    select {

      case name := <-ch1:
        names[names_idx] = name
        names_idx++
      case value := <-ch2:
        table_list[table_idx] = value
        table_idx++
    }
  }

  return names, table_list
}



func main() {
  args := os.Args
  
  arg1, _ := ioutil.ReadFile(args[1])
  arg2, _ := ioutil.ReadFile(args[2])

  to_iterate := strings.Split(string(arg1), ", ")
  clusters := strings.Split(string(arg2), ", ")
  out1, out2 := new_loop(to_iterate, clusters)
  fmt.Println(out1)
  fmt.Println(out2)
}

