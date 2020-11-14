# ancho

A simple CLI timer application, written in Go.

## Usage examples

### Starting a timebox
```
# Start a timebox for five minutes
ancho box -m 5
```


```
# Start a timebox for one minute and ten seconds
ancho box -m 1 -s 10
```


```
# Start a four-minute timebox and give that timebox a task label
ancho box -m 4 -l "listen to new busta rhymes single"
```


### Listing timeboxes
```
# List all completed timeboxes for the current date
ancho list
```


```
# List all completed timeboxes for November 1st, 2020
ancho list -d 2020-11-01
```
