# Generate ffprobe output

```
ffprobe -select_streams 0 -print_format compact= -show_frames /Users/andrey/Movies/bbb-540p30.mp4 >bbb-540p30.txt
```

# Run locally

```
go build && ./senderstub -i ~/Movies/bbb-540p30.mp4 -c bbb-540p30.txt
```
