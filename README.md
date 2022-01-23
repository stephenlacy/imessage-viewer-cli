# iMessage-viewer-cli
> Dump iMessage chats to pdf

### Building
```
go build -o imessage-viewer main.go
```

### Usage
[Ensure iTerm has full disk access](https://gitlab.com/gnachman/iterm2/-/wikis/Fulldiskaccess)

```
./imessage-viewer <command> <number>

```

### Reading from the local iMessage Sqlite DB
```
./imessage-viewer imessage +14155555555
> +14155555555.pdf
```


### Reading from an (unencrypted) iPhone backup
```
./imessage-viewer iphone +14155555555
> +14155555555.pdf
```


### Current limitations

- Emoji text does not render. This is due to lack of font support in [johnfercher/maroto](github.com/johnfercher/maroto)
- Images are not included. PRs welcome for additional features

