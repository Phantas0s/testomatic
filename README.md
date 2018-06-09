# testomatic

`testomatic` is a simple CLI which watch and run unit tests automatically on save.
The result will appear in the terminal as well as a desktop notification.


[Installation](#installation)  
[Configuration file](#configuration-file)  
[Examples](#examples)  
[Contributing](#contributing)  
[Disclaimer](#disclaimer)  

## Installation

You can find the releases of testomatic here: [Github releases](https://github.com/Phantas0s/testomatic/releases)

Here an easy way to install testomatic on Linux using your favorite shell:

```shell
cd /usr/local/bin && sudo curl -LO https://github.com/Phantas0s/testomatic/releases/download/v0.1.0/testomatic && sudo chmod 755 testomatic && cd -
```

You can now normally run `testomatic` from anywhere.

## How it works

When you run `testomatic` it will:

1. Watch any modifications occurring to a set of files. You can specify the files watched using a regex matching filenames.
2. Execute a command each time the watched files are saved. 
3. The path of the file saved is added as an argument to the command line. Depending of your configuration, it can be the path of the file itself, the path of the file's directory or directly the root folder watched.
3. Display the result of the command in the terminal (`stdout`) and in a desktop notification.

This general behavior can be used to run unit tests when you save them.

## Configuration file

testomatic need a valid `yaml` configuration file. The best is to create a `.testomatic` file at the root of your project.

You can as well use a different name for the configuration file: in that case you can run `testomatic --config my-config-name.yml`

The configuration file can change the behavior of `testomatic` drastically to adapt it to your needs. Here the details:

```yaml
watch:
  root: src/Tests
  regex: "Test.php"
  ignore:
    - vendor
    - cache
  ignore_hidden: true
command:
  bin: docker-compose
  scope: current
  ignore_path: true
  abs: false
  options:
    - exec
    - "-T"
    - 'php'
    - bin/phpunit
notification:
  disable: false
  img_failure: /home/superUser/.autotest/images/failure.png
  img_success: /home/superUser/.autotest/images/success.png
  regex_success: ok
  regex_failure: fail
  display_result: true
```

### watch

| attribute     | value                                                                                                            | value type | required | default |
|---------------|------------------------------------------------------------------------------------------------------------------|------------|----------|---------|
| root          | The root folder where your tests are. `testomatic` will watch into this folder and every subfolders recursively. | string     | yes      | *empty* |
| regex         | Every filename matching this regex will be watched.                                                              | string     | yes      | *empty* |
| ignore        | Files or folders you want to ignore (vendor for example)                                                         | array      | no       | *empty* |
| ignore_hidden | Any files or folders beginning by a point `.` won't be watched                                                   | boolean    | no       | `false` |

### command

| attribute   | value                                                                                                                                                           | value type                       | required | default   |
|-------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------|----------|-----------|
| bin         | Path of the command line interface to execute                                                                                                                   | string                           | yes      | *empty*   |
| scope       | The command use the path of the saved file(`current`), the directory of the saved file (`dir`) or simply the root folder defined in `watch` (`all`) as argument | string - `current`, `dir`, `all` | yes      | *empty* |
| abs         | Use the saved file absolute path instead of the relative one                                                                                                    | boolean                          | no       | `false`   |
| ignore_path | Doesn't use the path of the file saved as command line option                                                                                                   | boolean                          | yes      | `false` |
| options     | Options to pass to the command line interface                                                                                                                   | array                            | no       | *empty*   |

### notification

| attribute      | value                                                                | value type | required | default |
|----------------|----------------------------------------------------------------------|------------|----------|---------|
| disable        | Disable the desktop notifications                                    | boolean    | no       | false   |
| img_failure    | Path of image displayed when test is failing                         | string     | no       | *empty* |
| img_success    | Path of the image displayed when test is a success                   | string     | no       | *empty* |
| regex_success  | If the result of the command match this regex, the test is a success | string     | yes      | *empty* |
| regex_failure  | If the result of the command match this regex, the test is a failure | string     | yes      | *empty* |
| display_result | Display the return of the command in the notification box            | boolean    | no       | false   |


## Examples

You will find in the folder `examples` tested configuration files for running `PHPUnit` and `go test`.
I included in the `php` examples how to run tests in docker using `docker-compose` with or without notifications.

## Contributing

Pull request is the way ;)

## Disclaimer

- I only tested `testomatic` on Linux (Arch linux). It might not work on macOS or Windows.
- You can use testomatic to run `Golang` and `PHPUnit` tests automatically. 
The configuration should be flexible enough for you to use it with other test frameworks / languages.
