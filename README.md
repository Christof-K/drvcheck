**Drvcheck** is golang written wrapper around unix `df` command which output drives memory status.

**Drvcheck adds ability of:**
- custom configuration via yaml file
- select which drives you want keep eye on (default all)
- select output memory unit (from KB to JB, default GB)
- select which columns you want to store (default all)
- storing parsed output in CSV file (daily or continuously file, default daily).

# Run executable
Download latest realease package from here.
It contains `drvcheck` executable and `config.yaml` configuration file which have to be at the same location as executable.

Unpack the package and go to drvcheck directory, from there:

Make it executable:
`sudo chmod +x ./drvcheck`

Run it:
`./drvcheck`

> Edit `config.yaml` to specify **drives** and **csv.dir** or leave it as it is, to store csv in the executable file directory

## Crontab

You probably want to dump drive status regularly, to do that you can use `crontab`

type `crontab -e` then add this line
`*/1 * * * * /full/path/to/drvcheck/drvcheck >> /full/path/to/drvcheck/debug.log` 
where `/full/path/to/drvcheck/` is full path to drvcheck directory.

> If you operate on macos authorize **Full Disk Access** for terminal.app in Settings > Security & Privacy > Privacy



More about crontab [here](https://man7.org/linux/man-pages/man5/crontab.5.html "here").


## Licence
**Drvcheck** is under "Do What The F*ck You Want To Public License" (wtfpl) which is self explanatory.


## TODO
- CLI interactive mode with charts
- Alerts of some kind
- DB storing



