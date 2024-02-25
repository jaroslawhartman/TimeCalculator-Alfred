# Time (and date) difference
(Yet) another time/date diff workflow for Alfred

**WORK IN PROGRESS -- NOTHING USEFUL YET**

# Design and Requirements

Below list of requirements and design considerations

## Frontend

Workflow for Alfred 5. Project:

![alt text](img/Project.png)

## Backend

Developed in `golang`.

## Input formats:
- td `<hh:mm>` `<hh:mm>`
- td `<mm:ss>` `<mm:ss>`- ??
- td `<hh:mm>` `<hh:mm>`
- td `<hh:mm:ss>` `<hh:mm:ss>`
- td `<hh:mm>` `<hh:mm:ss>`- Error
- td `<hh:mm:ss>` `<hh:mm>`- Error
- td `<DD>/<MM>/<YYYY>` `<hh:mm:ss>` `<DD>/<MM>/<YYYY>` `<hh:mm:ss>`

## Parameter checks:
- must me an even number of parameters
- must be 2 or 4
- optional "-" character between
    - Then number 3 or 5
    - 1st and 3nd -- if 3
    - 2nd and 3rd -- if 5

## Output:
- `<d>` days, `<h>` hours, `<m>` minutes, and `<s>` seconds
- `<d.dddd>` days
- `<h.hhh>` hours
- `<m.mmm>` minutes
- `<s>` seconds

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
