# 📋 Listea - A task List viewer with a cup of tea

[![Go Report Card](https://goreportcard.com/badge/git.iglou.eu/Laboratory/listea)](https://goreportcard.com/report/git.iglou.eu/Laboratory/listea)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Description
Listea work with Gitea API and GTK (systray)  
Listea is a simple tickets/task viewer on 'notification area'

## Config

**A config file is created at first run** *[See sample file](https://git.iglou.eu/Laboratory/listea/raw/branch/master/sample/config.json)*

| OS        | URI                                                                                 |
| :-------- | :---------------------------------------------------------------------------------- |
| Windows   | `%APPDATA%\listea\config.json`                                                      |
| Mac Os    | `${HOME}/Library/Application Support/listea/config.json`                            |
| Unix Like | `${HOME}/.config/listea/config.json` **or** `${XDG_CONFIG_HOME}/listea/config.json` |

**On config file**   

| ID          | VAR                                       |
| :---------- | :---------------------------------------- |
| `api_url`   | `https://<YOUR_GITEA_INSTACE_DNS>/api/v1` |
| `api_token` | Any Gitea API key token type              |
| `list`      | List of all API request                   |

**Each 'list' instances are composed by**   

| ID              | VAR                                            |
| :-------------- | :--------------------------------------------- |
| `api_request`   | You can use `GET` type listing `​/repos​/issues​/search` and `​/repos​/{owner}​/{repo}​/issues` |
| `query_key`     | Any Parameters like on [Official Swagger/Doc](https://try.gitea.io/api/swagger#/issue/issueSearchIssues)            |

## Install
Install with make `make install` (Build at same time)    
First run create an empty config file at `~/.config/listea/config`

## FAQ
**How do you pronounce Listea ?**   
Listea is pronounced /liːz’ti:/ as in "lis-tea"..

## License
This project is licensed under the MIT License.   
See the [LICENSE](https://git.iglou.eu/Laboratory/listea/src/branch/master/LICENSE) file for the full license text.

## Screenshots
![Tray icon](https://git.iglou.eu/Laboratory/listea/raw/branch/master/screenshot/tray_view.png)   
![Tray Open menu](https://git.iglou.eu/Laboratory/listea/raw/branch/master/screenshot/tray_open.png)