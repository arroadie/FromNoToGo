package main

import "fmt"
import "os"
import "os/exec"
import "strings"
import "github.com/smallfish/simpleyaml"
import "io/ioutil"
import "gopkg.in/yaml.v2"
import "time"

func main()  {

  arguments := os.Args

  if len(arguments) != 3 {
    fmt.Println("You cant run this script without arguments")
    fmt.Println("Usage:\nadd_and_tag_package package.tar.gz relative-path-to-hostclass.yml")
    os.Exit(1)
  }

  fmt.Println(fmt.Sprintf("Will add %s to %s", arguments[1], arguments[2]))

  fmt.Println("Checking that package exists")

  package_name := ""
  hostclass_name := ""
  package_basename := ""

  if strings.Contains(arguments[1], "tar.gz") {
    package_name = arguments[1]
  } else {
    package_name = fmt.Sprintf("%s.tar.gz", arguments[1])
  }

  package_basename = strings.Split(package_name, "-")[0]

  if strings.Contains(arguments[2], ".yml") {
    hostclass_name = "arguments[2]"
  } else {
    hostclass_name = fmt.Sprintf("%s.yml", arguments[2])
  }

  if !strings.Contains(arguments[2], "hostclasses/") {
    hostclass_name = fmt.Sprintf("hostclasses/%s", hostclass_name)
  }

  // Only safe to do that here (?)
  hostclass_base_name := strings.Replace(strings.Replace(hostclass_name, "hostclasses/", "", 1), ".yml", "", 1)

  hostclass_file, err := ioutil.ReadFile(hostclass_name)
  if err != nil {
    fmt.Println("Hostclass file not found, check that the file is available. Hostclass path:", hostclass_name)
    os.Exit(1)
  }

  package_url := fmt.Sprintf("http://config/package/%s", package_name)

  execute("wget", []string{"-S", "--spider", package_url})

  hc_yaml, err := simpleyaml.NewYaml(hostclass_file)
  if err != nil {
    fmt.Println("There was a problem reading the configuration file. Please report it :)")
    os.Exit(1)
  }

  future_production_slice, _ := hc_yaml.Get("packages").Get("production").Array()

  //fmt.Println("slice before", future_production_slice)

  for i := 0; i < len(future_production_slice)-1; i++ {
    if str, ok := future_production_slice[i].(string); ok {
      if strings.HasPrefix(str, package_basename) {
        future_production_slice[i] = strings.Replace(package_name, ".tar.gz", "", 1)
      }
    }
  }

  //fmt.Println("slice after", future_production_slice)
  y, _ := hc_yaml.Map()
  if m, ok := y["packages"].(map[interface{}]interface{}); ok {
    m["production"] = future_production_slice
    y["packages"] = m
  }

  file,_ := yaml.Marshal(y)
  err = ioutil.WriteFile(hostclass_name, []byte(fmt.Sprintf("%s\n",file)), 0644)
  if err != nil {
    fmt.Println(err)
  }

  now := time.Now().Format("2006.01.02_15.04")

  execute("git", []string{"commit", "-am", fmt.Sprintf("'adding %s to %s'", package_name, hostclass_name)})
  execute("git", []string{"tag", fmt.Sprintf("%s-%s", hostclass_base_name, now), "-F",hostclass_name})
  execute("git", []string{"push", "--tags"})
  execute("git", []string{"remote", "update"})
  execute("bin/ops-config-queue", []string{})

  fmt.Println("Process executed successfully")
}

func execute(app string, vals []string) {
  cmd := exec.Command(app, vals...)
  err := cmd.Start()

  if err != nil {
    fmt.Println("Error preparing", app, "Please report errors")
    fmt.Println(err)
    os.Exit(1)
  }

  err = cmd.Wait()
  if err != nil {
    fmt.Println("Error running", app, "with arguments", vals, "Please report errors")
    fmt.Println(err)
    os.Exit(1)
  }
  fmt.Println("executing", cmd.Path, "with", cmd.Args)
}