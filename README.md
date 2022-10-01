# Wars Of Warp
An implementation of a simple map generator for "WarpWar."
The original game was designed by Howard M. Thompson and is copyright by Metagaming Concepts.

Join the Groups.io [WarpWar forum](https://groups.io/g/warpwar/messages) to find out more about this game.

Rules taken from [GC00](http://www.contrib.andrew.cmu.edu/usr/gc00/reviews/warpwar.html) archive.

# Building
1. Install Go.
2. Run `go build` in the root directory.

# Running
## Create a file locally
1. Update the data in `cli/create.go`.
2. Run `./wow create`. 

## Web Server
1. Run `./wow server`.
2. Open the page in your browser.
It defaults to [localhost:8080/wow](http://localhost:8080/wow).
3. There are links for the "standard" map in color or black-and-white, as well as for generating a random map.
4. To create a custom map, add your data to the text area and click the button.
5. The server will return an SVG that you can save.

# Data Example

## CSV
The CSV record contains "name, column, row, economic-value, warp-target".

A star may have multiple warp targets (that's the star the warp line leads to) by adding multiple names at the end of the record.

Here is the CSV data for the "standard" map:

    Adab, 6, 6, 0, Erech, Khafa, Byblos
    Akkad, 7, 16, 3, Kish
    Assur, 12, 10, 2, Nippur, Lagash
    Babylon, 6, 18, 4, Sumer
    Byblos, 2, 6, 3, Adab
    Calah, 8, 4, 1, Nippur
    Elam, 7, 12, 5, Lagash
    Erech, 4, 4, 3, Ur, Adab
    Eridu, 12, 16, 1, Kish, Ugarit
    Girsu, 8, 13, 1, Umma
    Jarmo, 11, 12, 3, Kish
    Isin, 1, 15, 1, Nineveh
    Khafa, 7, 9, 2, Adab
    Kish, 10, 15, 0, Jarmo, Eridu
    Lagash, 9, 11, 1, Assur
    Larsu, 11, 2, 2, Susa
    Mari, 6, 10, 1, Ubaid, Umma
    Mosul, 3, 1, 2, Sippur
    Nineveh, 3, 19, 2, Isin
    Nippur, 10, 7, 1, Calah, Susa, Assur, Lagash
    Sippur, 2, 4, 1, Mosul
    Sumarra, 2, 12, 2, Ubaid, Umma
    Sumer, 4, 16, 0, Umma, Babylon
    Susa, 12, 5, 0, Larsu, Nippur
    Ubaid, 3, 8, 5, Mari, Sumarra
    Ugarit, 11, 20, 2, Eridu
    Umma, 5, 14, 2, Sumarra, Mari, Girsu, Sumer
    Ur, 7, 2, 4, Erech

## JSON
The API expects a single JSON object with the following shape:

    {
        "mono": false,
        "nodes": [
            {"name": "Adab", "col": 6, "row": 6, "econ-value": 0, "warps": ["Erech", "Khafa", "Byblos"]},
            {"name": "Akkad", "col": 7, "row": 16, "econ-value": 3, "warps": ["Kish"]},
            {"name": "Assur", "col": 12, "row": 10, "econ-value": 2, "warps": ["Nippur", "Lagash"]},
            {"name": "Babylon", "col": 6, "row": 18, "econ-value": 4, "warps": ["Sumer"]},
            {"name": "Byblos", "col": 2, "row": 6, "econ-value": 3, "warps": ["Adab"]},
            {"name": "Calah", "col": 8, "row": 4, "econ-value": 1, "warps": ["Nippur"]},
            {"name": "Elam", "col": 7, "row": 12, "econ-value": 5, "warps": ["Lagash"]},
            {"name": "Erech", "col": 4, "row": 4, "econ-value": 3, "warps": ["Ur", "Adab"]},
            {"name": "Eridu", "col": 12, "row": 16, "econ-value": 1, "warps": ["Kish", "Ugarit"]},
            {"name": "Girsu", "col": 8, "row": 13, "econ-value": 1, "warps": ["Umma"]},
            {"name": "Jarmo", "col": 11, "row": 12, "econ-value": 3, "warps": ["Kish"]},
            {"name": "Isin", "col": 1, "row": 15, "econ-value": 1, "warps": ["Nineveh"]},
            {"name": "Khafa", "col": 7, "row": 9, "econ-value": 2, "warps": ["Adab"]},
            {"name": "Kish", "col": 10, "row": 15, "econ-value": 0, "warps": ["Jarmo", "Eridu"]},
            {"name": "Lagash", "col": 9, "row": 11, "econ-value": 1, "warps": ["Assur"]},
            {"name": "Larsu", "col": 11, "row": 2, "econ-value": 2, "warps": ["Susa"]},
            {"name": "Mari", "col": 6, "row": 10, "econ-value": 1, "warps": ["Ubaid", "Umma"]},
            {"name": "Mosul", "col": 3, "row": 1, "econ-value": 2, "warps": ["Sippur"]},
            {"name": "Nineveh", "col": 3, "row": 19, "econ-value": 2, "warps": ["Isin"]},
            {"name": "Nippur", "col": 10, "row": 7, "econ-value": 1, "warps": ["Calah", "Susa", "Assur", "Lagash"]},
            {"name": "Sippur", "col": 2, "row": 4, "econ-value": 1, "warps": ["Mosul"]},
            {"name": "Sumarra", "col": 2, "row": 12, "econ-value": 2, "warps": ["Ubaid", "Umma"]},
            {"name": "Sumer", "col": 4, "row": 16, "econ-value": 0, "warps": ["Umma", "Babylon"]},
            {"name": "Susa", "col": 12, "row": 5, "econ-value": 0, "warps": ["Larsu", "Nippur"]},
            {"name": "Ubaid", "col": 3, "row": 8, "econ-value": 5, "warps": ["Mari", "Sumarra"]},
            {"name": "Ugarit", "col": 11, "row": 20, "econ-value": 2, "warps": ["Eridu"]},
            {"name": "Umma", "col": 5, "row": 14, "econ-value": 2, "warps": ["Sumarra", "Mari", "Girsu", "Sumer"]},
            {"name": "Ur", "col": 7, "row": 2, "econ-value": 4, "warps": ["Erech"]}
        ]
    }


# systemd
See the
[DO Tutorial](https://www.digitalocean.com/community/tutorials/how-to-sandbox-processes-with-systemd-on-ubuntu-20-04)
for details on securing and locking down this as a service.

FWIW, this is my starter:

    /etc/systemd/system# cat wow.service
    [Unit]
    Description=WoW server
    StartLimitIntervalSec=0
    After=network-online.target
    
    [Service]
    Type=simple
    DynamicUser=yes
    PIDFile=/run/wow.pid
    WorkingDirectory=/var/www/wraith.dev/wow
    ExecStart=/var/www/bin/wow server
    ExecReload=/bin/kill -USR1 $MAINPID
    Restart=on-failure
    RestartSec=1
    
    [Install]
    WantedBy=multi-user.target

To get this to work, I built the project and copied the executable into the bin directory.
I then copied the contents of the public directory into the `WorkingDirectory` from the service file.

You can view logs with

    sudo journalctl -u wow.service

