# Time (and date) difference
(Yet) another time/date diff workflow for Alfred

**WORK IN PROGRESS -- NOTHING USEFUL YET**

# Design and Requirements

Below list of requirements and design considerations

## Frontend

Workflow for Alfred 5.

![alt text](img/Usage1.png)

## Backend

Developed in `golang`.

## Input components

Time component formats `<time>`:
 - [X] `<ss>`
 - [X] `<mm:ss>`
 - [X] `<hh:mm:ss>`

Date component formats `<date>`:
 - [ ] If configured `DD/MM/YYYY`
     - `<DD>/<MM>`
     - `<DD>/<MM>/<YYYY>`
 - [ ] If configured `MM/DD/YYYY`
     - `<MM>/<DD>`
     - `<MM>/<DD>/<YYYY>`

Timestamp component formats `<ts>`:
 - [X] Unix timestamp `<dddddddddd>u`, e.g. `1709420400u`

Compount duration component `<period>`:
 -  [X] `<d>d<h>h<m>m<s>s` - in any order
 -  [X] Any component can be ommited, e.g. `1d4h`

Number component `<number>` represents:
 -  [X] Number of seconds `60`
 -  [ ] A number for Span calculations `*` or `/`


## Valid queries
- Duration span (difference) where `<OP>` can be `-` or `+`:
    - [X] `td <time> <OP> <time>` - time difference
    - [ ] `td <date> <time> <OP> <date> <time>` - time difference
- Span calculations, where `<OP>` can be `-` or `+`:
    - [X] `td <time> <OP> <time>`
    - [ ] `td <date> <time> <OP> <time>`
    - [X] `td <time> <OP> <period>`
    - [ ] `td <date> <time> <OP> <period>`
    - [X] `td <timestamp> <OP> <period>`
    - [X] `td <timestamp> <OP> <time>`
- Span calculations, where `<OP>` can be `*` or `/`:
    - [ ] `td <time> <OP> <number>`


## Output:
- [X] `<d>` days, `<h>` hours, `<m>` minutes, and `<s>` seconds
- [ ] `hh:mm:ss` (or `<hh>h<mm>m<ss>s` ?) -- perhaps optional (with AM/PM)
- [X] `<d.ddd>` days
- [X] `<h.hh>` hours
- [X] `<m.mm>` minutes
- [X] `<s>` seconds
- If `<date>` specified, a date will be returned
    - [ ] `DD/MM/YYYY hh:mm:ss`, or
    - [ ] `MM/DD/YYYY hh:mm:ss`

## Unit formatted for singular/plural:
- day/days
- hour/hours
- minut/minutes
- second/seconds

Additionally:
- `0,2,...` is plural, `1` is singular)
- Numbers formatted with thusdands separators, e.g.:
- `999`
- `1,234`
- `1,234,567`
- `1,234,567.890`

## Configurations:

Based on Alfred 5 workflow configuration https://www.alfredapp.com/help/workflows/workflow-configuration/


![alt text](img/Configuration.png)

- Search Key - default `td`
- Date formats
    - `DD/MM/YYYY`
    - `DD/MM`
    - `MM/DD/YYYY`
    - `MM/DD`

## OneUpdater support

Must-have!

# References
* Icon - https://www.flaticon.com/free-icon/duration_5116345
* https://www.alfredapp.com/help/workflows/inputs/script-filter/json/
